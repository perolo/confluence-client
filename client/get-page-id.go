package client

import (
	"net/http"
	"fmt"
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
