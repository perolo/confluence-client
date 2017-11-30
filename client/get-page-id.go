package client

//import "net/url"

//SearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetPageById(id string) (results *ConfluencePage) {
	results = &ConfluencePage{}
	//req :=
	c.doRequest("GET", "/rest/api/content/"+id+"?expand=body.view", nil, results)
	return results
}
