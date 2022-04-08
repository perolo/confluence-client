package client

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
	"reflect"
)

type GroupsType struct {
	Total      int      `json:"total"`
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Groups     []string `json:"groups,omitempty"  structs:"groups,omitempty"`
	Message    string   `json:"message,omitempty" structs:"message,omitempty"`
	Status     string   `json:"status,omitempty" structs:"status,omitempty"`
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
type RemoveMembersResponseType struct {
	UsersRemoved []string `json:"usersRemoved,omitempty"  structs:"usersRemoved,omitempty"`
	UsersSkipped []string `json:"usersSkipped,omitempty"  structs:"usersSkipped,omitempty"`
	Message      string   `json:"message,omitempty" structs:"message,omitempty"`
	Status       string   `json:"status,omitempty" structs:"status,omitempty"`
}

type GroupUsersType struct {
	Users []string `json:"users,omitempty"  structs:"users,omitempty"`
}

type UsersType struct {
	Total      int `json:"total"`
	MaxResults int `json:"maxResults"`
	StartAt    int `json:"startAt"`
	Users      []struct {
		Business []struct {
			Location   string `json:"location"`
			Position   string `json:"position"`
			Department string `json:"department"`
		} `json:"business"`
		LastFailedLoginDate     interface{} `json:"lastFailedLoginDate"`
		LastSuccessfulLoginDate interface{} `json:"lastSuccessfulLoginDate"`
		FullName                string      `json:"fullName"`
		Personal                []struct {
			Website string `json:"website"`
			Im      string `json:"im"`
			Phone   string `json:"phone"`
		} `json:"personal"`
		UpdatedDate                   int64       `json:"updatedDate"`
		LastFailedLoginDateString     interface{} `json:"lastFailedLoginDateString"`
		CreatedDate                   int64       `json:"createdDate"`
		CreatedDateString             string      `json:"createdDateString"`
		Name                          string      `json:"name"`
		LastSuccessfulLoginDateString interface{} `json:"lastSuccessfulLoginDateString"`
		Key                           string      `json:"key"`
		Email                         string      `json:"email"`
		UpdatedDateString             string      `json:"updatedDateString"`
		HasAccessToUseConfluence      bool        `json:"hasAccessToUseConfluence"`
	} `json:"users"`
}

func (c *ConfluenceClient) GetGroups(options *GetGroupMembersOptions) (*GroupsType, error) {
	var err error
	err = nil
	apiEndpoint := "/rest/extender/1.0/group/getGroups"
	theURL := ""
	if options != nil {
		theURL, err = addOptions(apiEndpoint, options)
		if err != nil {
			return nil, err
		}
	} else {
		theURL = apiEndpoint
	}
	groups := new(GroupsType)
	c.doRequest("GET", theURL, nil, &groups)
	return groups, err
}

type GetGroupMembersOptions struct {
	StartAt             int  `url:"startAt,omitempty"`
	MaxResults          int  `url:"maxResults,omitempty"`
	ShowBasicDetails    bool `url:"showBasicDetails,omitempty"`
	ShowExtendedDetails bool `url:"showExtendedDetails,omitempty"`
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

func (c *ConfluenceClient) GetGroupMembers(groupname string, options *GetGroupMembersOptions) (*UsersType, *http.Response, error) {
	var err error
	err = nil
	apiEndpoint := "/rest/extender/1.0/group/getUsers/" + groupname
	theURL := ""
	if options != nil {
		theURL, err = addOptions(apiEndpoint, options)
		if err != nil {
			return nil, nil, err
		}
	} else {
		theURL = apiEndpoint
	}
	members := new(UsersType)
	_, resp := c.doRequest("GET", theURL, nil, &members)
	return members, resp, err
}

func (c *ConfluenceClient) AddGroups(groupnames []string) *AddGroupsResponseType {
	u := "/rest/extender/1.0/group/addGroups"
	var payload = new(AddGroupsType)
	payload.Groups = append(payload.Groups, groupnames...)
	groups := new(AddGroupsResponseType)
	c.doRequest("POST", u, payload, &groups)
	// fmt.Println("res: " + string(res))
	return groups
}

func (c *ConfluenceClient) AddGroupMembers(groupname string, members []string) *AddMembersResponseType {
	u := fmt.Sprintf("/rest/extender/1.0/group/addUsers/" + groupname)
	var payload = new(GroupUsersType)
	payload.Users = append(payload.Users, members...)
	response := new(AddMembersResponseType)
	c.doRequest("POST", u, payload, &response)
	// fmt.Println("res: " + string(res))
	return response
}

// RemoveGroupMembers {CONFLUENCE_URL}/rest/extender/1.0/group/removeUsers/{GROUP}
func (c *ConfluenceClient) RemoveGroupMembers(groupname string, members []string) *RemoveMembersResponseType {
	u := fmt.Sprintf("/rest/extender/1.0/group/removeUsers/" + groupname)
	var payload = new(GroupUsersType)
	payload.Users = append(payload.Users, members...)
	response := new(RemoveMembersResponseType)
	c.doRequest("POST", u, payload, &response)
	// fmt.Println("res: " + string(res))
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
		Name        string   `json:"name"`
		Key         string   `json:"key"`
	} `json:"spaces"`
	StartAt int `json:"startAt"`
}

type GroupOptions struct {
	Limit int
	Start int
}

// GetAllSpacesForGroupPermissions /rest/extender/1.0/permission/group/{GROUP_NAME}/getAllSpacesWithPermissions
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

// AddSpacePermissionsForGroup /rest/extender/1.0/permission/group/{GROUP_NAME}/getAllSpacesWithPermissions
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

func (c *ConfluenceClient) GetUsers(options *PaginationOptions) *UsersType {
	var u string
	if options == nil {
		u = "/rest/extender/1.0/group/getUsers"
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

func (c *ConfluenceClient) GetAllGroupsWithAnyPermission(spacekey string, options *PaginationOptions) *GetAllGroupsWithAnyPermissionType {
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

func (c *ConfluenceClient) GetGroupPermissionsForSpace(spacekey, group string) *GetPermissionsForSpaceType {
	u := fmt.Sprintf("/rest/extender/1.0/permission/group/%s/getPermissionsForSpace/space/%s", group, spacekey)
	permissions := new(GetPermissionsForSpaceType)
	c.doRequest("GET", u, nil, &permissions)
	return permissions
}

type PermissionsTypes []string

func (c *ConfluenceClient) GetPermissionTypes() *PermissionsTypes {
	u := "/rest/extender/1.0/permission/space/permissionTypes"
	types := new(PermissionsTypes)
	c.doRequest("GET", u, nil, &types)
	return types
}

type GetAllUsersWithAnyPermissionType struct {
	Total      int      `json:"total"`
	MaxResults int      `json:"maxResults"`
	Users      []string `json:"users"`
	StartAt    int      `json:"startAt"`
}

func (c *ConfluenceClient) GetAllUsersWithAnyPermission(spacekey string, options *PaginationOptions) (*GetAllUsersWithAnyPermissionType, *http.Response) {
	var u string
	if options == nil {
		u = fmt.Sprintf("/rest/extender/1.0/permission/space/%s/allUsersWithAnyPermission", spacekey)
	} else {
		u = fmt.Sprintf("/rest/extender/1.0/permission/space/%s/allUsersWithAnyPermission?startAt=%d&maxResults=%d", spacekey, options.StartAt, options.MaxResults)
	}
	users := new(GetAllUsersWithAnyPermissionType)
	_, resp := c.doRequest("GET", u, nil, &users)
	return users, resp
}

func (c *ConfluenceClient) GetUserPermissionsForSpace(spacekey, user string) (*GetPermissionsForSpaceType, *http.Response) {
	u := fmt.Sprintf("/rest/extender/1.0/permission/user/%s/getPermissionsForSpace/space/%s", user, spacekey)
	permissions := new(GetPermissionsForSpaceType)
	_, resp := c.doRequest("GET", u, nil, &permissions)
	return permissions, resp
}

type GetAllSpacesWithPermissionsType struct {
	Total      int `json:"total"`
	MaxResults int `json:"maxResults"`
	Spaces     []struct {
		Permissions []string `json:"permissions"`
		Name        string   `json:"name"`
		Key         string   `json:"key"`
	} `json:"spaces"`
	StartAt int `json:"startAt"`
}

func (c *ConfluenceClient) GetAllSpacesWithPermissions(user string, options *GetGroupMembersOptions) (*GetAllSpacesWithPermissionsType, *http.Response) {
	var u string
	if options == nil {
		u = fmt.Sprintf("/rest/extender/1.0/permission/user/%s/getAllSpacesWithPermissions?spacesAsArray=true", user)
	} else {
		u = fmt.Sprintf("/rest/extender/1.0/permission/user/%s/getAllSpacesWithPermissions?spacesAsArray=true&startAt=%d&maxResults=%d", user, options.StartAt, options.MaxResults)
	}
	users := new(GetAllSpacesWithPermissionsType)
	_, resp := c.doRequest("GET", u, nil, &users)
	return users, resp
}

type GetSpacePermissionAllActorsType struct {
	Permissions struct {
		Setpagepermissions struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"SETPAGEPERMISSIONS"`
		Removepage struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"REMOVEPAGE"`
		Editblog struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"EDITBLOG"`
		Removeowncontent struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"REMOVEOWNCONTENT"`
		Editspace struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"EDITSPACE"`
		Removemail struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"REMOVEMAIL"`
		Setspacepermissions struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"SETSPACEPERMISSIONS"`
		Viewspace struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"VIEWSPACE"`
		Removeblog struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"REMOVEBLOG"`
		Comment struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"COMMENT"`
		Createattachment struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"CREATEATTACHMENT"`
		Removeattachment struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"REMOVEATTACHMENT"`
		Removecomment struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"REMOVECOMMENT"`
		Exportspace struct {
			AnonymousAccess bool     `json:"anonymousAccess"`
			Groups          []string `json:"groups"`
			Users           []string `json:"users"`
		} `json:"EXPORTSPACE"`
	} `json:"permissions"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

// rest/extender/1.0/permission/space/DEMO/getSpacePermissionActors/ALL

func (c *ConfluenceClient) GetSpacePermissionAllActors(space string) (*GetSpacePermissionAllActorsType, *http.Response) {
	u := fmt.Sprintf("/rest/extender/1.0/permission/space/%s/getSpacePermissionActors/ALL", space)

	perm := new(GetSpacePermissionAllActorsType)
	_, resp := c.doRequest("GET", u, nil, &perm)
	return perm, resp
}
