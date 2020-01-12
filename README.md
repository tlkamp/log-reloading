# Hot Reloading for Logs
An example server that automatically reloads its logging configuration.

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

## Build & Execute
```
$ go build
$ ./log-reloading
INFO[0000] Server is starting at localhost:8080
INFO[0008] 127.0.0.1:53863 GET /-/config
INFO[0013] 127.0.0.1:53865 POST /-/reload
INFO[0013] Hashes match. Skipping reload.
```
