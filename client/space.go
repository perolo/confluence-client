package client

import (
	"fmt"
	"net/http"
)

type SpaceOptions struct {
	Limit    int    `url:"limit,omitempty"`
	Start    int    `url:"start,omitempty"`
	Label    string `url:"label,omitempty"`
	Type     string `url:"type,omitempty"`
	Status   string `url:"status,omitempty"`
	SpaceKey string `url:"spaceKey,omitempty"`
}

type ConfluenceSpaceResult struct {
	Results []SpaceType `json:"results,omitempty" structs:"results,omitempty"`
	Start   int         `json:"start,omitempty" structs:"start,omitempty"`
	Limit   int         `json:"limit,omitempty" structs:"limit,omitempty"`
	Size    int         `json:"size,omitempty"  structs:"size,omitempty"`
}
type SpaceType struct {
	ID         int               `json:"id,omitempty" structs:"id,omitempty"`
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

// GetSpaces searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetSpaces(options *SpaceOptions) (results *ConfluenceSpaceResult, resp *http.Response) {
	var req string
	req, _ = addOptions("/rest/api/space", options)
	results = new(ConfluenceSpaceResult)
	_, resp = c.doRequest("GET", req, nil, results)
	return results, resp
}

// GetSpaceProperties http://example.com/rest/experimental/space/TST/property?expand=space,version
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

// AddSpaceCategory /rest/extender/1.0/category/addSpaceCategory/space/{SPACE_KEY}/category/{CATEGORY_NAME}
func (c *ConfluenceClient) AddSpaceCategory(spaceKey string, category string) http.Response {
	var req string
	req = fmt.Sprintf("/rest/extender/1.0/category/addSpaceCategory/space/%s/category/%s", spaceKey, category)

	var res http.Response
	c.doRequest("PUT", req, nil, &res)
	return res
}

// RemoveSpaceCategory /rest/ui/1.0/space/{SPACE_KEY}/label/{LABEL_ID}
func (c *ConfluenceClient) RemoveSpaceCategory(spaceKey string, categoryid int) http.Response {
	var req string
	req = fmt.Sprintf("/rest/ui/1.0/space/%s/label/%v", spaceKey, categoryid)

	var res http.Response
	c.doRequest("DELETE", req, nil, &res)
	return res
}

type SpaceCategoriesResponseType struct {
	Name       string `json:"name"`
	Categories []struct {
		Name     string `json:"name"`
		ID       int    `json:"id"`
		NiceName string `json:"niceName"`
	} `json:"categories"`
	Key string `json:"key"`
}

// GetSpaceCategories {CONFLUENCE_URL}/rest/extender/1.0/category/getSpaceCategories/{SPACE_KEY}
func (c *ConfluenceClient) GetSpaceCategories(spaceKey string) (results *SpaceCategoriesResponseType) {
	var req string
	req = fmt.Sprintf("/rest/extender/1.0/category/getSpaceCategories/%s", spaceKey)

	results = new(SpaceCategoriesResponseType)
	c.doRequest("GET", req, nil, results)
	return results
}
