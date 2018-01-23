package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//LineClient Client for line api
type LineClient struct {
	isActive   bool
	endpoint   string
	userID     string
	token      string
	httpClient *http.Client
}

//TextMessage Line TextMessage
type TextMessage struct {
	TypeText string `json:"type"`
	Message  string `json:"text"`
}

//PushTextContent Line PushTextContent
type PushTextContent struct {
	To       string        `json:"to"`
	Messages []TextMessage `json:"messages"`
}

func (client *LineClient) lineClientError(message string) {
	fmt.Println("LineClient Error: " + message)
	client.isActive = false
}

//NewLineClient Create new LineClient
func NewLineClient() *LineClient {
	endpoint := "https://api.line.me"
	userID := os.Getenv("LINE_USER_ID")
	if userID == "" {
		panic("LINE_USER_ID")
	}
	token := os.Getenv("LINE_TOKEN")
	if token == "" {
		panic("LINE_TOKEN")
	}
	return &LineClient{true, endpoint, userID, token, &http.Client{}}
}

//NewRequest Make request of LineClient
func (client *LineClient) NewRequest(pathDir string, method string, bodyJSON string) *http.Request {
	if client.isActive == false {
		return nil
	}
	reqURL := client.endpoint + pathDir
	body := bytes.NewReader([]byte(bodyJSON))
	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		client.lineClientError("New Request")
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+client.token)
	req.Header.Set("Content-Type", "application/json")
	return req
}

//Do Do request
func (client *LineClient) Do(request *http.Request) (*http.Response, error) {
	if client.isActive == false {
		return nil, nil
	}
	return client.httpClient.Do(request)
}

//GetResponseBody Get response body string
func (client *LineClient) GetResponseBody(request *http.Request) []byte {
	if client.isActive == false {
		return nil
	}
	res, errReq := client.Do(request)
	if errReq != nil {
		client.lineClientError("Do Request")
		return nil
	}
	defer res.Body.Close()
	byteRet, errRet := ioutil.ReadAll(res.Body)
	if errRet != nil {
		client.lineClientError("Response Read")
		return nil
	}
	if res.StatusCode != 200 {
		client.lineClientError(fmt.Sprintf("Response Code = %d\n%s", res.StatusCode, string(byteRet)))
		return nil
	}
	return byteRet
}

//JSONDecode Decode response body json
func (client *LineClient) JSONDecode(byteBody []byte, jsonData interface{}) {
	if client.isActive == false {
		return
	}
	if err := json.Unmarshal(byteBody, jsonData); err != nil {
		client.lineClientError("JSON Decode")
		return
	}
}

//JSONEncode Encode post body json
func (client *LineClient) JSONEncode(jsonData interface{}) string {
	if client.isActive == false {
		return ""
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		client.lineClientError("JSON Encode")
		return ""
	}
	return string(jsonBytes)
}

//Post Post request
func (client *LineClient) Post(pathDir string, body interface{}, response interface{}) {
	if client.isActive == false {
		return
	}
	bodyString := client.JSONEncode(body)
	request := client.NewRequest(pathDir, "POST", bodyString)
	responseBody := client.GetResponseBody(request)
	if response == nil {
		return
	}
	client.JSONDecode(responseBody, response)
}

//PushTextMessage Push Line message
func (client *LineClient) PushTextMessage(message string) {
	if client.isActive == false {
		return
	}
	pathDir := "/v2/bot/message/push"
	textMessage := TextMessage{"text", message}
	body := PushTextContent{client.userID, []TextMessage{textMessage}}
	client.Post(pathDir, body, nil)
}
