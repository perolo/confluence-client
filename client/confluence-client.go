package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//ConfluenceClient is the primary client to the Confluence API
type ConfluenceClient struct {
	username string
	password string
	baseURL  string
	debug    bool
	usetoken bool
	client   *http.Client
}

//OperationOptions holds all the options that apply to the specified operation
type OperationOptions struct {
	Title         string
	SpaceKey      string
	Filepath      string
	BodyOnly      bool
	StripImgs     bool
	AncestorTitle string
	AncestorID    int64
}

//Client returns a new instance of the client
func Client(config *ConfluenceConfig) *ConfluenceClient {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	return &ConfluenceClient{
		username: config.Username,
		password: config.Password,
		baseURL:  config.URL,
		usetoken: config.UseToken,
		debug:    config.Debug,
		client: &http.Client{
			Timeout: 60 * time.Second, Transport: tr,
		},
	}
}

func SetTokenAuth(r *http.Request, password string) {
	r.Header.Set("Authorization", "Bearer "+password)
}

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
func (c *ConfluenceClient) doRequest(method, url string, content, responseContainer interface{}) ([]byte, *http.Response) {
	b := new(bytes.Buffer)
	if content != nil {
		json.NewEncoder(b).Encode(content)
	}
	furl := c.baseURL + url
	if c.debug {
		log.Println("Full URL", furl)
		log.Println("JSON Content:", b.String())
	}
	request, err := http.NewRequest(method, furl, b)
	if c.usetoken {
		SetTokenAuth(request, c.password)
	} else {
		request.SetBasicAuth(c.username, c.password)
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Fatal(err)
	}
	if c.debug {
		log.Printf("Sending request to services: \n %s", formatRequest(request))
	}
	response, err := c.client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if c.debug {
		log.Println("Response received, processing response...")
		log.Println("Response status code is", response.StatusCode)
		log.Println(response.Status)
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode < 200 || response.StatusCode > 300 {
		log.Println("Bad response code received from server: ", response.Status)
	} else {
		json.Unmarshal(contents, responseContainer)
	}
	return contents, response
}

func (c *ConfluenceClient) doGetPage(method, url string, content interface{}) ([]byte, *http.Response) {
	b := new(bytes.Buffer)
	if content != nil {
		json.NewEncoder(b).Encode(content)
	}
	//	furl := c.baseURL + url // How to fix this for Hierarchies report?
	furl := url
	if c.debug {
		log.Println("Full URL", furl)
		log.Println("JSON Content:", b.String())
	}
	request, err := http.NewRequest(method, furl, b)
	request.Header.Set("X-Atlassian-Token", "nocheck")
	if c.usetoken {
		SetTokenAuth(request, c.password)
	} else {
		request.SetBasicAuth(c.username, c.password)
	}
	if c.debug {
		log.Println("Sending request to services...")
	}
	response, err := c.client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if c.debug {
		log.Println("Response received, processing response...")
		log.Println("Response status code is", response.StatusCode)
		log.Println(response.Status)
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode < 200 || response.StatusCode > 300 {
		log.Println("Bad response code received from server: ", response.Status)
	}
	return contents, response
}
