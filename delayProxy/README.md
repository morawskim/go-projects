# DelayProxy

DelayProxy is a simple HTTP reverse proxy written in Go, designed to introduce a programmable delay before forwarding requests to a target destination. It is particularly useful for testing how applications handle slow network responses or timeouts.

## How it works

The proxy listens for incoming requests and expects a specific query parameter, `__targetUrl`, which defines where the request should eventually be sent. Before forwarding, the proxy pauses for a configured amount of time.

## Configuration Options

The application can be configured using the following command-line flags:

*   `-port`: The port on which the proxy server will listen.
    *   **Default**: `8080`
*   `-delay`: The duration to wait before forwarding the request. It accepts values in a format compatible with Go's `time.ParseDuration` (e.g., `5s`, `100ms`, `1m`).
    *   **Default**: `10s`

## Usage Requirements

Every request sent to the proxy must include the `__targetUrl` query parameter containing a valid URL. For example:
`http://localhost:8080/?__targetUrl=https://cdn.example.com/assets/styles.css`
