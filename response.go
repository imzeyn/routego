package routego

import (
	"encoding/json"
	"fmt"
)

func (response *Response) Cookie(cookie CookieObject) {
	if response.cookie == nil{
		response.cookie = make([]string, 0)
	}
	
	cookieVal := fmt.Sprintf("%s=%s; Max-Age=%v; SameSite=%s;", cookie.Name, cookie.Value, cookie.MaxAge, cookie.SameSite)

	if(cookie.Path == ""){
		cookieVal += " Path=/;"
	}else{
		cookieVal += " Path=" + cookie.Path + ";"
	}

	if(cookie.Domain != ""){
		cookieVal += " Domain=" + cookie.Domain + ";" 
	}

	if cookie.HttpOnly{
		cookieVal += " HttpOnly;"
	}

	if cookie.Secure{
		cookieVal += " Secure;"
	}

	
	response.cookie = append(response.cookie, cookieVal)
	cookieVal = ""
}

func (response *Response) Header(key, value string){
	if response.headers == nil{
		response.headers = make(map[string]string)
	}
	response.headers[key] = value
}

func (response *Response) Status(code int){
	response.status = code
}

func (response *Response) Write(w string){
	data := []byte(w)
	w = ""
	if response.headers["Content-Type"] == ""{
		response.Header("Content-Type", "text/html; charset=utf-8;")
	}

	response.data = &data
	response.done <- true
	<-response.writed
	data = nil
}

func (response *Response) JSON(w interface{}) error {
	data, e := json.Marshal(w)
	if e != nil{
		return e
	}
	
	w = nil
	response.Header("Content-Type", "application/json; charset=utf-8;")
	response.data = &data
	response.done <- true
	<-response.writed
	data = nil
	return nil
}
