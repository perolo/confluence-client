package client

import (
	"fmt"
)

type SpaceOptions struct {
	Limit int
	Start int
}

type ConfluenceSpaceResult struct {
	Results []SpaceType `json:"results,omitempty" structs:"results,omitempty"`
	Start   int         `json:"start,omitempty" structs:"start,omitempty"`
	Limit   int         `json:"limit,omitempty" structs:"limit,omitempty"`
	Size    int         `json:"size,omitempty"  structs:"size,omitempty`
}
type SpaceType struct {
	Id         int    `json:"id,omitempty" structs:"id,omitempty"`
	Key        string `json:"key,omitempty" structs:"key,omitempty"`
	Name       string `json:"name,omitempty" structs:"name,omitempty"`
	Type       string `json:"type,omitempty" structs:"type,omitempty"`
	Links      string `json:"links,omitempty" structs:"links,omitempty"`
	Start      string `json:"start,omitempty" structs:"start,omitempty"`
	Expandable string `json:"_expandable,omitempty" structs:"_expandable,omitempty"`
}

//SearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetSpaces(options *SpaceOptions) (results *ConfluenceSpaceResult) {
	var req string
	if options == nil {
		req = "/rest/api/space"
	} else {
		req = fmt.Sprintf("/rest/api/space?start=%d&limit=%d", options.Start, options.Limit)
	}
	results = new (ConfluenceSpaceResult)
	c.doRequest("GET", req, nil, results)
	return results
}
