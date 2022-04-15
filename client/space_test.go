package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	// testMux is the HTTP request multiplexer used with the test server.
	testMux *http.ServeMux

	// testClient is the Jira client being tested.
	testClient *ConfluenceClient

	// testServer is a test HTTP server used to provide mock API responses.
	testServer *httptest.Server
)

func TestGetSpaces(t *testing.T) {
	var config = ConfluenceConfig{}
	config.Username = "admin"
	config.Password = "admin"
	config.UseToken = false
	config.URL = "http://localhost:1990/confluence"
	config.Debug = true

	theClient := Client(&config)
	spopt := SpaceOptions{Start: 0, Limit: 20, Type: "global", Status: "current"}
	spaces, resp := theClient.GetSpaces(&spopt) //nolint:bodyclose
	defer CleanupH(resp)
	if resp.StatusCode != 200 {
		t.Errorf("Expected StatusCode 200, received: %v Spaces \n", resp.StatusCode)
	}
	if spaces.Size != 1 {
		t.Errorf("Expected 1 Space, received: %v Spaces \n", spaces.Size)
	}

}

func TestSpace_GetSpaces_Moc_Success(t *testing.T) {
	setup()
	defer teardown()
	//	testMux.HandleFunc("http://localhost:1990/confluence/rest/api/space", func(w http.ResponseWriter, r *http.Request) {
	testMux.HandleFunc("/rest/api/space", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/rest/api/space")

		_, err := fmt.Fprint(w, `{"results":[{"id":98305,"key":"ds","name":"Demonstration Space","type":"global","_links":{"webui":"/display/ds","self":"http://localhost:1990/confluence/rest/api/space/ds"},"_expandable":{"metadata":"","icon":"","description":"","retentionPolicy":"","homepage":"/rest/api/content/65551"}}],"start":0,"limit":25,"size":1,"_links":{"self":"http://localhost:1990/confluence/rest/api/space","base":"http://localhost:1990/confluence","context":"/confluence"}}`)
		if err != nil {
			t.Errorf("Error given: %s", err)
		}

	})

	spaces, resp := testClient.GetSpaces(nil) //nolint:bodyclose
	defer CleanupH(resp)
	if resp.StatusCode == 200 {
		if spaces == nil {
			t.Error("Expected Spaces. Spaces is nil")
		} else {
			if spaces.Size != 1 {
				t.Errorf("Expected 1 Space, received: %v Spaces \n", spaces.Size)
			}
		}
	} else {
		t.Error("Expected response 200.")
	}
}

func TestSpace_GetSpaces_Moc_File_Success(t *testing.T) {
	index := TestSpace_GetSpaces_Moc_File_S

	testAPIEndpoint := ConfluenceTest[index].APIEndpoint

	raw, err := ioutil.ReadFile("../" + ConfluenceTest[index].File)
	if err != nil {
		t.Error(err.Error())
	}

	setup()
	defer teardown()
	testMux.HandleFunc(testAPIEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, ConfluenceTest[index].Method)
		testRequestURL(t, r, testAPIEndpoint)

		_, err = fmt.Fprint(w, string(raw))
		if err != nil {
			t.Errorf("Error given: %s", err)
		}

	})

	spaces, resp := testClient.GetSpaces(nil) //nolint:bodyclose
	defer CleanupH(resp)

	if resp.StatusCode == 200 {

		if spaces == nil {
			t.Error("Expected Spaces. Spaces is nil")
		} else {
			if spaces.Size != 1 {
				t.Errorf("Expected 1 Space, received: %v Spaces \n", spaces.Size)
			}
		}
	} else {
		t.Error("Expected response 200.")
	}

}

// setup sets up a test HTTP server along with a jira.Client that is configured to talk to that test server.
// Tests should register handlers on mux which provide mock responses for the API method being tested.
func setup() {
	// Test server
	testMux = http.NewServeMux()
	testServer = httptest.NewServer(testMux)

	// jira client configured to use test server
	//testClient, _ = NewClient(nil, testServer.URL)
	var config = ConfluenceConfig{}
	config.Username = "admin"
	config.Password = "admin"
	config.UseToken = false
	//config.URL = "http://localhost:1990/confluence"
	config.URL = testServer.URL
	config.Debug = true

	testClient = Client(&config)

}

// teardown closes the test HTTP server.
func teardown() {
	//	testServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testRequestURL(t *testing.T, r *http.Request, want string) {
	if got := r.URL.String(); !strings.HasPrefix(got, want) {
		t.Errorf("Request URL: %v, want %v", got, want)
	}
}
