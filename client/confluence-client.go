package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type httpClient interface {
	Do(request *http.Request) (response *http.Response, err error)
}

// ConfluenceClient is the primary client to the Confluence API
type ConfluenceClient struct {
	client http.Client
	//	username string
	//	password string
	baseURL string
	debug   bool
	//	usetoken bool
}

// OperationOptions holds all the options that apply to the specified operation
type OperationOptions struct {
	Title         string
	SpaceKey      string
	Filepath      string
	BodyOnly      bool
	StripImgs     bool
	AncestorTitle string
	AncestorID    int64
}

// Client returns a new instance of the client
func NewClient(httpClient httpClient, baseURL string) (*ConfluenceClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	// ensure the baseURL contains a trailing slash so that all paths are preserved in later calls
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		client:  httpClient,
		baseURL: parsedBaseURL,
	}

	return c, nil
}

func (s *AuthenticationService) SetBasicAuth(username, password string) {
	s.username = username
	s.password = password
	s.authType = authTypeBasic
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
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func (c *ConfluenceClient) doRequest(method, url string, content, responseContainer interface{}) ([]byte, *http.Response) {
	b := new(bytes.Buffer)
	if content != nil {
		err := json.NewEncoder(b).Encode(content)
		if err != nil {
			log.Fatal(err)
		}
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
		err := json.Unmarshal(contents, responseContainer)
		if err != nil {
			log.Fatal(err)
		}
	}
	return contents, response
}

func (c *ConfluenceClient) DoGetPage(method, url string, content interface{}) ([]byte, *http.Response) {
	b := new(bytes.Buffer)
	if content != nil {
		err := json.NewEncoder(b).Encode(content)
		if err != nil {
			log.Fatal(err)
		}

	}
	furl := c.baseURL + url // How to fix this for Hierarchies report?
	//furl := url
	if c.debug {
		log.Println("Full URL", furl)
		log.Println("JSON Content:", b.String())
	}
	request, err2 := http.NewRequest(method, furl, b)
	if err2 != nil {
		log.Fatal(err2)
	}
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
