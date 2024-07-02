package routego

import (
	"encoding/json"
	"io"
	"net/http"	 
)

func ReadBody(r *Request) ([]byte, error){
	data, err := io.ReadAll(r.HTTP.Body)
	return data, err
}

func JsonBody(r *Request) (any, error) {
	data, err := ReadBody(r)
	if err != nil {
		return nil, err
	}

	var dat any
	if err := json.Unmarshal(data, &dat); err != nil {
		return nil, err
	}

	return dat, nil
}

func SetCookies(w http.ResponseWriter, c []Cookie) {
	for _, v := range c {
		w.Header().Add("Set-Cookie", buildCookie(v))
	}
}

func GetCookies(r *Request) map[string] string{
	values := map[string]string{}
	for _, v := range r.HTTP.Cookies() {
		values[v.Name] = v.Value
	}
 	return values
}