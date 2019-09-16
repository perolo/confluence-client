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

type UserCreateType struct {
	UserName     string `json:"name,omitempty"          structs:"name,omitempty"`
	Password     string `json:"password,omitempty"      structs:"password,omitempty`
	Email        string `json:"email,omitempty"         structs:"email,omitempty`
	DisplayName  string `json:"fullName,omitempty"      structs:"fullName,omitempty`
	Notification string `json:"notification,omitempty"  structs:"notification,omitempty`
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
