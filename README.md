[![Go Report Card](https://goreportcard.com/badge/github.com/tlkamp/log-reloading)](https://goreportcard.com/report/github.com/tlkamp/log-reloading)

# Hot Reloading for Logs
An example server that can live-reload logging configuration. The server itself just logs requests and manages its own config.

The logging configuration can be updated and reloaded at will. If the `/-/reload` endpoint detects that the configuration has not changed, no action will be taken (indicated in the logs).

## Demo
![log-reload-fps](https://user-images.githubusercontent.com/18516698/72223338-41b1d880-3533-11ea-9358-97ee0597ba6d.gif)

## Configuration
The server can be configured via flags on the command-line.


### Server

#### Available CLI Flags
| **Flag**         | **Type**    | **Default** | **Description**                              |
|------------------|-------------|-------------|----------------------------------------------|
| `-bind-address`  | String      | localhost   | The address at which to listen for requests. |
| `-config-file`     | String      | config.yaml  | The path to the config file to use.            |
| `-port`          | Int         | 8080        | The port on which to listen for connections. |

**Example**

```console
# Start the server on port 9090 using a config file called 'example.yaml'
$ ./log-reloading -port 9090 -config-file=example.yaml
INFO[0000] Server is starting at localhost:9090
```

### Logging
Logging configuration is in `YAML`.

```yaml
logging:
  level: info  # Any of [debug, info, warn, error, fatal, panic]
  colors: true # true or false
  format: text # text or json
```

## Endpoints
| **Endpoints**  | **Method** | **Description**                     |
|----------------|------------|-------------------------------------|
| `/-/config`     | GET        | View the current configuration.      |
| `/-/reload`    | POST       | Reload the configuration from disk.  |

## Build & Execute
```console
$ make run
Removing log-reloading
INFO[0000] Server is starting at localhost:8080
INFO[0008] 127.0.0.1:53863 GET /-/config
INFO[0013] 127.0.0.1:53865 POST /-/reload
INFO[0013] Hashes match. Skipping reload.
```
