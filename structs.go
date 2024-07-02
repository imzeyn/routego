package routego

import (
	"net/http"
	"regexp"
)

type Path struct {
	Name string
	AllowedMethods []string
	Handler	func(request Request) Reply
	Paths []Path
	Key string
}

type rePath struct{
	Path Path
	Paramas map[string]int
	Re *regexp.Regexp
	 
}

type newPattern struct{
	Name string
	Params map[string]int
	Re *regexp.Regexp
}

type Reply struct {
	Content any
	Json bool
	Status int
	Redirect string
	SetCookie []Cookie
	DelCookie []string
	SetHeader map[string]string
	ContnetType string
}

type Cookie struct{
	Name string
	Value string
	MaxAge int
	HttpOnnly bool
	SameSite string
	Secure bool
	Path string
}
 
type Server struct {
	Addr string
	Urls []Path
	PublicFolder string
	ServeFiles map[string]string
	urlmap map[string]Path
	reurls []rePath
	Handle ServerHandle
}

type ServerHandle struct{
	On404 func(request Request) Reply
	OnBefore func(request Request, key string) *Reply
}

type Request struct{
	HTTP http.Request
	Params map[string]string
}