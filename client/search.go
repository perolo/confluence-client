package client

// Search searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) Search(cql string) (results *ConfluencePageSearch) {
	results = &ConfluencePageSearch{}
	req := "/rest/api/content/search?cql=" + cql + "&expand=body.view"
	//	fmt.Println(req)
	c.doRequest("GET", req, nil, results)
	return results
}
