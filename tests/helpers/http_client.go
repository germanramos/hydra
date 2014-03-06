package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HTTPClientHelper struct{}

func NewHTTPClientHelper() *HTTPClientHelper {
	return new(HTTPClientHelper)
}

// Creates a new HTTP client with KeepAlive disabled.
func newHTTPClient() *http.Client {
	return &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	// return &http.Client{}
}

// Reads the body from the response and closes it.
func (h HTTPClientHelper) ReadBody(resp *http.Response) []byte {
	if resp == nil {
		return []byte{}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}

// Reads the body from the response and parses it as JSON Object.
func (h HTTPClientHelper) ReadBodyJsonArray(resp *http.Response) ([]map[string]interface{}, error) {
	var m []map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&m)
	if err != nil {
		panic(fmt.Sprintf("HTTP body JSON parse error: %v", err))
	}
	return m, nil
}

// Reads the body from the response and parses it as JSON Object.
func (h HTTPClientHelper) ReadBodyJsonObject(resp *http.Response) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&m)
	if err != nil {
		panic(fmt.Sprintf("HTTP body JSON parse error: %v", err))
	}
	return m, nil
}

func (h HTTPClientHelper) Get(url string) (*http.Response, error) {
	return h.send("GET", url, "", nil)
}

func (h HTTPClientHelper) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	return h.send("POST", url, bodyType, body)
}

func (h HTTPClientHelper) PostForm(url string, data url.Values) (*http.Response, error) {
	return h.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (h HTTPClientHelper) Put(url string, bodyType string, body io.Reader) (*http.Response, error) {
	return h.send("PUT", url, bodyType, body)
}

func (h HTTPClientHelper) PutForm(url string, data url.Values) (*http.Response, error) {
	return h.Put(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (h HTTPClientHelper) Delete(url string, bodyType string, body io.Reader) (*http.Response, error) {
	return h.send("DELETE", url, bodyType, body)
}

func (h HTTPClientHelper) DeleteForm(url string, data url.Values) (*http.Response, error) {
	return h.Delete(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (h HTTPClientHelper) send(method string, url string, bodyType string, body io.Reader) (*http.Response, error) {
	c := newHTTPClient()
	fmt.Println(method)
	fmt.Println(url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println("ERROR")
		return nil, err
	}
	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	return c.Do(req)
	// res, err := c.Do(req)
	// if err != nil {
	// 	fmt.Println("ERROR")
	// 	fmt.Println(err.Error())
	// }
	// b, _ := ioutil.ReadAll(res.Body)
	// res.Body.Close()
	// fmt.Println(string(b))
	// return res, err
}
