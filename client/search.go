package client

// Search searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) Search(cql string) (results *ConfluencePageSearch) {
	results = &ConfluencePageSearch{}
	req := "/rest/api/content/search?cql=" + cql + "&expand=body.view"
	//	fmt.Println(req)
	_, resp := c.doRequest("GET", req, nil, results) //nolint:bodyclose
	defer CleanupH(resp)
	return results
}
