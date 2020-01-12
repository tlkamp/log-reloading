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
