package client

import "net/url"

//SearchPages searche	s for pages in the space that meet the specified criteria
func (c *ConfluenceClient) SearchPages(title, spaceKey string) (results *ConfluencePageSearch) {
	results = &ConfluencePageSearch{}
	c.doRequest("GET", "/rest/api/content?title="+url.QueryEscape(title)+"&spaceKey="+url.QueryEscape(spaceKey)+"&expand=version,body.view", nil, results)
	return results
}
