package client

import (
	"fmt"
)

type GroupsType struct {
	Groups  []string `json:"groups,omitempty"  structs:"groups,omitempty`
	Message string   `json:"message,omitempty" structs:"message,omitempty"`
	Status  string   `json:"status,omitempty" structs:"status,omitempty"`
}

type MembersType struct {
	Users  []map[string]string `json:"users,omitempty"  structs:"users,omitempty`
	Status string              `json:"status,omitempty" structs:"status,omitempty"`
}

type AddGroupsType struct {
	Groups []string `json:"groups,omitempty"  structs:"groups,omitempty`
}

type AddGroupsResponseType struct {
	GroupsAdded   []string `json:"groupsAdded,omitempty"  structs:"groupsAdded,omitempty`
	GroupsSkipped []string `json:"groupsSkipped,omitempty"  structs:"groupsSkipped,omitempty`
	Message       string   `json:"message,omitempty" structs:"message,omitempty"`
	Status        string   `json:"status,omitempty" structs:"status,omitempty"`
}
type AddMembersResponseType struct {
	UsersAdded   []string `json:"usersAdded,omitempty"  structs:"usersAdded,omitempty`
	UsersSkipped []string `json:"usersSkipped,omitempty"  structs:"usersSkipped,omitempty`
	Message      string   `json:"message,omitempty" structs:"message,omitempty"`
	Status       string   `json:"status,omitempty" structs:"status,omitempty"`
}

type AddUsersType struct {
	Users []string `json:"users,omitempty"  structs:"users,omitempty`
}

func (c *ConfluenceClient) GetGroups() *GroupsType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/group/getGroups")

	groups := new(GroupsType)
	res, _ := c.doRequest("GET", u, nil, &groups)

	fmt.Println("res: " + string(res))

	return groups
}

func (c *ConfluenceClient) GetGroupMembers(groupname string) *MembersType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/group/getUsers/" + groupname)

	members := new(MembersType)
	c.doRequest("GET", u, nil, &members)

	//fmt.Println("res: " + string(res))

	return members
}

func (c *ConfluenceClient) AddGroups(groupnames []string) *AddGroupsResponseType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/group/addGroups")

	var payload = new(AddGroupsType)
	payload.Groups = append(payload.Groups, groupnames...)

	groups := new(AddGroupsResponseType)
	res, _ := c.doRequest("POST", u, payload, &groups)

	fmt.Println("res: " + string(res))

	return groups
}

func (c *ConfluenceClient) AddGroupMembers(groupname string, members []string) *AddMembersResponseType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/group/addUsers/" + groupname)

	var payload = new(AddUsersType)
	payload.Users = append(payload.Users, members...)

	response := new(AddMembersResponseType)
	res, _ := c.doRequest("POST", u, payload, &response)

	fmt.Println("res: " + string(res))

	return response
}

type PermissionsType struct {
	Permissions []string `json:"permissions"`
}

type SpaceGroupPermissionType struct {
	Total      int `json:"total"`
	MaxResults int `json:"maxResults"`
	Spaces     []struct {
		Permissions []string `json:"permissions"`
		Name        string          `json:"name"`
		Key         string          `json:"key"`
	} `json:"spaces"`
	StartAt int `json:"startAt"`
}

type GroupOptions struct {
	Limit int
	Start int
}

///rest/extender/1.0/permission/group/{GROUP_NAME}/getAllSpacesWithPermissions
func (c *ConfluenceClient) GetAllSpacesForGroupPermissions(groupname string, options *GroupOptions) *SpaceGroupPermissionType {
	var req string
	if options == nil {
		req = "/rest/extender/1.0/permission/group/" + groupname + "/getAllSpacesWithPermissions?spacesAsArray=true"
	} else {
		req = fmt.Sprintf("/rest/extender/1.0/permission/group/%s/getAllSpacesWithPermissions?spacesAsArray=true&startAt=%d&maxResults=%d", groupname, options.Start, options.Limit)
	}

	spaces := new(SpaceGroupPermissionType)
	c.doRequest("GET", req, nil, &spaces)
	return spaces
}

type AddResponseType struct {
	Total   int           `json:"total"`
	Added   []string      `json:"added"`
	Skipped []interface{} `json:"skipped"`
}

///rest/extender/1.0/permission/group/{GROUP_NAME}/getAllSpacesWithPermissions
func (c *ConfluenceClient) AddSpacePermissionsForGroup(spacekey string, groupname string, pp []string) *AddResponseType {
	req := "/rest/extender/1.0/permission/space/" + spacekey + "/group/" + groupname + "/addSpacePermissions"
	var perm PermissionsType
	perm.Permissions = pp
	resp := new(AddResponseType)
	c.doRequest("PUT", req, perm, &resp)
	return resp
}
