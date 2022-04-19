package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"

	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/infrastructure"
)

type PactConf struct {
	BrokerHost   string `env:"BROKER_HOST" envDefault:"http://3.229.36.112"`
	BrokerPort   string `env:"BROKER_PORT" envDefault:"80"`
	ProviderHost string `env:"PROVIDER_HOST" envDefault:"http://localhost"`
	ProviderPort string `env:"PROVIDER_PORT" envDefault:"8080"`
	PactPath     string `env:"PACTS_PATH" envDefault:"./pacts"`
}

// Detail has the consumer version details of a pact test
type Detail struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	Href  string `json:"href"`
}

// ConsumerVersion represents the data of the consumer version of a pact test
type ConsumerVersion struct {
	Details Detail `json:"pb:consumer-version"`
}

// JSONTemp stores the variables from the json of the pact test
// that we are going to use
type JSONTemp struct {
	Links        ConsumerVersion `json:"_links"`
	Interactions []interface{}   `json:"interactions"`
}

// A temporary logger is created to be used by the HTTPhandler
type loggerMock struct {
	*testing.T
}

// It has all the logger functions that a normal logger has
func (m *loggerMock) Debug(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Info(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Warn(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Error(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Crit(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}
func (m *loggerMock) Success(format string, params ...interface{}) {
	fmt.Sprintf(format, params...) // nolint: vet,megacheck
}

// Example Provider Pact: How to run me!
// 1. Start the daemon with `./pact-go daemon`
// 2. cd <pact-go>/examples
// 3. go test -v -run TestProvider

func TestProvider(t *testing.T) {
	var conf PactConf
	fmt.Printf("Pact directory: %+v", conf.PactPath)
	infrastructure.LoadFromEnv(&conf)
	var pact = &dsl.Pact{
		Consumer: "ads-recommender",
	}
	files, err := IOReadDir(conf.PactPath)
	if err != nil {
		fmt.Printf("Error in reading files. Error %+v", err)
	}
	for _, file := range files {
		// Verify the Provider with local Pact Files
		h := types.VerifyRequest{
			ProviderBaseURL:       conf.ProviderHost + ":" + conf.ProviderPort,
			PactURLs:              []string{conf.PactPath + "/" + file},
			CustomProviderHeaders: []string{"Authorization: basic e5e5e5e5e5e5e5"},
		}
		_, err := pact.VerifyProvider(t, h)
		if err != nil {
			fmt.Printf("Error verifying the provider.Error %+v\n", err)
			return
		}
	}
}

func TestSendBroker(*testing.T) {
	pactPublisher := &dsl.Publisher{}
	var conf PactConf
	newVer := 1.0
	sendCond := false

	infrastructure.LoadFromEnv(&conf)

	oldPactResponse, currentVer, errOld := getContractInfo(conf.BrokerHost +
		"/pacts/provider/profile-ms/consumer/ads-recommender/latest")
	if errOld != nil {
		if errOld.Error() != "the error code was 404" {
			fmt.Printf("Error getting the contract from the broker: +%v\n", errOld)
			return
		}
	}

	newPactResponse, errNew := getJSONPactFile(conf)
	if errNew != nil {
		fmt.Printf("Error getting the contract from the file: +%v\n", errNew)
		return
	}

	if oldPactResponse == nil {
		sendCond = true
	} else if !reflect.DeepEqual(oldPactResponse, newPactResponse) {
		sendCond = true
		newVer = currentVer + 0.1
	}

	if sendCond {
		err := pactPublisher.Publish(types.PublishRequest{
			PactURLs:        []string{"./pacts/ads-recommender.json"},
			PactBroker:      conf.BrokerHost + ":" + conf.BrokerPort,
			ConsumerVersion: fmt.Sprintf("%.1f", newVer),
			Tags:            []string{"ads-recommender"},
		})
		if err != nil {
			fmt.Printf("Error with the Pact Broker server. Error %+v\n", err)
			return
		}
	}
}

func IOReadDir(root string) ([]string, error) {
	files := make([]string, 0)
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

// Method that gets the last version of the contract between consumer (ads-recommender) and it's provider
// Returns the contract and the last version of the contract or error
func getContractInfo(url string) (interface{}, float64, error) {
	var confContracts infrastructure.Config
	var result JSONTemp
	infrastructure.LoadFromEnv(&confContracts)
	logger := &loggerMock{}

	HTTPHandler := infrastructure.NewHTTPHandler(logger)
	httprequest := HTTPHandler.NewRequest().
		SetMethod("GET").
		SetPath(url)
	publishedContract, err := HTTPHandler.Send(httprequest)
	if err != nil {
		return nil, -1, err
	}
	resp := fmt.Sprintf("%s", publishedContract)

	err = json.Unmarshal([]byte(resp), &result)

	if (err != nil) || (len(result.Interactions) < 1) {
		return nil, -1, err
	}
	pactResponse := result.Interactions[1].(map[string]interface{})
	delete(pactResponse, "_id")
	versionFloat, err := strconv.ParseFloat(result.Links.Details.Name, 64)
	if err != nil {
		return nil, -1, err
	}
	return pactResponse, versionFloat, nil
}

func getJSONPactFile(conf PactConf) (interface{}, error) {
	var result JSONTemp
	file, err := IOReadDir(conf.PactPath)
	if err != nil {
		return nil, err
	}
	pactFileToSend, err := ioutil.ReadFile(conf.PactPath + "/" + file[0])
	if err != nil {
		return nil, err
	}
	resp := fmt.Sprintf("%s", pactFileToSend)
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return nil, err
	}
	pactResponse := result.Interactions[1].(map[string]interface{})
	return pactResponse, err
}
