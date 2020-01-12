# Hot Reloading for Logs
An example server that automatically reloads its logging configuration.

The server configuration is loaded **once** on startup. The logging configuration can be updated and reloaded at will. If the `/-/reload` endpoint detects that the configuration has not changed, no action will be taken (indicated in the logs).

## Demo
// TODO: Add gif


## Configuration
Configuration is done in `YAML`. The `server` section is required.

```yaml
server:
  address: localhost
  port: "8080"
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
