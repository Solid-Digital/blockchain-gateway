# Readme

http-api implements the pipeline interfaces.

## Supported features

* consume json request
* consume other request
* token authentication

## Testing

Run `make test` from the package root to run the tests.

## Configuration 

| item         | description           |
|--------------|-----------------------|
| apiKeys      | list of api keys      |
| CORS		   | Allow CORS (bool)			   |

Configuration of api keys is optional. If you don't configure any, the api is public.

### Example configuration

````toml
apiKeys = ["foo", "bar"]
CORS = false
````

### Outputs

| key    | description                      |
|--------|----------------------------------|
| header | map of string to list of strings |
| body   | string                           |