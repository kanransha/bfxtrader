package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//BitClient Client for bitflyer private api
type BitClient struct {
	endpoint   string
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

func bitClientError(message string) {
	panic("BitClient Error: " + message)
}

//NewBitClient Create new BitClient
func NewBitClient() *BitClient {
	endpoint := "https://api.bitflyer.jp"
	apiKey := os.Getenv("BITFLYER_KEY")
	if apiKey == "" {
		panic("API Key")
	}
	apiSecret := os.Getenv("BITFLYER_SECRET")
	if apiSecret == "" {
		panic("API Secret")
	}
	return &BitClient{endpoint, apiKey, apiSecret, &http.Client{}}
}

//NewRequest Make request of BitClient
func (client *BitClient) NewRequest(pathDir string, method string, bodyJSON string) *http.Request {
	reqURL := client.endpoint + pathDir
	body := bytes.NewReader([]byte(bodyJSON))
	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		bitClientError("New Request")
	}
	req.Header.Set("ACCESS-KEY", client.apiKey)
	timeStamp := fmt.Sprint(time.Now().Unix())
	req.Header.Set("ACCESS-TIMESTAMP", timeStamp)
	message := timeStamp + method + pathDir + bodyJSON
	hash := hmac.New(sha256.New, []byte(client.apiSecret))
	hash.Write([]byte(message))
	req.Header.Set("ACCESS-SIGN", hex.EncodeToString(hash.Sum(nil)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

//Do Do request
func (client *BitClient) Do(request *http.Request) (*http.Response, error) {
	return client.httpClient.Do(request)
}

//GetResponseBody Get response body string
func (client *BitClient) GetResponseBody(request *http.Request) []byte {
	res, errReq := client.Do(request)
	if errReq != nil {
		bitClientError("Do Request")
	}
	defer res.Body.Close()
	byteRet, errRet := ioutil.ReadAll(res.Body)
	if errRet != nil {
		bitClientError("Response Read")
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 401 {
			fmt.Println("Request Header: ", request.Header)
		}
		bitClientError(fmt.Sprintf("Response Code = %d\n%s", res.StatusCode, string(byteRet)))
	}
	return byteRet
}

//JSONDecode Decode response body json
func (client *BitClient) JSONDecode(byteBody []byte, jsonData interface{}) {
	if err := json.Unmarshal(byteBody, jsonData); err != nil {
		bitClientError("JSON Decode")
	}
}

//JSONEncode Encode post body json
func (client *BitClient) JSONEncode(jsonData interface{}) string {
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		bitClientError("JSON Encode")
	}
	return string(jsonBytes)
}

//Post Post request
func (client *BitClient) Post(pathDir string, body interface{}, response interface{}) {
	bodyString := client.JSONEncode(body)
	request := client.NewRequest(pathDir, "POST", bodyString)
	responseBody := client.GetResponseBody(request)
	if response == nil {
		return
	}
	client.JSONDecode(responseBody, response)
}

//Get Get request
func (client *BitClient) Get(pathDir string, query string, response interface{}) {
	pathWithQuery := pathDir
	if query != "" {
		pathWithQuery = pathWithQuery + "?" + query
	}
	request := client.NewRequest(pathWithQuery, "GET", "")
	responseBody := client.GetResponseBody(request)
	client.JSONDecode(responseBody, response)
}
