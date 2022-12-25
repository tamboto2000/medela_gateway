# Medela Gateway
API Gateway for PT. Medela Potentia interview task

## Installation
Clone this repository, enter to project root directory and run this command to build the project:
`make build`

To run the build, execute this command:
`make run`

This command will run the app and will also trying to find default config `config.json`.
You can provide the config file by creating one.

## Configuration file overview
A configuration file is needed to run the app so that it can function properly.
The configuration file consist of registered endpoints pointed to determined server or backend, example:
```json
{
    "endpoints": [
        {
            "endpoint": "/api/v1/user/:id",
            "method": "GET",
            "backend": {
                "host": "https://userdata.com",
                "url_pattern": "/user/:id",
                "method": "GET"
            }
        }
    ]
}
```

From above example, the config file is telling the gateway to add route to `GET https://userdata.com/user/:id` by endpoint `/api/v1/user/:id`.
The `:id` is path parameter, `endpoint` and `backend.url_pattern` must has the same path parameters (if any, this is optional) or the http requests will not propely work

*Currently the gateway can only support JSON response data because of the time constraints to complete this project*

## Add an endpoint
To add an endpoint, append a new object inside `endpoints` on the config file, example:
```json
{
    "endpoints": [
        {
            "endpoint": "/api/v1/user/:id",
            "method": "GET",
            "backend": {
                "host": "https://userdata.com",
                "url_pattern": "/user/:id",
                "method": "GET"
            }
        }
    ]
}
```
The configuration attributes of endppoints are:
| Attribute name | Mandatory? | Type     | Description                                                                                                                                                                     |
|----------------|------------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `endpoint`     | Y          | `string` | The URL you want to expose. You can add path parameter to it by using preffix `:`, example: `/product/:label`                                                                   |
| `method`       | Y          | `string` | The supported method to request this endpoint                                                                                                                                   |
| `backend`      | Y          | `object` | The server that will handle the request to this endpoint. [See the backend configuration](#backend-configuration)                                                               |
| `middlewares`  | N          | `array`  | Collection of middlewares. When present, every request to this endpoint will be passed through the middlewares first. See [middleware configuration](#middleware-configuration) |

## Backend configuration
A backend is the server that resposinble to handle a request to an endpoint. Currently every endpoint can only support exactly 1 backend server. Example:
```json
{
    "host": "https://userdata.com",
    "url_pattern": "/user/:id",
    "method": "GET"
}
```

Attributes of backend are:
| Attribute name | Mandatory? | Type     | Description                                                                                                                                               |
|----------------|------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
| `host`         | Y          | `string` | The server/backend host, can be a domain or IP. You need to specify the scheme as the preffix, example: `https://foo.com`, `foo.com` is not a valid value |
| `url_pattern`  | Y          | `string` | The URL that will handle the request to an endpoint                                                                                                       |
| `method`       | Y          | `string` | The supported method by this backend                                                                                                                      |

## Technical references
- Reverse proxy
https://itnext.io/why-should-you-write-your-own-api-gateway-from-scratch-378074bfc49e
- Response modifier
https://stackoverflow.com/questions/31535569/how-to-read-response-body-of-reverseproxy
https://forum.golangbridge.org/t/how-to-build-an-http-reverse-proxy-to-dynamically-change-the-content-of-the-upstream-response/1313