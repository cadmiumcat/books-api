books-api
================
Digital Publishing Books API


### Getting started
[Swagger specification](https://cadmiumcat.github.io/books-api/)

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
| MONGODB_BOOKS_COLLECTION     | books           | The MongoDB books collection                                                                                       |
| MONGODB_REVIEWS_COLLECTION   | reviews         | The MongoDB reviews collection                                                                                     |
| MONGODB_DATABASE             | bookStore       | MongoDB database                                                                                                   |
| DEFAULT_MAXIMUM_LIMIT        | 1000            | Pagination: maximum number of items returned                                                                       |
| DEFAULT_LIMIT                | 20              | Pagination: default number of items returned                                                                       |
| DEFAULT_OFFSET               | 0               | Pagination: default number of documents into the full list that a response starts at                               |

### Electronic Library Design

See [ARCHITECTURE](architecture/README.md) Source of truth for processing of book data for search.

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020 - 2021 Catalina Garcia

Released under MIT license, see [LICENSE](LICENSE.md) for details.