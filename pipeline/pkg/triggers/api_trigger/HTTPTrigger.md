# HTTP-trigger

This trigger component exposes an HTTP API. Requests can be sent in various
formats, and responses are given in JSON.

## Initialization configuration

This component may be initialized with a some options for authentication and
Cross Origin Resource Requests (CORS).

**Bearer authentication**

A list of api keys which, which may be used for bearer authentication. This
cannot be used at the same time as basic authentication.

```toml
apiKeys = ["foo", "bar"]
```

**Basic authentication**

A list of username / password combinations. This cannot be used at the same time
as bearer authentication.

```toml
[[basicAuth]]
username = "jim"
password = "secret"
```

**Allowed Origins (for CORS)**

A list of origin domains. Specifies which domains may connect to the API.

In order to interact with the HTTP API with a browser the domain name on which
the application is hosted needs to be included in this list. If the special
`"*"` value is present in the list, all origins will be allowed.

```toml
AllowedOrigins = ["app.example.com"]
```

## Usage

The API accepts an HTTP POST request with a payload on the root of the component
address. The component accepts data in any format, and will output it as string
under the `http-trigger` key in the invocation state object.

The raw request body will be JSON-stringified and added to the ‘body’ field.

A special case is when the header `Content-Type: application/json` is set. In
this case the component will attempt to parse the request body as JSON before
adding it to the ‘json’ field. If the header is set, but the body cannot be
parsed to JSON an error will be returned.

**Optional Headers:**

These headers will be accepted by the component.

| Header        | required | example                        |
| ------------- | -------- | ------------------------------ |
| Authorization | false    | Bearer MyApiAuthKey            |
|               |          | Basic dXNlcm5hbWU6cGFzc3dvcmQ= |
| Content-Type  | false    | application/json               |

When Bearer authentication is used you may use one of the specified apiKeys.
When Basic authentication is used the header should confirm to the Basic Auth
standard.

### Outputs

This component will wait until all actions in the pipeline have completed, and
will then return with a response code and a body.

If the pipeline completed successfully the response code will be 200. Using the
pipeline response configuration you may specify what this component returns.

When an error occurs in the pipeline, a 500 error status code is returned with
an error message in the body.

For example; to echo the body that was received you may specify the body field
from the http-trigger component.

```
body = $.http-trigger.body
```

- header  
  A map each value representing one of the received headers
- body  
  An object containing the string of the request body, or the parsed JSON, if
  content-type was set to JSON
- raw  
  Raw is the original http request.
