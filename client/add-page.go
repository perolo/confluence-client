package client

import (
	"io/ioutil"
	"log"
	"strconv"

	"github.com/perolo/confluence-client/utility"
	"net/http"
)

// AddOrUpdatePage checks for an existing page then calls AddPage or UpdatePage depending on the result
func (c *ConfluenceClient) AddOrUpdatePage(options OperationOptions) bool {
	var retry = 0

	for retry < 3 {
		results := c.SearchPages(options.Title, options.SpaceKey)
		ancestorID := options.AncestorID
		if options.AncestorTitle != "" {
			ancestorResults := c.SearchPages(options.AncestorTitle, options.SpaceKey)
			if ancestorResults.Size < 1 {
				log.Fatal("Ancestor title not found!")
			} else {
				ancestorIDint, err := strconv.Atoi(ancestorResults.Results[0].ID)
				log.Println("Found ancestor ID", ancestorIDint)
				if err != nil {
					log.Fatal(err)
				}
				ancestorID = int64(ancestorIDint)
			}
		}
		var res = true
		if results.Size == 1 {
			log.Println("Page found, updating page...")
			item := results.Results[0]
			res = c.UpdatePage(options.Title, options.SpaceKey, options.Filepath, options.BodyOnly, options.StripImgs, item.ID, item.Version.Number+1, ancestorID)
		} else if results.Size == 0 {
			log.Println("Page not found, adding page...")
			c.AddPage(options.Title, options.SpaceKey, options.Filepath, options.BodyOnly, options.StripImgs, ancestorID)
		} else {
			log.Println("Multiple Pages found, retry")
			res = false
		}

		if res {
			return res
		} else {
			retry++
		}
	}
	return false
}

// AddPage adds a new page to the space with the given title
func (c *ConfluenceClient) AddPage(title, spaceKey, filepath string, bodyOnly, stripImgs bool, ancestor int64) {
	page := newPage(title, spaceKey)
	if ancestor > 0 {
		page.Ancestors = []ConfluencePageAncestor{
			ConfluencePageAncestor{ancestor},
		}
	}
	response := &ConfluencePage{}
	page.Body.Storage.Value = getBodyFromFile(filepath, bodyOnly, stripImgs)
	_, resp := c.doRequest("POST", "/rest/api/content/", page, response) //nolint:bodyclose
	defer CleanupH(resp)
	//log.Println("ConfluencePage Object Response", response)
}

func CleanupH(resp *http.Response) {
	// fmt.Println("Running Cleanup...")
	err := resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// UpdatePage adds a new page to the space with the given title
func (c *ConfluenceClient) UpdatePage(title, spaceKey, filepath string, bodyOnly, stripImgs bool, id string, version, ancestor int64) bool {
	page := newPage(title, spaceKey)
	page.ID = id
	page.Version = &ConfluencePageVersion{version}
	if ancestor > 0 {
		page.Ancestors = []ConfluencePageAncestor{
			ConfluencePageAncestor{ancestor},
		}
	}
	response := &ConfluencePage{}
	page.Body.Storage.Value = getBodyFromFile(filepath, bodyOnly, stripImgs)
	_, resp := c.doRequest("PUT", "/rest/api/content/"+id, page, response) //nolint:bodyclose
	defer CleanupH(resp)
	log.Println("ConfluencePage Object Response", response)
	return response.Status == "current"
}

func (c *ConfluenceClient) DeletePage(title string, spaceKey string) bool {
	results := c.SearchPages(title, spaceKey)
	if results.Size == 1 {
		log.Println("Page found, Deleting page...")
		var response *http.Response
		_, resp := c.doRequest("DELETE", "/rest/api/content/"+results.Results[0].ID, nil, response) //nolint:bodyclose
		defer CleanupH(resp)
		return true
	} else {
		return false
	}

}

func getBodyFromFile(filepath string, bodyOnly, stripImgs bool) string {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	if !bodyOnly {
		return string(buf)
	}
	return utility.StripHTML(buf, bodyOnly, stripImgs)
}
