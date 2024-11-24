## Overview

This Go application is a simple HTTP proxy server that forwards requests to an upstream server (only GET requests).
Proxy server sends periodic updates to keep the connection alive during long upstream requests. 

The server is designed for scenarios where you need to proxy client requests, because
client expect to receive response in 30 seconds, but your server require much more time to send response.

## API Usage

### Endpoint

`GET /?url=<upstream_url>`

#### Parameters

* `url`: URL to forward the request to (required).

#### Headers

The proxy includes headers from the client, excluding `Host`,
and adds the `X-Forwarded-For` header with the client's IP address.

### Response Structure

#### Successful Response

The server returns the upstream response wrapped in a JSON object:

```
{
    "statusCode": 200,
    "headers": {
        "Content-Type": ["application/json"]
    },
    "body": "Response body from the upstream service"
}
```
#### Error Response

In case of errors, the server responds with a JSON object:

```
{
    "error": "Error message"
}
```
