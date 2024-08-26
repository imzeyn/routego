package routego

import (
	"net/http"
	"regexp"
	"strings"
)

func (s *Server) Listen(){
	s.clearPattern.fixed = map[string]Path{}
	s.clearPattern.re = make([]Path, 0)

	s.tempURL = map[string]temporaryURL{}

	s.resolvePaths(&s.UrlPatterns, "")
	s.UrlPatterns = nil
	s.buildPath()
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)
	http.ListenAndServe(s.Addr, mux)
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request){
	path, ok := s.getPath(r)
		
	response := Response{
		done: make(chan bool),
		writed: make(chan bool),
		WriteMiddleware: make(chan bool),
	}

	defer close(response.done)
	defer close(response.writed)
	defer close(response.WriteMiddleware)

	params := make(map[string] string)

	if ok {
		if len(path.paramList) > 0{
			paramURL := make([]string, 0)
			for _, v := range  strings.Split(s.clearURL(r.URL.Path), "/") {
				if strings.TrimSpace(v) != ""{
					paramURL = append(paramURL, v)
				}
			}
			paramLen := len(paramURL)
			for k, v := range path.paramList {
				if paramLen > v - 1{
					params[k] = paramURL[v - 1]
				}
			}
		}
	}
	
	request := Request{
		HTTP: r,
		Params: &params,
		PathExist: ok,
	}

	if ok{
		request.PathID = path.ID
		request.Additional = &path.Additional

	}else{
		request.PathID = ""
	}


	isWriteMiddleware := false
	for _, v := range s.Middleware {
		go v(&request, &response)
	 
		if <-response.WriteMiddleware{
			isWriteMiddleware = true
			break
		}
	}

	isAllowedMethod := isAllowedMethod(r.Method, path.AllowedMethods)
	if !isWriteMiddleware{
		if !isAllowedMethod{
			if s.NotAllowed != nil{
				go s.NotAllowed(&request, &response)
			}else{
				go s.notAllowed(&request, &response)
			}
		}else if ok{
			if(path.Handler == nil){
				if(s.NotFound != nil){
					go s.NotFound(&request, &response)
				}else{
					go s.notFound(&request, &response)
				}
			}else{
				go path.Handler(&request, &response)
			}
		}else if(s.NotFound != nil){
			go s.NotFound(&request, &response)
		}else{
			go s.notFound(&request, &response)
		}
	}
	
	<-response.done

	for k, v := range response.headers{
		w.Header().Add(k, v)
	}

	for _, v := range response.cookie {
		w.Header().Add("Set-Cookie", v)
	}

	if response.status == 0{
		if isAllowedMethod{
			response.status = 200
		}else{
			response.status = 403
		}
	}

	w.WriteHeader(response.status)
	w.Write(*response.data)

	response.data, response.headers, response.cookie = nil, nil, nil
	//w, r, params = nil, nil, nil

	response.writed <- true
}

func (s *Server) resolvePaths(patterns *[]Path, parent string){
	for _, v := range *patterns {
		key := randomString()
		if v.AllowedMethods == nil{
			v.AllowedMethods = []string{"GET"}
		}
		
		s.tempURL[key] = temporaryURL{
			parent: parent,
			path: Path{
				ID: v.ID,
				URL: v.URL,
				stringURL: v.URL.(string),
				AllowedMethods: v.AllowedMethods,
				Additional: v.Additional,
				Handler: v.Handler,
			},
		}
		if v.Include != nil{
			s.resolvePaths(&v.Include, key)
		}
	}
}

func (s *Server) buildPath(){
	for key := range s.tempURL {
		url, isre, params := s.reURL(s.clearURL(s.createURL(key)))
		
		id := s.tempURL[key].path.ID
		handler := s.tempURL[key].path.Handler
		additional := s.tempURL[key].path.Additional
		allowed_methods := s.tempURL[key].path.AllowedMethods
		 
		if !isre{
			if val, ok := s.clearPattern.fixed[url.(string)]; ok && handler != nil && val.Handler != nil{
				panic(url.(string) + " This URL has been defined more than once, there are multiple handlers for this URL")
			}else if handler == nil && val.Handler != nil{
				handler = val.Handler
				additional = val.Additional
				id = val.ID
				allowed_methods = val.AllowedMethods
			}
		}
		path := Path{
			ID: id,
			URL: url,
			Handler: handler,
			AllowedMethods: allowed_methods,
			Additional: additional,
			paramList: params,
			stringURL: s.tempURL[key].path.stringURL,
		}
		
		if isre{
			if val, ok, i := s.isRePathExists(*url.(*regexp.Regexp)); ok{
				if path.Handler != nil && (*val).Handler != nil{
					urlname := s.tempURL[key].path.URL.(string)
					if urlname != val.stringURL{
						urlname = val.stringURL
					}
						
					if urlname == "/" || urlname == ""{
						if strings.ReplaceAll(val.stringURL, "/", "") == ""{
							urlname = s.tempURL[key].path.URL.(string)
						}else{
							urlname = val.stringURL
						}
					}

					panic(urlname + " This URL has been defined more than once, there are multiple handlers for this URL")					
				}else if path.Handler == nil && (*val).Handler != nil{
					continue
				}else if path.Handler != nil && (*val).Handler == nil{
					
					s.clearPattern.re[i].Handler = path.Handler
					s.clearPattern.re[i].Additional = path.Additional
					s.clearPattern.re[i].AllowedMethods = path.AllowedMethods
					s.clearPattern.re[i].ID = path.ID
					
					
				}
			}else{

				s.clearPattern.re = append(s.clearPattern.re, path)
			}

			continue
		}

		s.clearPattern.fixed[url.(string)] = path
	}

	s.tempURL = nil
}

func (s *Server) isRePathExists(url regexp.Regexp) (*Path, bool, int){
	 
	for i, v := range s.clearPattern.re {
		if(v.URL.(*regexp.Regexp).String() == url.String()){
			return &v, true, i
		}
	}
	return nil, false, -1
}

func (s *Server) createURL(key string) string{
	var url string

	if entry, exists := s.tempURL[key]; exists {
		url = entry.path.URL.(string)
		if entry.parent == ""{
			return url
		}
		parentURL := s.createURL(entry.parent)
		if parentURL != ""{
			return parentURL + url
		}
	}

	return url
}

func (s *Server) reURL(url string) (interface{}, bool, map[string] int){
	
	if !singleBracePattern.MatchString(url){
		return url, false, nil
	}

	newURL := ""
	paramList := map[string]int{}

	for i, v := range strings.Split(url, "/") {
		paramArr := singleBracePattern.FindStringSubmatch(v)
		if len(paramArr) > 0{
			paramName := paramArr[0]
			paramName = strings.ReplaceAll(paramName, "{", "")
			paramName = strings.ReplaceAll(paramName, "}", "")
			if doubleBracePattern.MatchString(v){
				newURL += doubleBracePattern.ReplaceAllString(v, `?([\p{L}\p{N}\p{M}.@_-]*)?/`)
				paramList[paramName] = i			
			}else if singleBracePattern.MatchString(v){
				newURL += singleBracePattern.ReplaceAllString(v, `[\p{L}\p{N}\p{M}.@_-]+/`)
				paramList[paramName] = i
			}else{
				newURL += v + "/"
			}
		}else{
			newURL += v + "/"
		}
	}

	newURL += ""

	return regexp.MustCompile( "^" + s.clearURL(newURL) + "$" ), true, paramList
}


func (s *Server) clearURL(url string) string{
	newURL := "/"
	for _, v := range strings.Split(url, "/") {
		if strings.TrimSpace(v) != ""{
			newURL += v + "/"
		}
	}

	return newURL
}

func (s *Server) getPath(request *http.Request) (*Path, bool){
	url := s.clearURL(request.URL.Path)
	
	path, ok := s.clearPattern.fixed[url]; 

	if !ok{
		for _, v := range s.clearPattern.re {
			re := v.URL.(*regexp.Regexp)
			 	
			if re.MatchString(url){
				ok = true
				path = v
			}

		}
	}

	return &path, ok
}

func isAllowedMethod(method string, method_list []string) bool{
	for _, v := range method_list{
		if method == v{
			return true
		}
	}
	return false
}

func (s *Server) notFound(_ *Request, response *Response){
	response.Status(404)
	response.Write("")
}

func (s *Server) notAllowed(_ *Request, response *Response){
	response.Status(403)
	response.Write("")
}

