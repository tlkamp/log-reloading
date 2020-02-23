[![Go Report Card](https://goreportcard.com/badge/github.com/tlkamp/log-reloading)](https://goreportcard.com/report/github.com/tlkamp/log-reloading)

# Hot Reloading for Logs
An example server that can live-reload logging configuration.

The logging configuration can be updated and reloaded at will. If the `/-/reload` endpoint detects that the configuration has not changed, no action will be taken (indicated in the logs).

## Demo
![log-reload-fps](https://user-images.githubusercontent.com/18516698/72223338-41b1d880-3533-11ea-9358-97ee0597ba6d.gif)

## Configuration
The server can be configured via flags on the command-line:
```shell
$ ./log-reloading -h
Usage of ./log-reloading:
  -bind-address string
        the address to bind to. (default "localhost")
  -port int
        the port to listen on. (default 8080)

# Start the server on port 9090
$ ./log-reloading -port 9090
INFO[0000] Server is starting at localhost:9090
```

Logging configuration is done in `YAML`.

```yaml
logging:
  level: info
  colors: true
  format: text
```

## Endpoints
| **Endpoints**  | **Method** | **Description**                     |
|----------------|------------|-------------------------------------|
| `/-/config`    | GET        | View the current configuration.     |
| `/-/reload`    | POST       | Reload the configuration from disk. |

## Build & Execute
```
$ go build
$ ./log-reloading
INFO[0000] Server is starting at localhost:8080
INFO[0008] 127.0.0.1:53863 GET /-/config
INFO[0013] 127.0.0.1:53865 POST /-/reload
INFO[0013] Hashes match. Skipping reload.
```
