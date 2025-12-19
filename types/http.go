package types

type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	PATCH   HTTPMethod = "PATCH"
	DELETE  HTTPMethod = "DELETE"
	OPTIONS HTTPMethod = "OPTIONS"
	HEAD    HTTPMethod = "HEAD"
)

type HTTPQueryParam struct {
	Key   string
	Value string
}

type HTTPRequest struct {
	Path        string
	Method      string
	Query       []HTTPQueryParam
	Params      []HTTPQueryParam
	QueryString string
	IP          string
	URL         string
}
