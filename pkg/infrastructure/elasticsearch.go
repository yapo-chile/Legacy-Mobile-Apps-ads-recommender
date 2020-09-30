package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.mpi-internal.com/Yapo/pro-carousel/pkg/interfaces/loggers"
)

type ElasticItem struct {
	ID   string
	Data interface{}
}
type ElasticHandler struct {
	client        *elasticsearch.Client
	batchSize     int
	searchTimeout time.Duration
	logger        loggers.Logger
}
type BulkResponse struct {
	Errors bool       `json:"errors"`
	Items  []ItemBulk `json:"items"`
}
type ItemBulk struct {
	Index IndexBulk `json:"index"`
}
type IndexBulk struct {
	ID     string    `json:"_id"`
	Result string    `json:"result"`
	Status int       `json:"status"`
	Error  ErrorBulk `json:"error"`
}
type ErrorBulk struct {
	Type   string    `json:"type"`
	Reason string    `json:"reason"`
	Cause  CauseBulk `json:"caused_by"`
}
type CauseBulk struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// NewElasticHandlerHandler will create a new instance of a custom http request handler
func NewElasticHandlerHandler(
	maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost, idleConnTimeout, batchSize int, searchTimeout time.Duration,
	addresses string,
	logger loggers.Logger) *ElasticHandler {
	cfg := elasticsearch.Config{
		Addresses: []string{
			addresses,
		},
		Transport: &http.Transport{
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			MaxConnsPerHost:     maxConnsPerHost,
			IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
		},
	}
	es7, _ := elasticsearch.NewClient(cfg)
	return &ElasticHandler{
		client:        es7,
		batchSize:     batchSize,
		searchTimeout: searchTimeout,
		logger:        logger,
	}
}

// Info gets elastic cluster info
func (es *ElasticHandler) Info() (interface{}, error) {
	res, err := es.client.Info()
	if err != nil {
		return nil, err
	}
	res.String() //nolint https://github.com/elastic/go-elasticsearch/issues/123
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		err = fmt.Errorf("error: %s", res.String())
	}
	return res, err
}

// Create generates a new elastic index
func (es *ElasticHandler) Create(index string) error {
	if _, err := es.client.Indices.Delete([]string{index}); err != nil {
		return fmt.Errorf("cannot delete index: %s", err)
	}
	res, err := es.client.Indices.Create(index)
	// issue: https://github.com/elastic/go-elasticsearch/issues/123
	res.String() // nolint
	if err != nil || res.IsError() {
		return fmt.Errorf("cannot create index: %s", err)
	}
	return nil
}

// PutMapping put mapping on a index on elastic search
// mapping byte slice with the poperties json data
// index string with index name
func (es *ElasticHandler) PutMapping(mapping []byte, index string) error {
	res, err := es.client.Indices.PutMapping(
		[]string{index},
		strings.NewReader(string(mapping)),
	)
	// issue: https://github.com/elastic/go-elasticsearch/issues/123
	res.String() // nolint
	if err != nil || res.IsError() {
		return fmt.Errorf("cannot put mapping on index: %s", err)
	}
	return nil
}

// Search gets response from elastic using query
// index string with index name
// query string with query string
// size how many hits are returned
// from in which item the search process begins
func (es *ElasticHandler) Search(index, query string, size, from int) (string, error) {
	if size <= 0 {
		size = 10
	}
	res, err := es.client.Search(
		es.client.Search.WithIndex(index),
		es.client.Search.WithBody(strings.NewReader(query)),
		es.client.Search.WithSize(size),
		es.client.Search.WithFrom(from),
		es.client.Search.WithTimeout(es.searchTimeout),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// Check response status
	response, err := ioutil.ReadAll(res.Body)
	return string(response), err
}

// Bulk insert a data collection in elastic
// collection items to be send
// index string with index name
// action string to indicate which action should be done, ex: index, update, ...
func (es *ElasticHandler) Bulk(collection []ElasticItem, index, action string) error {
	var blk *BulkResponse
	var buf bytes.Buffer
	var res *esapi.Response
	var raw map[string]interface{}
	numErrors, numIndexed, numItems, numBatches, currBatch := 0, 0, 0, 0, 0
	count := len(collection)
	if count%es.batchSize == 0 {
		numBatches = (count / es.batchSize)
	} else {
		numBatches = (count / es.batchSize) + 1
	}
	// Loop over the collection
	for i, item := range collection {
		numItems++
		currBatch = i / es.batchSize
		if i == count-1 {
			currBatch++
		}
		// Prepare the metadata payload
		meta := []byte(fmt.Sprintf(`{ "%s" : { "_id" : "%s" } }%s`, action, item.ID, "\n"))
		// Prepare the data payload: encode to JSON
		data, err := json.Marshal(item.Data)
		if err != nil {
			es.logger.Error("Cannot encode item %d: %s", item.ID, err)
			continue
		}
		// Append newline to the data payload
		data = append(data, "\n"...)
		// Append payloads to the buffer (ignoring write errors)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
		if i > 0 && i%es.batchSize == 0 || i == count-1 {
			es.logger.Error("Bulk [%d de %d] ", currBatch, numBatches)
			// execute bulk
			res, err = es.client.Bulk(bytes.NewReader(buf.Bytes()), es.client.Bulk.WithIndex(index))
			if err != nil {
				es.logger.Error("ElasticSearch Bulk Fails: indexing batch %d: %s", currBatch, err)
			}
			// If the whole request failed, print error and mark all documents as failed
			// issue: https://github.com/elastic/go-elasticsearch/issues/123
			res.String() // nolint
			if res.IsError() {
				numErrors += numItems
				es.bulkError(res, raw)
			} else {
				numErrors, numIndexed = es.bulkOK(res, blk, numErrors, numIndexed)
			}
			res.Body.Close()
			buf.Reset()
			numItems = 0
		}
	}
	return nil
}

func (es *ElasticHandler) bulkError(res *esapi.Response, raw map[string]interface{}) {
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		es.logger.Error("Failure to to parse response body: %s", err)
	} else {
		es.logger.Error("  Error: [%d] %s: %s",
			res.StatusCode,
			raw["error"].(map[string]interface{})["type"],
			raw["error"].(map[string]interface{})["reason"],
		)
	}
}

func (es *ElasticHandler) bulkOK(res *esapi.Response, blk *BulkResponse, numErrors, numIndexed int) (int, int) {
	if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
		es.logger.Error("Failure to to parse response body: %s", err)
	} else {
		for _, d := range blk.Items {
			// ... so for any HTTP status above 201 ...
			if d.Index.Status > 201 {
				// ... increment the error counter ...
				numErrors++
				// ... and print the response status and error information ...
				es.logger.Error("  Error: [%d]: %s: %s: %s: %s",
					d.Index.Status,
					d.Index.Error.Type,
					d.Index.Error.Reason,
					d.Index.Error.Cause.Type,
					d.Index.Error.Cause.Reason,
				)
			} else {
				// ... otherwise increase the success counter.
				numIndexed++
			}
		}
	}
	return numErrors, numIndexed
}
