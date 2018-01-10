package client

import "net/http"

//import "net/url"

//SearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetPageById(id string) (results *ConfluencePage) {
	results = &ConfluencePage{}
	c.doRequest("GET", "/rest/api/content/"+id+"?expand=body.view", nil, results)
	return results
}


func (c *ConfluenceClient) GetPage(url string) ([]byte,  *http.Response){
	contents, response := c.doGetPage("GET", url, nil)
	return contents, response
}
