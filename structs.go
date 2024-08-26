package routego

import (
	"net/http"
)

type RouteHandler func(request *Request, response *Response)

type Path struct {
	ID			   string
	URL            interface{}
	stringURL	   string
	Handler        RouteHandler
	AllowedMethods []string
	Include        []Path
	Additional	   interface{}
	paramList	   map[string]int

}

type urlmap struct{
	fixed	map[string]Path
	re		[]Path
}

type temporaryURL struct{
	path 	Path
	parent	string
}

type Server struct {
	Addr         string
	UrlPatterns  []Path
	clearPattern urlmap
	tempURL		 map[string] temporaryURL
	NotFound	 RouteHandler
	NotAllowed 	 RouteHandler
	Middleware	 []RouteHandler
}

type Response struct {
	data			  *[]byte
	writed 			chan bool
	done    		chan bool
	headers	map[string]string
	cookie  		 []string
	status  			  int
	WriteMiddleware	chan bool
}

type Request struct {
	HTTP   	   *http.Request
	Params 	   *map[string]string
	PathExist   bool
	PathID      string
	Additional	*interface{}
}

type CookieObject struct{
	Name		string
	Value		string
	Path		string
	Domain		string
	SameSite	string

	HttpOnly	bool
	Secure		bool

	MaxAge		int
}

type MimeSignature struct {
	Type      string
	Signature []byte
	Category  MimeCategory
	Extensions []string
}

type MimeSignatureList []MimeSignature
type MimeCategory string