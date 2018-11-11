package common

type Route struct {
	Method  string
	URL     string
	Handler HandlerFunc
}
