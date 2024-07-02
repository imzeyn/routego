package routego

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"mime"
)

func (s *Server) Listen(){
	s.urlmap = make(map[string]Path)
	

	for i := 0; i < len(s.Urls); i++ {
		s.urlReslover(s.Urls[i], "")
	}

	s.Urls = nil

 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if s.ServeFiles != nil {
			url := r.URL.Path
			if !strings.HasSuffix(url, "/") {
				url += "/"
			}
			if strings.HasPrefix(url, s.ServeFiles["url"]) {
				path := filepath.Join(s.ServeFiles["folder"], strings.Replace(url, s.ServeFiles["url"], "", 1))
				file, err := os.Open(path)
				if err == nil {
 
					fileInfo, err := file.Stat()
					if err != nil {
						http.Error(w, "File information could not be retrieved", http.StatusInternalServerError)
						return
					}
 
					mimeType := mime.TypeByExtension(filepath.Ext(path))
					if mimeType == "" {
						buffer := make([]byte, 512)
						_, err := file.Read(buffer)
						if err == nil {
							mimeType = http.DetectContentType(buffer)
						}
					}

					w.Header().Set("Content-Type", mimeType)
					file.Seek(0, 0)
					http.ServeContent(w, r, path, fileInfo.ModTime(), file)

					file.Close()
				} else {

					http.Error(w, "File not found", http.StatusNotFound)
					return
				}
				return
			}
		}
		s.handle(w, r)
	})

	http.ListenAndServe(s.Addr, nil)
}


func (s *Server) handle(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path

	if !strings.HasSuffix(url, "/"){
		url += "/"
	}
 
	page, res := s.urlmap[url]
	params := map[string]string{}
	
	if !res{
		arrUrl := urlArr(url)
		arrlen := len(arrUrl) - 1
		for i := 0; i < len(s.reurls); i++ {
			if s.reurls[i].Re.MatchString(url){
				for k,v  := range s.reurls[i].Paramas {
					if v > arrlen{
						break
					}else{
						params[k] = arrUrl[v]
					}
				}
				page = s.reurls[i].Path
				res = true
				break
			}
		}
	}

	if !res{
		if s.Handle.On404 != nil{
			r := s.Handle.On404(Request{HTTP: *r})
			if r.Status == 0{
				r.Status = 404
			}
			resloveReply(r, w)
			return
		}

		resloveReply(Reply{Content: "", Json: false, Status: 404}, w)
		return
	}

	if !checkMethod(r.Method, page.AllowedMethods){
		w.WriteHeader(405)
		return
	}

	if s.Handle.OnBefore != nil{
		before := s.Handle.OnBefore(Request{*r, params}, page.Key)
		if before != nil{
			resloveReply(*before, w)
			return
		}
	}

	resloveReply(page.Handler(Request{*r, params}), w)
}

 