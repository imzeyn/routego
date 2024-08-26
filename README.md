# RouteGo Library

RouteGo is a library used to route and process HTTP requests in the Go language. This library provides a routing mechanism that matches HTTP requests against specific URL patterns and forwards them to the appropriate handlers.

# Features

- **URL Handling**: 
  - Support for URLs with or without parameters.
  
- **Parameter Management**: 
  - Ability to extract mandatory and optional parameters from requests.

- **Customizable Error Handling**: 
  - Configurable responses for 404 Not Found and 403 Forbidden errors.

- **Middleware Support**: 
  - Integrate and manage multiple middleware functions for request processing.

- **File Serving**: 
  - Capability to serve static files efficiently.

- **Cookies and Headers Management**: 
  - Advanced handling of cookies and HTTP headers.

- **JSON Responses**: 
  - Rapidly generate JSON responses for API endpoints.

- **Request Body Parsing**: 
  - Read and parse JSON data from request bodies.

# Setup

Create a go.mod file in your project:

    go mod init myapp
Include the routego library in your project:

    go get -u github.com/imzeyn/routego

# Configuration and Usage
Below is an example of creating a simple HTTP server using `routego.Server`:

## Quick Start
```go
package main

import "github.com/imzeyn/routego"

func main() {
    urls := []routego.Path{
        {
			URL: "/",
			Handler: func(request *routego.Request, response *routego.Response) {
				template := `
				<html>
					<body style="background-color:black; color:white;">
						<h1>Hello!</h1>
					</body>
				</html>
				`
				response.Write(template) // This is how we return the HTTP response. Must be a string.
                // You can take action even after the response is returned
			},
            AllowedMethods: []string{"GET","POST"}, // You can specify the allowed methods here, if left blank the allowed method is only GET.
		},

        {
			URL: "/json-example",
			Handler: func(request *routego.Request, response *routego.Response) {
				json := map[string] string{
                    "msg": "Hello !"
                }
				response.JSON(json) // This way we can give JSON responses, it can be of any type.
                // You can take action even after the response is returned
			},
		},
    }
    s := routego.Server{
        Addr: ":80", // or provide "127.0.0.1:8080" etc..
        UrlPatterns: urls
    }

    s.Listen()
}
```

## URL Structure

In `routego`, we define URLs using the `routego.Path` method. This method allows you to set up nested routes and specify how each route should be handled. Below is an example demonstrating how to create a URL structure:

```go
url := routego.Path{
	ID: "mypathid", //Optional
	URL: "/user",
	Handler: func(request *routego.Request, response *routego.Response) {
		response.Write("User page")
	}, 
	
    Include: []routego.Path{
		{
			ID: "profile",
			URL: "/profile",
			Handler: func(request *routego.Request, response *routego.Response) {
				response.Write("User profile page")
			},
			Include: []routego.Path{...},
		},
	}, //Optional

	AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"}, //Optional
	Additional: "Something else...", //Optional
}
```

### URL Hierarchy

In the example above, we define a base URL `/user` and include a nested route `/profile`. The `Include` field allows for additional sub-routes to be specified. The hierarchical structure results in the following URL patterns:

- `/user`
- `/user/profile`

### URL Matching

The defined URL structure will match the following URLs:

- `/user` or `/user/` 
  - Matches the base route and will trigger the handler associated with `/user`.

- `/user/profile` or `/user/profile/`
  - Matches the nested route and will trigger the handler associated with `/user/profile`.

If you want to give a mandatory parameter, you should write it as `{parameter}`. If it is not mandatory and optional, it should be as `{{parameter}}`.
Example;
```go
routego.Path{
	URL: "/user/{name}",
	Handler: func(request *routego.Request, response *routego.Response) {
		response.JSON(request.Params)
	},
}

routego.Path{
	URL: "/user/{{name}}",
	Handler: func(request *routego.Request, response *routego.Response) {
		response.JSON(request.Params)
	},
}
```
### Allowed Methods
The `AllowedMethods` parameter is optional and if not entered, it can only be accessed via the `GET` method.

### Include
It should be an array consisting of `routego.Path`. If this part is added it will be concatenated with the parent `URL`.
With this structure, you can pull your URL list from different files and specify it here. 

### Summary

- **Base URL**: `/user`
- **Nested URL**: `/user/profile`
- **Matching URLs**: Both base and nested URLs can be accessed with or without a trailing slash.

This structure provides a clear and organized way to manage routes and their handlers, ensuring that nested paths are properly matched and handled.

# Handling requests


### `Request` Struct

The `Request` struct in the `routego` package provides methods to handle various aspects of HTTP requests, including cookies, headers, IP addresses, device information, JSON data, and file uploads.

### Methods

#### `Cookie() map[string]http.Cookie`

Returns a map of cookies present in the HTTP request. The map's keys are cookie names, and the values are `http.Cookie` objects.

**Example Usage:**
```go
cookies := request.Cookie()
```

#### `Header() map[string][]string`

Returns a map of HTTP headers from the request. The map's keys are header names, and the values are slices of header values.

**Example Usage:**
```go
headers := request.Header()
```

#### `IP() string`

Determines the IP address of the client making the request. It first checks the `X-Forwarded-For` and `X-Real-IP` headers. If these headers are not present, it falls back to the IP address from `RemoteAddr`.

**Example Usage:**
```go
clientIP := request.IP()
```

#### `DeviceInfo() (string, string)`

Extracts and returns the browser and operating system information from the `User-Agent` header. It returns two strings: the browser name and the operating system.

**Example Usage:**
```go
browser, os := request.DeviceInfo()
```

#### `JSON() (*interface{}, error)`

Parses the JSON body of the request and returns it as an `interface{}`. It returns a pointer to the parsed data and an error if any occurs during parsing.

**Example Usage:**
```go
data, err := request.JSON()
if err != nil {
    // Handle error
}
```

#### `Upload(to, name string) (bool, error)`

Handles file uploads from a multipart request. It saves the file with the given name to the specified destination path. Returns a boolean indicating success and an error if any occurs.

**Example Usage:**
```go
success, err := request.Upload("/path/to/save/file", "fileFieldName")
if err != nil {
    // Handle error
}
```


### `UploadIfValid(name, to string, signature *MimeSignatureList) (bool, string, error)`

This method handles file uploads from a multipart HTTP request and performs validation checks on the uploaded file. It ensures that the file's MIME type and extension match the provided signature list.

#### Parameters

- **`name`**: The name of the form field in the multipart request that contains the file.
- **`to`**: The destination path where the file should be saved if it is valid.
- **`signature`**: A pointer to a `MimeSignatureList` that contains MIME type signatures and associated file extensions for validation.

#### Returns

- **`bool`**: Returns `true` if the file is successfully validated and saved; otherwise, returns `false`.
- **`string`**: The file extension of the uploaded file if valid.
- **`error`**: An error if any occurs during processing, such as file read errors, validation failures, or file system issues.

#### Functionality

1. **Create a Multipart Reader**: 
   - Uses `r.HTTP.MultipartReader()` to create a multipart reader from the request's body.

2. **Iterate Through Parts**:
   - Iterates through each part of the multipart request using `mr.NextPart()`.

3. **Check Form Name**:
   - For each part, checks if the `FormName()` matches the specified `name` parameter. If it does, it processes the part as a file.

4. **File Validation**:
   - Extracts the file's MIME type and extension.
   - Reads the first 512 bytes of the file to obtain its signature.
   - Compares the file's signature and extension against the `MimeSignatureList` provided.

5. **Save Valid File**:
   - If the file's MIME signature and extension match one of the entries in the `MimeSignatureList`, it saves the file to the specified destination (`to`) with the appropriate extension.

6. **Error Handling**:
   - Returns `false` and an appropriate error if the file's signature or extension does not match, if the file cannot be read, or if saving the file fails.

#### Example Usage

```go
signatureList := &MimeSignatureList{
    {"image/jpeg", []byte{0xFF, 0xD8, 0xFF}, "image", []string{"jpg", "jpeg"}},
    {"image/png", []byte{0x89, 0x50, 0x4E, 0x47}, "image", []string{"png"}},
    // Add more MIME signatures as needed
}

success, ext, err := request.UploadIfValid("fileFieldName", "/path/to/save/file", signatureList)
if err != nil {
    // Handle error
}
if success {
    // File is valid and saved with extension ext
}
```

### Key Components

- **`MimeSignature`**: Represents a MIME type signature, including the MIME type, the signature bytes used for validation, the category of the file, and supported file extensions.
- **`MimeSignatureList`**: A list of `MimeSignature` entries used to validate uploaded files.
- **`MimeCategory`**: A category for grouping MIME types (e.g., `image`, `video`, `audio`, `document`).

This method ensures that uploaded files adhere to the expected MIME type and file extension criteria before saving, enhancing security and file integrity.


---

### `Response` Methods

The `Response` struct provides methods to manage HTTP responses, including setting cookies, headers, status codes, and writing response data in various formats.

#### `Cookie(cookie CookieObject)`

Sets a cookie in the response.

**Parameters:**
- **`cookie`**: A `CookieObject` that defines the properties of the cookie.

**Functionality:**
1. **Initialize Cookies Slice**: Checks if the `response.cookie` slice is `nil` and initializes it if necessary.
2. **Build Cookie String**: Constructs a cookie string from the `CookieObject` properties, including name, value, max age, path, domain, and flags (HttpOnly and Secure).
3. **Add Cookie to Response**: Appends the constructed cookie string to the `response.cookie` slice.

**Example Usage:**
```go
cookie := CookieObject{
    Name:     "session_id",
    Value:    "abc123",
    Path:     "/",
    Domain:   "example.com",
    MaxAge:   3600,
    SameSite: "Lax",
    HttpOnly: true,
    Secure:   true,
}

response.Cookie(cookie)
```

#### `Header(key, value string)`

Sets an HTTP header in the response.

**Parameters:**
- **`key`**: The name of the HTTP header.
- **`value`**: The value of the HTTP header.

**Functionality:**
1. **Initialize Headers Map**: Checks if the `response.headers` map is `nil` and initializes it if necessary.
2. **Set Header**: Adds or updates the header in the `response.headers` map with the provided key and value.

**Example Usage:**
```go
response.Header("Content-Type", "text/html; charset=utf-8;")
```

#### `Status(code int)`

Sets the HTTP status code for the response.

**Parameters:**
- **`code`**: The HTTP status code to set.

**Functionality:**
1. **Set Status Code**: Assigns the provided status code to the `response.status` field.

**Example Usage:**
```go
response.Status(200) // Sets the status code to 200 OK
```

#### `Write(w string)`

Writes a plain text response.

**Parameters:**
- **`w`**: The response body as a string.

**Functionality:**
1. **Set Default Content-Type**: If no `Content-Type` header is set, defaults to `"text/html; charset=utf-8;"`.
2. **Assign Response Data**: Converts the response body to a byte slice and assigns it to `response.data`.
3. **Signal Response Completion**: Sends a signal through the `done` channel indicating that the response is ready to be written.
4. **Wait for Writing**: Waits until the response has been fully written (synchronized using the `writed` channel).

**Example Usage:**
```go
response.Write("<html><body>Hello, world!</body></html>")
```

#### `JSON(w interface{}) error`

Writes a JSON response.

**Parameters:**
- **`w`**: The data to be serialized to JSON.

**Returns:**
- **`error`**: Returns an error if JSON marshaling fails.

**Functionality:**
1. **Marshal to JSON**: Converts the provided data to a JSON byte slice using `json.Marshal()`.
2. **Set Content-Type Header**: Sets the `Content-Type` header to `"application/json; charset=utf-8;"`.
3. **Assign Response Data**: Assigns the JSON data to `response.data`.
4. **Signal Response Completion**: Sends a signal through the `done` channel indicating that the response is ready to be written.
5. **Wait for Writing**: Waits until the response has been fully written (synchronized using the `writed` channel).

**Example Usage:**
```go
data := map[string]string{"message": "Hello, world!"}
err := response.JSON(data)
if err != nil {
    // Handle error
}
```

### Summary

The `Response` methods are designed to manage various aspects of an HTTP response in a structured manner. They allow setting cookies, headers, status codes, and writing data in both plain text and JSON formats. Each method ensures that the response is appropriately configured and sent back to the client.

Feel free to ask if you need further details or explanations!