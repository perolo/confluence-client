package client

import (
	"net/http"
	"fmt"
	"os"
	"bytes"
	"mime/multipart"
	"path/filepath"
	"io"
	"log"
	"io/ioutil"
	)

type PageOptions struct {
	// StartAt: The starting index of the returned projects. Base index: 0.
	Start int `url:"start,omitempty"`
	// MaxResults: The maximum number of projects to return per page. Default: 50.
	Limit int `url:"limit,omitempty"`
	// Expand: Expand specific sections in the returned issues
	Type string `url:"type,omitempty"`
}

//import "net/url"

//SearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetPageById(id string) (results *ConfluencePage) {
	results = &ConfluencePage{}
	c.doRequest("GET", "/rest/api/content/"+id+"?expand=body.view", nil, results)
	return results
}

func (c *ConfluenceClient) GetPages(space string, options *PageOptions) (results *ConfluencePages) {
	//path := fmt.Sprintf("rest/api/space/%s/content", space)
	//type=page&start=25&limit=3
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


func (c *ConfluenceClient) GetPage(url string) ([]byte,  *http.Response){
	contents, response := c.doGetPage("GET", url, nil)
	return contents, response
}

func (c *ConfluenceClient) GetPageAttachmentById(id string, name string) (results *ConfluenceAttachmnetSearch) {
	path := fmt.Sprintf("/rest/api/content/%s/child/attachment??filename=%s", id, name)

	results = &ConfluenceAttachmnetSearch{}
	c.doRequest("GET", path, nil, results)
	return results
}
/*
func (c *ConfluenceClient) UpdateAttachment(title, spaceKey, filepath string, bodyOnly, stripImgs bool, ID string, version, ancestor int64) {
	page := newPage(title, spaceKey)
	page.ID = ID
	page.Version = &ConfluencePageVersion{version}
	if ancestor > 0 {
		page.Ancestors = []ConfluencePageAncestor{
			ConfluencePageAncestor{ancestor},
		}
	}
	response := &ConfluencePage{}
	page.Body.Storage.Value = getBodyFromFile(filepath, bodyOnly, stripImgs)
	c.doRequest("PUT", "/rest/api/content/"+ID, page, response)
	//log.Println("ConfluencePage Object Response", response)
}
*/


func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile2(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		writer.WriteField2(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Content: %s\n", body)

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Atlassian-Token", "nocheck")
	return req, err
}


func (c *ConfluenceClient) UpdateAttachment(id string, attid string, filepath string) ([]byte){
//	response := &ConfluenceAttachmnetSearch{}
	path := fmt.Sprintf("/rest/api/content/%s/child/attachment/%s/data", id,attid)
//	c.doRequest("POST", path, bbody, response)
//	return  response

	//file=@myfile.txt" -F "minorEdit=true" -F "comment
	extraParams := map[string]string{
//		"file":       "data.json",
		"minorEdit":      "false",
		"comment": "Testing comment",
	}
	htt, err := newfileUploadRequest(path, extraParams, "file", filepath )
	if err != nil {
		log.Fatal(err)
	}

	response, err := c.client.Do(htt)
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
	/*
	if response.StatusCode < 200 || response.StatusCode > 300 {
		log.Println("Bad response code received from server: ", response.Status)
	} else {
		json.Unmarshal(contents, responseContainer)
	}
	return contents, response	*/
	return contents
}
