package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

type PageOptions struct {
	// StartAt: The starting index of the returned projects. Base index: 0.
	Start int `url:"start,omitempty"`
	// MaxResults: The maximum number of projects to return per page. Default: 50.
	Limit int `url:"limit,omitempty"`
	// Expand: Expand specific sections in the returned issues
	Type string `url:"type,omitempty"`
}

//SearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetPageById(id string) (results *ConfluencePage) {
	results = &ConfluencePage{}
	c.doRequest("GET", "/rest/api/content/"+id+"?expand=body.view", nil, results)
	return results
}

//SearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetPageByIdAncestor(id string) (results *ConfluencePage2) {
	results = &ConfluencePage2{}
	c.doRequest("GET", "/rest/api/content/"+id+"?expand=ancestors", nil, results)
	return results
}

func (c *ConfluenceClient) GetPages(space string, options *PageOptions) (results *ConfluencePages) {
	var path string
	if options == nil {
		path = fmt.Sprintf("/rest/api/space/%s/content", space)
	} else {
		path = fmt.Sprintf("/rest/api/space/%s/content?start=%v&limit=%v&type=%s", space, options.Start, options.Limit, options.Type)
	}

	results = &ConfluencePages{}
	c.doRequest("GET", path, nil, results)
	return results
}

func (c *ConfluenceClient) GetContent(content string, options *PageOptions) (results *ConfluencePageSearch) {
	var path string
	if options == nil {
		path = fmt.Sprintf("/rest/api/content?%s", content)
	} else {
		path = fmt.Sprintf("/rest/api/content?%s&start=%v&limit=%v", content, options.Start, options.Limit)
	}

	results = &ConfluencePageSearch{}
	c.doRequest("GET", path, nil, results)
	return results
}

func (c *ConfluenceClient) GetPage(url string) ([]byte, *http.Response) {
	contents, response := c.doGetPage("GET", url, nil)
	return contents, response
}

func (c *ConfluenceClient) GetPageAttachmentById(id string, name string) (results *ConfluenceAttachmnetSearch, data []byte, err error) {

	u := url.URL{
		Path: fmt.Sprintf("/rest/api/content/%s/child/attachment", id),
	}
	uv := url.Values{}
	if name != "" {
		uv.Add("filename", name)
	}

	u.RawQuery = uv.Encode()

	//	path := fmt.Sprintf("/rest/api/content/%s/child/attachment?filename=%s", id, name)

	results = &ConfluenceAttachmnetSearch{}
	c.doRequest("GET", u.String(), nil, results)

	if results.Size == 1 {
		fmt.Printf("Attachment: %s\n", results.Results[0].Title)

		content, resp := c.GetPage(results.Results[0].Links["self"])

		if resp.StatusCode == 200 {
			//			fmt.Printf("Content: %s\n", content)
			return results, content, nil
		} else {
			return results, nil, fmt.Errorf("Bad response code received from server: %v", resp.Status)
		}
	}
	return results, nil, fmt.Errorf("Failed to get attachment: %s", name)
}

func (c *ConfluenceClient) GetPageAttachmentById2(id string, name string) (retv *ConfluenceAttachment, data []byte, err error) {
	path := fmt.Sprintf("/rest/api/content/%s/child/attachment?filename=%s", id, url.QueryEscape(name))

	results := &ConfluenceAttachmnetSearch{}
	c.doRequest("GET", path, nil, results)

	for _, theRes := range results.Results {

		fmt.Printf("Attachment: %s\n", theRes.Title)

		if theRes.Title == name {
			content, resp := c.GetPage(theRes.Links["download"])

			if resp.StatusCode == 200 {
				//fmt.Printf("Content: %s\n", content)
				return &theRes, content, nil
			} else {
				return &theRes, nil, fmt.Errorf("Bad response code received from server: %v", resp.Status)
			}
		}
	}
	return nil, nil, fmt.Errorf("Failed to get attachment: %s", name)
}

func (c *ConfluenceClient) UpdateAttachment(id string, attid string, attName string, newFilePath string, com string) (contents []byte, retType *ConfluenceAttachment, err error) {

	path := fmt.Sprintf("/rest/api/content/%s/child/attachment/%s/data", id, attid)

	// Open the file
	file, err := os.Open(newFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to open file %s", newFilePath)
	}
	// Close the file later
	defer file.Close()

	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", attName)
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form file")
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to copy file %s", newFilePath)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("minorEdit")
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form field ")
	}

	_, err = fieldWriter.Write([]byte("true"))
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form field value")
	}
	// Populate other fields
	fieldWriter2, err := multiPartWriter.CreateFormField("comment")
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form field ")
	}

	_, err = fieldWriter2.Write([]byte(com))
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form field value")
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	req, err := http.NewRequest("POST", c.baseURL+path, &requestBody)

	if err != nil {
		return nil, nil, err
	}

	//fmt.Println(requestBody.String())

	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	req.Header.Set("X-Atlassian-Token", "nocheck")
	if c.usetoken {
		SetTokenAuth(req, c.password)
	} else {
		req.SetBasicAuth(c.username, c.password)
	}

	// Do the request
	response, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()
	if c.debug {
		log.Println("Response received, processing response...")
		log.Println("Response status code is", response.StatusCode)
		log.Println(response.Status)
	}
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return contents, nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 300 {
		log.Println("Bad response code received from server: ", response.Status)
		return contents, nil, fmt.Errorf("Bad response code received from server: %s ", response.Status)
	} else {
		json.Unmarshal(contents, retType)
	}
	return contents, retType, nil
}

type ConfluenceAddAttachment struct {
	Results []struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Status  string `json:"status"`
		Title   string `json:"title"`
		Version struct {
			By struct {
				Type           string `json:"type"`
				Username       string `json:"username"`
				UserKey        string `json:"userKey"`
				ProfilePicture struct {
					Path      string `json:"path"`
					Width     int    `json:"width"`
					Height    int    `json:"height"`
					IsDefault bool   `json:"isDefault"`
				} `json:"profilePicture"`
				DisplayName string `json:"displayName"`
				Links       struct {
					Self string `json:"self"`
				} `json:"_links"`
				Expandable struct {
					Status string `json:"status"`
				} `json:"_expandable"`
			} `json:"by"`
			When      time.Time `json:"when"`
			Message   string    `json:"message"`
			Number    int       `json:"number"`
			MinorEdit bool      `json:"minorEdit"`
			Hidden    bool      `json:"hidden"`
			Links     struct {
				Self string `json:"self"`
			} `json:"_links"`
			Expandable struct {
				Content string `json:"content"`
			} `json:"_expandable"`
		} `json:"version"`
		Container struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Status     string `json:"status"`
			Title      string `json:"title"`
			Extensions struct {
				Position string `json:"position"`
			} `json:"extensions"`
			Links struct {
				Webui  string `json:"webui"`
				Edit   string `json:"edit"`
				Tinyui string `json:"tinyui"`
				Self   string `json:"self"`
			} `json:"_links"`
			Expandable struct {
				Container    string `json:"container"`
				Metadata     string `json:"metadata"`
				Operations   string `json:"operations"`
				Children     string `json:"children"`
				Restrictions string `json:"restrictions"`
				History      string `json:"history"`
				Ancestors    string `json:"ancestors"`
				Body         string `json:"body"`
				Version      string `json:"version"`
				Descendants  string `json:"descendants"`
				Space        string `json:"space"`
			} `json:"_expandable"`
		} `json:"container"`
		Metadata struct {
			Comment    string `json:"comment"`
			MediaType  string `json:"mediaType"`
			Expandable struct {
				Currentuser string `json:"currentuser"`
				Properties  string `json:"properties"`
				Frontend    string `json:"frontend"`
				EditorHTML  string `json:"editorHtml"`
			} `json:"_expandable"`
		} `json:"metadata"`
		Extensions struct {
			MediaType string `json:"mediaType"`
			FileSize  int    `json:"fileSize"`
			Comment   string `json:"comment"`
		} `json:"extensions"`
		Links struct {
			Webui    string `json:"webui"`
			Download string `json:"download"`
			Self     string `json:"self"`
		} `json:"_links"`
		Expandable struct {
			Operations   string `json:"operations"`
			Children     string `json:"children"`
			Restrictions string `json:"restrictions"`
			History      string `json:"history"`
			Ancestors    string `json:"ancestors"`
			Body         string `json:"body"`
			Descendants  string `json:"descendants"`
			Space        string `json:"space"`
		} `json:"_expandable"`
	} `json:"results"`
	Size  int `json:"size"`
	Links struct {
		Base    string `json:"base"`
		Context string `json:"context"`
	} `json:"_links"`
}

func (c *ConfluenceClient) AddAttachment(id string, attName string, newFilePath string, com string) (contents []byte, retType *ConfluenceAddAttachment, err error) {

	path := fmt.Sprintf("/rest/api/content/%s/child/attachment", id)

	// Open the file
	file, err := os.Open(newFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to open file %s", newFilePath)
	}
	// Close the file later
	defer file.Close()

	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", attName)
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form file")
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to copy file %s", newFilePath)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("minorEdit")
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form field ")
	}

	_, err = fieldWriter.Write([]byte("true"))
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form field value")
	}
	// Populate other fields
	fieldWriter2, err := multiPartWriter.CreateFormField("comment")
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form field ")
	}

	_, err = fieldWriter2.Write([]byte(com))
	if err != nil {
		return nil, nil, fmt.Errorf("Confluence client: Failed to create Form field value")
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	req, err := http.NewRequest("POST", c.baseURL+path, &requestBody)

	if err != nil {
		return nil, nil, err
	}

	//fmt.Println(requestBody.String())

	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	req.Header.Set("X-Atlassian-Token", "nocheck")
	if c.usetoken {
		SetTokenAuth(req, c.password)
	} else {
		req.SetBasicAuth(c.username, c.password)
	}

	// Do the request
	response, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()
	if c.debug {
		log.Println("Response received, processing response...")
		log.Println("Response status code is", response.StatusCode)
		log.Println(response.Status)
	}
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return contents, nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 300 {
		log.Println("Bad response code received from server: ", response.Status)
		return contents, nil, fmt.Errorf("Bad response code received from server: %s ", response.Status)
	} else {
		json.Unmarshal(contents, retType)
	}
	return contents, retType, nil
}
