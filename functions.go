package routego

import (
	"fmt"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

var (
    doubleBracePattern = regexp.MustCompile(`{{.*?}}`)
    singleBracePattern = regexp.MustCompile(`{.*?}`)
	clearBrancePatter =  regexp.MustCompile(`[{}]`)
)

func urlArr(text string) []string{
	list := []string{}
	for _, v := range strings.Split(text, "/") {
		if strings.Trim(v, " ") != ""{
			list = append(list, v)
		}
	}
	return list
}

func goPattern(text string) newPattern {
	newStr := `^/`
	paramList := map[string]int{}
	
	for i, v := range urlArr(text) {
		if doubleBracePattern.MatchString(v){
			newStr += `?([\p{L}\p{N}\p{M}.@_-]*/)?`
			paramList[clearBrancePatter.ReplaceAllString(v, "")] = i
		}else if singleBracePattern.MatchString(v){
			newStr += `[\p{L}\p{N}\p{M}.@_-]+/`
			paramList[clearBrancePatter.ReplaceAllString(v, "")] = i
		}else{
			newStr += v + "/"
		}
	}

	newStr += `$`
 
	return newPattern{
        Name: newStr,
        Params: paramList,
		Re: regexp.MustCompile(newStr),
    }
}

func (s *Server) urlReslover(u Path, parent string){
	urlstr := parent + u.Name

	if !strings.HasSuffix(urlstr, "/"){
		urlstr += "/"
	} 
 
	if len(u.AllowedMethods) == 0{
		u.AllowedMethods = []string{"GET"}
	}

	copy_path := Path{
		urlstr,
		u.AllowedMethods,
		u.Handler,
		[]Path{},
		u.Key,
	}

	if singleBracePattern.MatchString(urlstr) {

		pattern := goPattern(copy_path.Name)

		copy_path.Name = pattern.Name
		
		s.reurls = append(s.reurls, rePath{
			copy_path, 
			pattern.Params,
			pattern.Re,
		})

	}else{
		s.urlmap[urlstr] = copy_path
	}

	for i := 0; i < len(u.Paths); i++ {
		s.urlReslover(u.Paths[i], urlstr)
	}

	u.Paths = nil
}

func checkMethod(method string, method_list []string) bool{
	for _, v := range method_list {
		if v == method{
			return true
		}
	}
	return false
}

func resloveReply(response Reply, w http.ResponseWriter) {
	
	status := 200
	 
	if(response.Status != 0){
		status = response.Status
	}

	if response.Redirect != ""{
		if response.Status == 0{
			status = 303
		}

		w.Header().Add("Location", response.Redirect)
		w.WriteHeader(status)
		return
	}
	
	if response.ContnetType == ""{
		response.ContnetType = "text/html; charset=UTF-8"
	}

	w.Header().Add("Content-Type", response.ContnetType)

	if response.Json{
		j, err := json.Marshal(response.Content)
		
		if(err != nil){
			panic(err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(j)
	}else if value, ok := response.Content.(string); ok{
		
		if response.SetCookie != nil{
			SetCookies(w, response.SetCookie)
		}

		if response.SetHeader != nil {
			setHeaders(w, response.SetHeader)
		}

		w.WriteHeader(status)
		w.Write([]byte(value)) 
	}else{
		panic("Reply: While the response was expected to be of string type, it was given an unknown type of response")
	}

}

func setHeaders(w http.ResponseWriter, headers map[string]string) {
	for k, v := range headers {
		w.Header().Add(k, v)
	}
}

func buildCookie(c Cookie) string{
	if(c.Path == ""){
		c.Path = "/"
	}

	value := c.Name + "=" + c.Value + ";path=" + c.Path + ";Max-Age=" + fmt.Sprintf("%d", c.MaxAge) + ";SameSite=" + c.SameSite + ";"
	if c.HttpOnnly{
		value += "HttpOnly;"
	}

	if(c.Secure){
		value += "Secure;"
	}

	return  value
}