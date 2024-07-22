# HTTP Action

Call an HTTP API or webservice.

Limitation: *Only method HTTP POST available*

## Input fields

* Url
`string` - the url to make a call to; requires to be prepended by http:// or https://.
* Method
`string` - HTTP method, only "POST" available
* RequestBody
`[]byte` - Request body. 
* ContentType
`string` - the content type, e.g. "application/json"

## Response fields

* ResponseBody
`[]byte` - Response body that is returned by the webservice / API.

* ResponseHeaders 
`map[string][]string` - Response headers of the HTTP request.

* ResponseStatus
`int` - Response code of the request, e.g. 200, 500 