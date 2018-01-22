package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

//BitClient Client for bitflyer private api
type BitClient struct {
	endpoint   *url.URL
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

//NewBitClient Create new BitClient
func NewBitClient() *BitClient {
	epStr := "https://api.bitflyer.jp"
	endpoint, err := url.ParseRequestURI(epStr)
	if err != nil {
		panic("Endpoint URL is wrong")
	}
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
func (client *BitClient) NewRequest(pathDir string, method string, bodyJSON string) (*http.Request, error) {
	reqURL := *client.endpoint
	reqURL.Path = path.Join(client.endpoint.Path, pathDir)
	body := bytes.NewReader([]byte(bodyJSON))
	req, err := http.NewRequest(method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("ACCESS-KEY", client.apiKey)
	timeStamp := fmt.Sprint(time.Now().Unix())
	req.Header.Set("ACCESS-TIMESTAMP", timeStamp)
	message := timeStamp + method + pathDir + bodyJSON
	hash := hmac.New(sha256.New, []byte(client.apiSecret))
	hash.Write([]byte(message))
	req.Header.Set("ACCESS-SIGN", hex.EncodeToString(hash.Sum(nil)))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

//Do Do request
func (client *BitClient) Do(request *http.Request) (*http.Response, error) {
	return client.httpClient.Do(request)
}
