# RouteGo Library

RouteGo is a library used to route and process HTTP requests in the Go language. This library provides a routing mechanism that matches HTTP requests against specific URL patterns and forwards them to the appropriate handlers.

# Features

 - URLs without/with parameters
 - Getting mandatory parameters or optional parameters
 - Customizable 404
 - Any action before accessing a page
 - Serve Files
 - Managing cookies
 - Managing headers
 - Quickly provide json responses
 - Reading json data in request body


# Setup

Create a go.mod file in your project:

    go mod init myapp
Include the routego library in your project:

    go get github.com/imzeyn/routego

# Configuration and Usage

To use the routego library, you must first determine the server address and port, then specify the URL list of your website.

## Simple Operation
```go
package main
import "github.com/imzeyn/routego"

func main() {
	myUrls := [] routego.Path {{
		Name: "/",
		Handler: home,
	}} // Your URL list here
	server := routego.Server {
		Addr: "0.0.0.0:8080", //Or :8080
		Urls: myUrls,
	}
	server.Listen() // Open the web server
}

func home(request routego.Request) routego.Reply {
	return routego.Reply {
		Content: "Wow! Its home page.",
	}
}
```

## Creating a URL
We use `routego.Path{}` to create the url. You need to add the generated URLs into `routego.Server.Urls`. If the URL matches, the `routego.Path.Handler` method works.


#### Example of the usage
```go
myurl := routego.Path{
	Name: "/about", // Your URL Address,
	AllowedMethods: []string{"GET"}, // Specify here which methods you want to allow. Default GET.
	Handler: func(request  routego.Request) routego.Reply {
		return  routego.Reply{Content: "About Page"}
	},
	Key: "about",
}
```
Matches the following URL:
	
 	1. /about
  	2. /about/

### Specifying a required parameter
```go
myurl  :=  routego.Path{
	Name: "/user/{id}", // Your URL Address,
	AllowedMethods: []string{"GET"},
	Handler: myFunc,
	Key: "user-with-param",
}
```
This URL matches the following lists:
 1. GET /user/@blabla123
 2. GET /user/@blabla123/

### Specifying an optional parameter
```go
myurl  :=  routego.Path{
	Name: "/users/{{id}}", // Your URL Address,
	AllowedMethods: []string{"GET"},
	Handler: myFunc,
	Key: "user-with-param-optional",
}
```
This URL matches the following lists:
1. GET /users
2. GET /users/
3. GET /users/@blabla123
4. GET /users/@blabla123/


### Specify a sub-URL
```go
myurl := routego.Path{
	Name: "/foo", // Your URL Address,
	AllowedMethods: []string{"GET"}, // Specify here which methods you want to allow. Default GET.
	Handler: myFunc,
	Key: "user-with-param",
	Paths: []routego.Path{{
		Name: "/bazz",
		AllowedMethods: []string{"GET", "POST"},
		Handler: bazz,
		Key: "bazz",
		Paths: [] routego.Path{{
			Name: "/bar",
			AllowedMethods: []string{"DELETE"},
			Handler: bar,
			Key: "bar",
        		Paths: [] routego.Path{{...}},
		}},
	}},
}
```
Then you need to get a URL like below:

 1. /foo
 2. /foo/bazz
 3. /foo/bazz/bar

## Give a response
We use `routego.Reply` to easily respond as a string or json.

`routego.Path.Handler` takes a function and this function must return `routego.Reply`. This function takes only 1 parameter and this parameter should be `routegor.Request`.
### Example of code
```go
package main
import "github.com/imzeyn/routego"
func main() {
	myurl := routego.Path{
		Name: "/",
		Handler: home,
	}
	myUrls := []routego.Path{myurl} // Your URL list here
	server := routego.Server{
		Addr: "0.0.0.0:8080", //Or :8080
		Urls: myUrls,
	}
	server.Listen() // Open the web server
}


func home(request routego.Request) routego.Reply {
	return routego.Reply{Content: "Wow! Its home page."}
}
```
## Give a JSON response
To return a JSON response, just say `routego.Reply.Json = true`.
### Example of code
```go
func json_response(request routego.Request) routego.Reply {
	json := map[string] any{
		"key": "value",
		"list": [3]int{1, 2, 3},
	}
	return routego.Reply{Content: json, Json: true}
}
```
## Specify a status code
To specify a status code, simply write the status code in `int` type in the `routego.Reply.Status` section.
This value is `200` by default, 303 for redirected pages and `404` when there are no pages. If a request is made due to an invalid method, it is automatically set to `405`.
### Example of code
```go
func server_error(request routego.Request) routego.Reply {
	return routego.Reply{Content: "...", Status: 500}
}
```

## Adding a cookies
You must provide the cookies you want to add in a list as `[]routego.Cookie`. If you then give this list to `routego.Reply` the cookie will be added to the user's browser.
### Example of code
```go
func add_cookie(request routego.Request) routego.Reply {
	myCookies := []routego.Cookie{
		{
			Name: "authtoken",
			Value: "...",
			MaxAge: 3600, //Seconds, 3600 seconds is 1 hour
			HttpOnnly: true, //Default is false. Optional
			SameSite: "lax", //Optional
			Secure: true, // Optional
			Path: "/mypath", // Optional. Default is /
		},
		{
			Name: "authtoken123",
			Value: "...",
			MaxAge: 3600, //Seconds, 3600 seconds is 1 hour
		},
	}
	return routego.Reply{
		Content: "...",
		Json: true, // Or false default is false
		SetCookie: myCookies,
	}
}
```
## Deleting a cookies
You must specify the cookies you want to delete as `[]string`. If you then provide this list to `routego.Reply`, the cookie will be deleted from the user's browser. The content of this list should be the names of the cookies.
### Example of code
```go
func delete_cookie(request routego.Request) routego.Reply {
	return routego.Reply{
		Content: "...",
		Json: true, // Or false default is false
		DelCookie: []string{"authtoken"},
	}
}
```
## Adding a header
You must specify the header information you want to add as `map[string]string`. If you then provide this map to `routego.Reply`. Titles will be added.
### Example of code
```go
func add_header(request routego.Request) routego.Reply {
	return routego.Reply{
		Content: "...",
		SetHeader: map[string]string{
			"X-Foo": "Bar",
			"X-Bazz": "Foo",
		},
	}
}
```
## Change content type
```go
func text_file(request routego.Request) routego.Reply {
	return routego.Reply{
		Content: "Im a text",
		ContnetType: "text/plain", //Default is text/html; charset=UTF-8
	}
}
```

## Handling requests
We access it thanks to the `routego.Request` parameter that `routego.Path.Handler` receives when processing requests. `routego.Request` has two keys, one is `HTTP` and the other is `Params`. `routego.Request.HTTP` this parameter defaults to `http.Request` in the `net/http` library. `Params`, if you have specified a parameter for the URL, returns it to you as `map[string]string`.
### Example codes
#### Reading cookies
```go
func read_cookies(request routego.Request) routego.Reply {
	cookie_map := routego.GetCookies(&request) // Returns as map[string]string
	value, err := cookie_map["my_cookie_name"]
	fmt.Println(value, err)
	return routego.Reply{Content: "..."}
}
```

#### Reading body
```go
func read_body(request routego.Request) routego.Reply {
	data, err := routego.ReadBody(&request) // Returns as ([]byte, error)
	fmt.Println(data, err)
	return routego.Reply{Content: "..."}
}
```

#### Read body as JSON
```go
func read_body_as_json(request routego.Request) routego.Reply {
	data, err := routego.JsonBody(&request) // Returns as (any, error)
	fmt.Println(data, err)
	return routego.Reply{Content: "..."}
}
```
#### Request method, headers etc. getting things
`request.HTTP` is identical to `http.Request` in `net/http` package
```go
func handle(request routego.Request) routego.Reply {
	fmt.Println(request.HTTP.Method)
	fmt.Println(request.HTTP.RequestURI)
	fmt.Println(request.HTTP.URL.Path)
	fmt.Println(request.HTTP.Header.Get("Accept-Language"))
	fmt.Print(request.HTTP.FormValue("foo")) // Get the multipart data
	return routego.Reply{Content: "..."}
}
```
#### Getting parameter values ​​in parameterized URL
```go
package main

import "github.com/imzeyn/routego"

func main() {
	myUrls := []routego.Path{
		{
			Name: "/user/{id}",
			Handler: get_user,
		},
		{
			Name: "/blog/{{blogID}}",
			Handler: blog_list_or_content,
		},
	} // Your URL list here
	server := routego.Server{
		Addr: "0.0.0.0:8080", //Or :8080
		Urls: myUrls,
	}
	server.Listen() // Open the web server
}


func get_user(request routego.Request) routego.Reply {
	id := request.Params["id"] // Returns as string
	return routego.Reply{Content: "User ID is " + id}
}

func blog_list_or_content(request routego.Request) routego.Reply{
	content := "You are showing blog list."
	if blogID, ok := request.Params["blogID"]; ok{
		content = "Blog content for " + blogID
	}
	return routego.Reply{Content: content}
}
```
## Custom 404 page
```go
server := routego.Server{
	Addr: "0.0.0.0:8080", //Or :8080
	Urls: myUrls,
	Handle: routego.ServerHandle{
		On404: func(request routego.Request) routego.Reply {
			return routego.Reply{Content: "Page is not found"}
		},
	},
}

server.Listen() // Open the web server
```
## Performing an action before accessing the page
When there is any match in your URL list and you want to perform an action before `Handler`. We use `OnBefore`.

```go
server := routego.Server{
	Addr: "0.0.0.0:8080", //Or :8080
	Urls: myUrls,
	Handle: routego.ServerHandle{
		On404: func(request routego.Request) routego.Reply {
			return routego.Reply{Content: "Page is not found"}
		},
		OnBefore: func(request routego.Request, key string) *routego.Reply {
			token := request.HTTP.Header.Get("X-Token") != "..." ;
			if key == "api" && token{
				response := map[string] string{
					"error": "Invalid API token.",
				}
				return &routego.Reply{Content: response, Json: true}
			}
			fmt.Println("URL", request.HTTP.URL)
			return nil
		},
	},
}

server.Listen() // Open the web server
```
If `OnBefore` returns `nil` then `routego.Path.Handler` in the URL will work. If `OnBefore` returns a `routego.Reply` then `routego.Path.Handler` will not work.
The second parameter taken here, `key`, is the value you define in `routego.Path.Key`. Whether you define it uniquely or not, the purpose here is to understand which page you visited with `OnBefore`.
```go
myUrls := []routego.Path{
	{
		Name: "/user/{id}",
		Handler: get_user,
		Key: "user",
	},
	{
		Name: "/blog/{{blogID}}",
		Handler: blog_list_or_content,
		Key: "blog",
	},
}  
```
## Serve Files
```go
package main

import (
    "goroute"
)

func main() {
    server := goroute.Server{
        Addr: ":8080",
        ServeFiles: map[string]string{
            "url":    "/static/",
            "folder": "./public",
        },
    }

    server.Listen()
}
```
In this example, the server serves requests under the `/static/` URL path from the local `./public` directory. For example, a request to `http://localhost:8080/static/css/styles.css` in the browser causes the server to serve `./public/css/styles.css`.
## Note
Although the ServeFiles feature is suitable for small projects and development environments, it is not recommended in a production environment. Instead, it will be safer and more performant to use professional web servers such as Apache2 and Nginx to serve static files. These servers perform better in high-traffic environments and offer more configuration options for security.
