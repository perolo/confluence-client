package client

import (
	"net/http"
	"fmt"
	"os"
	"bytes"
	"mime/multipart"
			"log"
	"io/ioutil"
	"io"
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


// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", uri, body)
}

func (c *ConfluenceClient) UpdateAttachment2(id string, attid string, filepath string) ([]byte){
//	response := &ConfluenceAttachmnetSearch{}
	path := fmt.Sprintf("/rest/api/content/%s/child/attachment/%s/data", id,attid)
//	c.doRequest("POST", path, bbody, response)
//	return  response

	//file=@myfile.txt" -F "minorEdit=true" -F "comment
	extraParams := map[string]string{
//		"file":       "data.json",
		"minorEdit":      "\"false\"",
		"comment": "\"Testing comment\"",
	}
//	htt, err := newfileUploadRequest(path, extraParams, "file", filepath )
	htt, err := newfileUploadRequest(path, extraParams, "file", "data.json" )
	if err != nil {
		log.Fatal(err)
	}

	htt.Header.Set("X-Atlassian-Token", "nocheck")

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
func (c *ConfluenceClient) UpdateAttachment(id string, attid string, filepath string) ([]byte){
	//	response := &ConfluenceAttachmnetSearch{}
	path := fmt.Sprintf("/rest/api/content/%s/child/attachment/%s/data", id,attid)


	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err)
	}
	// Close the file later
	defer file.Close()

	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file", "data.json")
	if err != nil {
		log.Fatalln(err)
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatalln(err)
	}
//	"minorEdit":      "\"false\"",
//		"comment": "\"Testing comment\"",

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("minorEdit")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = fieldWriter.Write([]byte("true"))
	if err != nil {
		log.Fatalln(err)
	}
	// Populate other fields
	fieldWriter2, err := multiPartWriter.CreateFormField("comment")
	if err != nil {
		log.Fatalln(err)
	}

//	_, err = fieldWriter2.Write([]byte("\"Test7\""))
	_, err = fieldWriter2.Write([]byte("Test7"))
	if err != nil {
		log.Fatalln(err)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	req, err := http.NewRequest("POST", c.baseURL+ path, &requestBody)
	//req, resp := c.doGetPage2("POST", path, &requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(requestBody.String())
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	req.Header.Set("X-Atlassian-Token", "nocheck")
	req.SetBasicAuth(c.username, c.password)

	// Do the request
	//client := &http.Client{}
	response, err := c.client.Do(req)
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
