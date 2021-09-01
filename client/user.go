package client

import (
	"fmt"
	"net/http"
)

/*
{
"name": "&lt; USER NAME &gt;",
"password": "&lt; USER PASSWORD &gt;",
"emailAddress": "&lt; USER EMAIL &gt;",
"displayName": "&lt; USER DISPLAY NAME &gt;",
"notification" : "&lt; true / false &gt;" //whether to sent notification to given email address or not
}
{
"name":"dummy",
"password":"1234",
"emailAddress":"per.olofsson@assaabloy.com",
"displayName":"dummy",
"notification":"false"
}
*/
/*
type	"known"
username	"perolo"
userKey	"ff80818160bb5dce0160c6c7d1c40006"
profilePicture	{…}
displayName	"Perolo"
_links	{…}
_expandable	{…}

*/

type UserType struct {
	Type           string `json:"type,omitempty"  structs:"type,omitempty`
	UserName       string `json:"username,omitempty" structs:"username,omitempty"`
	UserKey        string `json:"userKey,omitempty" structs:"userKey,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"  structs:"profilePicture,omitempty`
	DisplayName    string `json:"displayName,omitempty"  structs:"displayName,omitempty`
	Links          string `json:"_links,omitempty"  structs:"_links,omitempty`
	Expandable     string `json:"_expandable,omitempty"  structs:"_expandable,omitempty`
}

type UserDetailType struct {
	Business []struct {
		Location   string `json:"location"`
		Position   string `json:"position"`
		Department string `json:"department"`
	} `json:"business"`
	LastFailedLoginDate     int64  `json:"lastFailedLoginDate"`
	LastSuccessfulLoginDate int64  `json:"lastSuccessfulLoginDate"`
	FullName                string `json:"fullName"`
	Personal                []struct {
		Website string `json:"website"`
		Im      string `json:"im"`
		Phone   string `json:"phone"`
	} `json:"personal"`
	UpdatedDate                   int64  `json:"updatedDate"`
	LastFailedLoginDateString     string `json:"lastFailedLoginDateString"`
	CreatedDate                   int64  `json:"createdDate"`
	CreatedDateString             string `json:"createdDateString"`
	HasAccessToUseConfluence      bool   `json:"hasAccessToUseConfluence"`
	Name                          string `json:"name"`
	LastSuccessfulLoginDateString string `json:"lastSuccessfulLoginDateString"`
	Key                           string `json:"key"`
	Email                         string `json:"email"`
	UpdatedDateString             string `json:"updatedDateString"`
}

type UserCreateType struct {
	UserName     string `json:"name,omitempty"          structs:"name,omitempty"`
	Password     string `json:"password,omitempty"      structs:"password,omitempty`
	Email        string `json:"email,omitempty"         structs:"email,omitempty`
	DisplayName  string `json:"fullName,omitempty"      structs:"fullName,omitempty`
	Notification string `json:"notification,omitempty"  structs:"notification,omitempty`
}

type MessageType struct {
	Message string `json:"message"`
}
func (c *ConfluenceClient) GetUser(name string) (*UserType, *http.Response) {
	var u string
	u = fmt.Sprintf("/rest/api/user?username=" + name)

	user := new(UserType)
	_, res2 := c.doRequest("GET", u, nil, &user)

//	fmt.Println("res: " + string(res))

	return user, res2
}


func (c *ConfluenceClient) CreateUser(newUser UserCreateType) (*http.Response) {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/user/add")

	//	payload =
	user := new(UserCreateType)
	res, res2 := c.doRequest("PUT", u, newUser, &user)

	fmt.Println("res: " + string(res))

	return res2
}
func (c *ConfluenceClient) GetUserDetails(name string) (*UserDetailType, *http.Response) {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/user/getUserDetails/" + name)

	user := new(UserDetailType)
	_, res2 := c.doRequest("GET", u, nil, &user)

	//	fmt.Println("res: " + string(res))

	return user, res2
}

func (c *ConfluenceClient) DeactivateUser(name string) (*MessageType, *http.Response) {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/user/deactivate/" + name)

	user := new(MessageType)
	_, res2 := c.doRequest("POST", u, nil, &user)

	//	fmt.Println("res: " + string(res))

	return user, res2
}
