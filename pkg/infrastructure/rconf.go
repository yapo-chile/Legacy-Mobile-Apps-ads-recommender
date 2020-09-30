package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

// RconfLogger Interface to define log functions needed
type RconfLogger interface {
	Debug(format string, params ...interface{})
	Error(format string, params ...interface{})
	Info(format string, params ...interface{})
}

// Rconf struct to contain data or remote configuration service
// now is working with Etcd Service
type Rconf struct {
	Log     RconfLogger
	Content *EtcdContent
}

// EtcdContent response from Etcd Service
type EtcdContent struct {
	Action string   `json:"action"`
	Node   EtcdNode `json:"node"`
}

// EtcdNode content one item from conf
type EtcdNode struct {
	Key           string     `json:"key"`
	Value         string     `json:"value"`
	ModifiedIndex int        `json:"modifiedIndex"`
	CreatedIndex  int        `json:"createdIndex"`
	Nodes         []EtcdNode `json:"nodes"`
	IsDir         bool       `json:"dir"`
}

func readAll(resp *http.Response, log RconfLogger) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		// In case of 404 response
		if resp.StatusCode == http.StatusNotFound {
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				log.Error("ReadError %s", readErr.Error())
			} else {
				log.Error("Response: %s", string(body))
			}
		}

		return nil, fmt.Errorf("response code invalid: %s", resp.Status)
	}

	// Read the body from response
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Error("Error on read body: %s", readErr.Error())
		return nil, readErr
	}

	return body, nil
}

// NewRconf remote configuration loader
// This method get the configuration file from remote host.
// On success returns a Rconf struct, if not success returns an Error and Rconf
// with an EtcdContent empty
func NewRconf(host, path, prefix string, log RconfLogger) (*Rconf, error) {
	// Init Rconf
	rconf := &Rconf{Log: log, Content: &EtcdContent{}}

	var netClient = &http.Client{
		Timeout: time.Second * 100, //nolint: gomnd
	}
	// Build URL
	url := fmt.Sprintf("%s%s%s", host, prefix, path)
	resp, err := netClient.Get(url)

	if err != nil {
		log.Error("Error to get url %s", url)
		return rconf, err
	}
	// Close body after all
	defer resp.Body.Close() // nolint: errcheck

	body, err := readAll(resp, log)
	if err != nil {
		return rconf, err
	}
	// Parse body to json
	if jsonErr := json.Unmarshal(body, &rconf.Content); jsonErr != nil {
		log.Error("Error %s", jsonErr.Error())
		return rconf, jsonErr
	}

	log.Info("Conf %s loaded", url)

	return rconf, nil
}

// Get gets the result of a GET method with the given key
func (v Rconf) Get(key string) string {
	if v.Content == nil {
		v.Log.Error("Empty conf")
		return ""
	}

	if v.Content.Node.IsDir {
		v.Log.Error("Conf %s is a dir, i can not get some value", v.Content.Node.Key)
		return ""
	}

	return gjson.Get(v.Content.Node.Value, key).String()
}

// Translate implemented for language use
func (v Rconf) Translate(key string) string {
	msg := v.Get(key)
	v.Log.Debug("Translate: %s, %s", key, msg)

	return msg
}
