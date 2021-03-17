books-api
================
Digital Publishing Books API

### Getting started

#### Pre-requisites

Install and run a mongoDB

- Using homebrew:
    - `brew install mongodb`
    - If not automatically started, use `brew services start mongodb`
    - Stop mongoDB by running `brew services stop mongodb`

#### Other Dependencies

* No further dependencies other than those defined in `go.mod`

### Running the Application and unit tests

- Run application with `make debug`
- Run unit test with `make test`

### Configuration

| Environment variable         | Default         | Description
| ---------------------------- | --------------- | ------------------------------------------------------------------------------------------------------------------ |
| BIND_ADDR                    | :8080           | The host and port to bind to                                                                                       |
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s              | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  |
| HEALTHCHECK_INTERVAL         | 30s             | Time between self-healthchecks (`time.Duration` format)                                                            |
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s             | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| MONGODB_BIND_ADDR            | localhost:27017 | The MongoDB bind address                                                                                           |
| MONGODB_COLLECTION           | books           | The MongoDB images database                                                                                        |
| MONGODB_DATABASE             | bookStore       | MongoDB collection                                                                                                 |

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.