package client

import (
	"fmt"
	"net/http"
)

type SpaceOptions struct {
	Limit  int    `url:"limit,omitempty"`
	Start  int    `url:"start,omitempty"`
	Label  string `url:"label,omitempty"`
	Type   string `url:"type,omitempty"`
	Status string `url:"status,omitempty"`
}

type ConfluenceSpaceResult struct {
	Results []SpaceType `json:"results,omitempty" structs:"results,omitempty"`
	Start   int         `json:"start,omitempty" structs:"start,omitempty"`
	Limit   int         `json:"limit,omitempty" structs:"limit,omitempty"`
	Size    int         `json:"size,omitempty"  structs:"size,omitempty"`
}
type SpaceType struct {
	Id         int               `json:"id,omitempty" structs:"id,omitempty"`
	Key        string            `json:"key,omitempty" structs:"key,omitempty"`
	Name       string            `json:"name,omitempty" structs:"name,omitempty"`
	Type       string            `json:"type,omitempty" structs:"type,omitempty"`
	Links      map[string]string `json:"_links,omitempty" structs:"_links,omitempty"`
	Start      string            `json:"start,omitempty" structs:"start,omitempty"`
	Expandable map[string]string `json:"_expandable,omitempty" structs:"_expandable,omitempty"`
}

type ConfluenceSpacePropertyResult struct {
	Results []string `json:"results,omitempty" structs:"results,omitempty"`
	Start   int      `json:"start,omitempty" structs:"start,omitempty"`
	Limit   int      `json:"limit,omitempty" structs:"limit,omitempty"`
	Size    int      `json:"size,omitempty"  structs:"size,omitempty"`
	Links   string   `json:"_links,omitempty" structs:"_links,omitempty"`
	Base    string   `json:"base,omitempty" structs:"base,omitempty"`
	Context string   `json:"context,omitempty" structs:"context,omitempty"`
}

type WatchResponseType struct {
	Watching bool `json:"watching,omitempty" structs:"watching,omitempty"`
}

//GetSpaces searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetSpaces(options *SpaceOptions) (results *ConfluenceSpaceResult) {
	var req string
	//var err error
	/*
		if options == nil {
			req = "/rest/api/space"
		} else {
			if options.Label == "" {
				req = fmt.Sprintf("/rest/api/space?start=%d&limit=%d", options.Start, options.Limit)
			} else {
				req = fmt.Sprintf("/rest/api/space?start=%d&limit=%d&label=%s", options.Start, options.Limit, options.Label)
			}

		}

	*/
	req, _ = addOptions("/rest/api/space", options)
	results = new(ConfluenceSpaceResult)
	c.doRequest("GET", req, nil, results)
	return results
}

//http://example.com/rest/experimental/space/TST/property?expand=space,version
func (c *ConfluenceClient) GetSpaceProperties(spacekey string) (results *ConfluenceSpacePropertyResult) {
	var req string
	req = fmt.Sprintf("/rest/api/space/%s/property?expand=space,version", spacekey)
	results = new(ConfluenceSpacePropertyResult)
	c.doRequest("GET", req, nil, results)
	return results
}

func (c *ConfluenceClient) AddWatcher(spaceKey string, user string) {
	var req string
	req = fmt.Sprintf("/rest/api/user/watch/space/%s?username=%s", spaceKey, user)

	var res *http.Response
	c.doRequest("POST", req, nil, res)
}

func (c *ConfluenceClient) GetWatcher(spaceKey string, user string) (results *WatchResponseType) {
	var req string
	req = fmt.Sprintf("/rest/api/user/watch/space/%s?username=%s", spaceKey, user)

	results = new(WatchResponseType)
	c.doRequest("GET", req, nil, results)
	return results
}
