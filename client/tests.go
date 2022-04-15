package client

type ConfluenceTestType struct {
	APIEndpoint string
	File        string
	Method      string
	Type        string
	TypeFile    string
}

const (
	TestSpace_GetSpaces_Moc_File_S = iota
	C1
	C2
)

var ConfluenceTest = []ConfluenceTestType{
	{APIEndpoint: "/rest/api/space", File: "mocks/spaces.json", Method: "GET", Type: "ConfluenceSpaceResult", TypeFile: "client/space-dtos.go"},
}
