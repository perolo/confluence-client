package client

import "fmt"

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
	UserName     string `json:"name,omitempty" structs:"name,omitempty"`
	Password     string `json:"password,omitempty"  structs:"password,omitempty`
	Email        string `json:"emailAddress,omitempty"  structs:"emailAddress,omitempty`
	DisplayName  string `json:"displayName,omitempty"  structs:"displayName,omitempty`
	Notification string `json:"notification,omitempty"  structs:"notification,omitempty`
}

func (c *ConfluenceClient) GetUser(name string) (*UserType) {
	var u string
	u = fmt.Sprintf("/rest/api/user?username=" + name)

	user := new(UserType)
	res, _ := c.doRequest("GET", u, nil, &user)

	fmt.Println("res: " + string(res))

	return user
}
func (c *ConfluenceClient) CreateUser(newUser UserCreateType) (*UserCreateType) {
	var u string
	u = fmt.Sprintf("/rest/api/user")

	//	payload =
	user := new(UserCreateType)
	res, _ := c.doRequest("POST", u, newUser, &user)

	fmt.Println("res: " + string(res))

	return user
}
