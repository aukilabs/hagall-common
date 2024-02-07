# Hagall Common Libraries

This repository contains packages that are commonly used by Hagall & other backend services in Aukilabs.

## Packages

| Package                | Description                                                      |
| ---------------------- | ---------------------------------------------------------------- |
| [crypt](crypt)         | Package that provides cryptography related functionalities.      |
| [hdsclient](hdsclient) | Package with a client interface to the Hagall Discovery Service. |
| [http](http)           | Package with common HTTP functionalities.                        |
| [messages](messages)   | Package with the definition of Hagall modules protobuf messages  |
| [ncsclient](ncsclient) | Package with a client interface to the Network Credit Service.   |
| [scenario](scenario)   | Package to support Hagall protocol simulation using websocket.   |
| [smoketest](smoketest) | Package that provides smoketest functionality.                   |
| [testing](testing)     | Package contains functions to support Hagall testing.            |
| [websocket](websocket) | Package with functions to manage websocket communications.       |

## Generating Protobuf

To regenerate protobuf messages after updated .proto files.

```shell
make proto
```
