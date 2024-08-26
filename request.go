package routego

import (
	"bytes"
	"encoding/json"
	"errors"
 
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)



func (r *Request) Cookie() map[string] http.Cookie{
	obj := make(map[string]http.Cookie, 0)
	
	for _, v := range r.HTTP.Cookies() {
		obj[v.Name] = *v 
	}
	
	return obj
}

func (r *Request) Header() map[string] []string{
	data := make(map[string] []string)

	for name, values := range r.HTTP.Header {
		data[name] = values
	}

	return data
}

func (r *Request) IP() string{
	forwarded := r.HTTP.Header.Get("X-Forwarded-For")
	
	if forwarded != "" {
	   ips := strings.Split(forwarded, ",")
	   return strings.TrimSpace(ips[0])
	}
   
	realIP := r.HTTP.Header.Get("X-Real-IP")
	if realIP != "" {
	   return realIP
	}
   
	ip, _, err := net.SplitHostPort(r.HTTP.RemoteAddr)
	if err != nil {
	   return ""
	}
	
	return ip
}

func (r *Request) DeviceInfo() (string, string){
    userAgent := r.HTTP.Header.Get("User-Agent")

    browser := "unknown Browser"
    os := "Unknown OS"

    if strings.Contains(userAgent, "Firefox/") {
        browser = "Mozilla Firefox"
    } else if strings.Contains(userAgent, "Chrome/") {
        browser = "Google Chrome"
    } else if strings.Contains(userAgent, "Safari/") && !strings.Contains(userAgent, "Chrome/") {
        browser = "Apple Safari"
    } else if strings.Contains(userAgent, "Edge/") {
        browser = "Microsoft Edge"
    } else if strings.Contains(userAgent, "MSIE") || strings.Contains(userAgent, "Trident/") {
        browser = "Microsoft Internet Explorer"
    }

    if strings.Contains(userAgent, "Windows NT") {
        os = "Windows"
    } else if strings.Contains(userAgent, "Mac OS X") {
        os = "Mac OS X"
    } else if strings.Contains(userAgent, "Linux") {
        os = "Linux"
    } else if strings.Contains(userAgent, "Android") {
        os = "Android"
    } else if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") {
        os = "iOS"
    }

	return browser, os
}

func (r *Request) JSON() (*interface{}, error){
	var data interface{}

	decoder := json.NewDecoder(r.HTTP.Body)
    err := decoder.Decode(&data)
    if err != nil {
        return nil, err
    }
    defer r.HTTP.Body.Close()

	return &data, nil
}


func (r *Request) Upload(to, name string) (bool, error) {
	mr, err := r.HTTP.MultipartReader()
	if err != nil {
		return false, err
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break 
		}
		if err != nil {
			return false, err
		}

		if part.FormName() == name {
			dst, err := os.Create(to)
			if err != nil {
				return false, err
			}
			defer dst.Close()

			_, err = io.Copy(dst, part)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (r *Request) UploadIfValid(name, to string, signature *MimeSignatureList) (bool, string, error) {
    mr, err := r.HTTP.MultipartReader()
    if err != nil {
        return false, "", err
    }

    for {
        part, err := mr.NextPart()
        if err == io.EOF {
            break
        }
        if err != nil {
            return false, "", err
        }

        if part.FormName() == name {
            file_name := strings.Split(part.FileName(), ".")
            file_ext := strings.ToLower(file_name[len(file_name)-1])
            file_name = nil
            header := make([]byte, 512)
            n, err := part.Read(header)
            if err != nil && err != io.EOF {
                return false, "", err
            }

            for _, v := range *signature {
                if bytes.HasPrefix(header[:n], v.Signature) {
                    for _, ext := range v.Extensions {
                        if file_ext == ext {
                            dst, err := os.Create(to + "." + ext)
                            if err != nil {
                                return false, "", err
                            }
                            defer dst.Close()

                            if _, err := dst.Write(header[:n]); err != nil {
                                return false, "", err
                            }

                            if _, err := io.Copy(dst, part); err != nil {
                                return false, "", err
                            }

                            return true, ext, nil
                        }
                    }
                    return false, "", errors.New("file extension does not match")
                }
            }

            return false, "", errors.New("signature does not match")
        }
    }

    return false, "", errors.New("file not found")
}