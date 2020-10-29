package client

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"reflect"
)

type GroupsType struct {
	Groups  []string `json:"groups,omitempty"  structs:"groups,omitempty"`
	Message string   `json:"message,omitempty" structs:"message,omitempty"`
	Status  string   `json:"status,omitempty" structs:"status,omitempty"`
}

type MembersType struct {
	Users  []map[string]string `json:"users,omitempty"  structs:"users,omitempty"`
	Status string              `json:"status,omitempty" structs:"status,omitempty"`
}

type AddGroupsType struct {
	Groups []string `json:"groups,omitempty"  structs:"groups,omitempty"`
}

type AddGroupsResponseType struct {
	GroupsAdded   []string `json:"groupsAdded,omitempty"  structs:"groupsAdded,omitempty"`
	GroupsSkipped []string `json:"groupsSkipped,omitempty"  structs:"groupsSkipped,omitempty"`
	Message       string   `json:"message,omitempty" structs:"message,omitempty"`
	Status        string   `json:"status,omitempty" structs:"status,omitempty"`
}
type AddMembersResponseType struct {
	UsersAdded   []string `json:"usersAdded,omitempty"  structs:"usersAdded,omitempty"`
	UsersSkipped []string `json:"usersSkipped,omitempty"  structs:"usersSkipped,omitempty"`
	Message      string   `json:"message,omitempty" structs:"message,omitempty"`
	Status       string   `json:"status,omitempty" structs:"status,omitempty"`
}

type AddUsersType struct {
	Users []string `json:"users,omitempty"  structs:"users,omitempty"`
}

type UsersType struct {
	Total      int `json:"total"`
	MaxResults int `json:"maxResults"`
	StartAt    int `json:"startAt"`
	Users      []struct {
		Name     string `json:"name"`
		FullName string `json:"fullName"`
		Key      string `json:"key"`
		Email    string `json:"email"`
	} `json:"users"`
}

func (c *ConfluenceClient) GetGroups() *GroupsType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/group/getGroups")
	groups := new(GroupsType)
	c.doRequest("GET", u, nil, &groups)
	return groups
}
type GetGroupMembersOptions struct {
	StartAt int `url:"startAt,omitempty"`
	MaxResults int `url:"maxResults,omitempty"`
	ShowBasicDetails  bool `url:"showBasicDetails,omitempty"`
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func (c *ConfluenceClient) GetGroupMembers(groupname string, options *GetGroupMembersOptions) (*UsersType, error) {
	var err error
	err = nil
	apiEndpoint := "/rest/extender/1.0/group/getUsers/" + groupname
	url :=""
	if (options != nil) {
		url, err = addOptions(apiEndpoint, options)
		if err != nil {
			return nil,  err
		}
	}
	members := new(UsersType)
	c.doRequest("GET", url, nil, &members)
	return members, err
}
/*
func (c *ConfluenceClient) GetGroupMembers(groupname string) *MembersType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/group/getUsers/" + groupname)
	members := new(MembersType)
	c.doRequest("GET", u, nil, &members)
	return members
}
*/
func (c *ConfluenceClient) AddGroups(groupnames []string) *AddGroupsResponseType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/group/addGroups")
	var payload = new(AddGroupsType)
	payload.Groups = append(payload.Groups, groupnames...)
	groups := new(AddGroupsResponseType)
	c.doRequest("POST", u, payload, &groups)
	//fmt.Println("res: " + string(res))
	return groups
}

func (c *ConfluenceClient) AddGroupMembers(groupname string, members []string) *AddMembersResponseType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/group/addUsers/" + groupname)
	var payload = new(AddUsersType)
	payload.Users = append(payload.Users, members...)
	response := new(AddMembersResponseType)
	c.doRequest("POST", u, payload, &response)
	//fmt.Println("res: " + string(res))
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

type PaginationOptions struct {
	// StartAt: The starting index of the returned projects. Base index: 0.
	StartAt int `url:"startAt,omitempty"`
	// MaxResults: The maximum number of projects to return per page. Default: 50.
	MaxResults int `url:"maxResults,omitempty"`
	// Expand: Expand specific sections in the returned issues
}

func (c *ConfluenceClient) GetUsers( options *PaginationOptions) *UsersType {
	var u string
	if options == nil {
		u = fmt.Sprintf("/rest/extender/1.0/group/getUsers")
	} else {
		u = fmt.Sprintf("/rest/extender/1.0/group/getUsers&startAt=%d&maxResults=%d", options.StartAt, options.MaxResults)
	}
	users := new(UsersType)
	c.doRequest("GET", u, nil, &users)
	return users
}

type GetAllGroupsWithAnyPermissionType struct {
	Total      int      `json:"total"`
	MaxResults int      `json:"maxResults"`
	Groups     []string `json:"groups"`
	StartAt    int      `json:"startAt"`
}

func (c *ConfluenceClient) GetAllGroupsWithAnyPermission( spacekey string, options *PaginationOptions) *GetAllGroupsWithAnyPermissionType {
	var u string
	if options == nil {
		u = fmt.Sprintf("/rest/extender/1.0/permission/space/%s/allGroupsWithAnyPermission", spacekey)
	} else {
		u = fmt.Sprintf("/rest/extender/1.0/permission/space/%s/allGroupsWithAnyPermission?startAt=%d&maxResults=%d", spacekey, options.StartAt, options.MaxResults)
	}
	groups := new(GetAllGroupsWithAnyPermissionType)
	c.doRequest("GET", u, nil, &groups)
	return groups
}

type GetPermissionsForSpaceType struct {
	Permissions []string `json:"permissions"`
	Name        string   `json:"name"`
	Key         string   `json:"key"`
}

func (c *ConfluenceClient) GetGroupPermissionsForSpace( spacekey, group string ) *GetPermissionsForSpaceType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/permission/group/%s/getPermissionsForSpace/space/%s", group, spacekey)
	permissions := new(GetPermissionsForSpaceType)
	c.doRequest("GET", u, nil, &permissions)
	return permissions
}

type PermissionsTypes []string

func (c *ConfluenceClient) GetPermissionTypes(  ) *PermissionsTypes {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/permission/space/permissionTypes")
	types := new(PermissionsTypes)
	c.doRequest("GET", u, nil, &types)
	return types
}

type GetAllUsersWithAnyPermissionType struct {
	Total      int      `json:"total"`
	MaxResults int      `json:"maxResults"`
	Users     []string  `json:"users"`
	StartAt    int      `json:"startAt"`
}

func (c *ConfluenceClient) GetAllUsersWithAnyPermission( spacekey string, options *PaginationOptions) *GetAllUsersWithAnyPermissionType {
	var u string
	if options == nil {
		u = fmt.Sprintf("/rest/extender/1.0/permission/space/%s/allUsersWithAnyPermission", spacekey)
	} else {
		u = fmt.Sprintf("/rest/extender/1.0/permission/space/%s/allUsersWithAnyPermission?startAt=%d&maxResults=%d", spacekey, options.StartAt, options.MaxResults)
	}
	users := new(GetAllUsersWithAnyPermissionType)
	c.doRequest("GET", u, nil, &users)
	return users
}

func (c *ConfluenceClient) GetUserPermissionsForSpace( spacekey, user string ) *GetPermissionsForSpaceType {
	var u string
	u = fmt.Sprintf("/rest/extender/1.0/permission/user/%s/getPermissionsForSpace/space/%s", user, spacekey)
	permissions := new(GetPermissionsForSpaceType)
	c.doRequest("GET", u, nil, &permissions)
	return permissions
}
