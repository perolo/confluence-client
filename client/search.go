package client

//import "net/url"

//SearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) Search(cql string) (results *ConfluencePageSearch) {
	results = &ConfluencePageSearch{}
	//req :=
	c.doRequest("GET", "/rest/api/content/search?cql="+cql+"&expand=body.view", nil, results)
//	c.doRequest("GET", "/rest/api/content/search?cql="+cql, nil, results)
	return results
}
