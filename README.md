# hawkeye_management_svc
Hawkeye Management Service

# Hawkeye Management Service (Utility Service)

[![GoDoc](https://godoc.org/github.com/ramesharun/hawkeye_management_svc?status.png)](http://godoc.org/github.com/ramesharun/hawkeye_management_svc)
[![Build Status](https://github.com/ramesharun/hawkeye_management_svc/workflows/build/badge.svg)](https://github.com/ramesharun/hawkeye_management_svc/actions?query=workflow%3Abuild)
[![Code Coverage](https://codecov.io/gh/ramesharun/hawkeye_management_svc/branch/master/graph/badge.svg)](https://codecov.io/gh/ramesharun/hawkeye_management_svc)
[![Go Report](https://goreportcard.com/badge/github.com/ramesharun/hawkeye_management_svc)](https://goreportcard.com/report/github.com/ramesharun/hawkeye_management_svc)

This starter kit is designed to get you up and running with a project structure optimized for developing
RESTful API services in Go. It promotes the best practices that follow the [SOLID principles](https://en.wikipedia.org/wiki/SOLID)
and [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). 
It encourages writing clean and idiomatic Go code. 

The kit provides the following features right out of the box:

* RESTful endpoints in the widely accepted format
* Standard CRUD operations of a database table
* JWT-based authentication
* Environment dependent application configuration management
* Structured logging with contextual information
* Error handling with proper error response generation
* Database migration
* Data validation
* Full test coverage
* Live reloading during development
 
The kit uses the following Go packages which can be easily replaced with your own favorite ones
since their usages are mostly localized and abstracted. 

* Routing: [ozzo-routing](https://github.com/go-ozzo/ozzo-routing)
* Database access: [ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
* Database migration: [golang-migrate](https://github.com/golang-migrate/migrate)
* Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
* Logging: [zap](https://github.com/uber-go/zap)
* JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://golang.org/doc/install) to
install Go on your computer. The kit requires **Go 1.17 or above**.

[Docker](https://www.docker.com/get-started) is also needed if you want to try the kit without setting up your
own database server. The kit requires **Docker 17.05 or higher** for the multi-stage build support.

After installing Go and Docker, run the following commands to start experiencing this starter kit:

```shell
# download the starter kit
git clone https://github.com/ramesharun/hawkeye-management-svc.git

cd hawkeye-management-svc

# start a PostgreSQL database server in a Docker container
#make db-start

# seed the database with some test data
#make testdata

# run the RESTful API server
make run

# or run the API server with live reloading, which is useful during development
# requires fswatch (https://github.com/emcrisostomo/fswatch)
make run-live
```

At this time, you have a RESTful API server running at `http://127.0.0.1:8080`. It provides the following endpoints:

* `GET /healthcheck`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
* `POST /v1/login`: authenticates a user and generates a JWT
* `GET /v1/monitors`: returns a paginated list of the monitors
* `GET /v1/monitors/:id`: returns the detailed information of an monitors
* `POST /v1/monitors`: creates a new monitor
* `PUT /v1/monitors/:id`: updates an existing monitor
* `DELETE /v1/monitors/:id`: deletes an monitor

Try the URL `http://localhost:8080/healthcheck` in a browser, and you should see something like `"OK v1.0.0"` displayed.

If you have `cURL` or some API client tools (e.g. [Postman](https://www.getpostman.com/)), you may try the following 
more complex scenarios:

```shell
# authenticate the user via: POST /v1/login
curl -X POST -H "Content-Type: application/json" -d '{"username": "demo", "password": "pass"}' http://localhost:8080/v1/login
# should return a JWT token like: {"token":"...JWT token here..."}

# with the above JWT token, access the album resources, such as: GET /v1/albums
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/v1/monitors
# should return a list of monitor records in the JSON format
```

To use the starter kit as a starting point of a real project whose package name is `github.com/abc/xyz`, do a global 
replacement of the string `github.com/ramesharun/hawkeye-management-svc` in all of project files with the string `github.com/abc/xyz`.


## Project Layout

The service uses the following project layout:
 
```
.
├── cmd                  main applications of the project
│   └── server           the API server application
├── config               configuration files for different environments
├── internal             private application and library code
│   ├── monitor          monitor-related features
│   ├── auth             authentication feature
│   ├── config           configuration library
│   ├── entity           entity definitions and domain logic
│   ├── errors           error types and handling
│   ├── healthcheck      healthcheck feature
│   └── test             helpers for testing purpose
├── migrations           database migrations
├── pkg                  public library code
│   ├── accesslog        access log middleware
│   ├── graceful         graceful shutdown of HTTP server
│   ├── log              structured and context-aware logger
│   └── pagination       paginated list
└── testdata             test data scripts
```

The top level directories `cmd`, `internal`, `pkg` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Within `internal` and `pkg`, packages are structured by features in order to achieve the so-called
[screaming architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html). For example, 
the `monitor` directory contains the application logic related with the album feature. 

Within each feature package, code are organized in layers (API, service, repository), following the dependency guidelines
as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).


## Common Development Tasks

This section describes some common development tasks using this starter kit.

### Implementing a New Feature

Implementing a new feature typically involves the following steps:

1. Develop the service that implements the business logic supporting the feature. Please refer to `internal/monitor/service.go` as an example.
2. Develop the RESTful API exposing the service about the feature. Please refer to `internal/monitor/api.go` as an example.
3. Develop the repository that persists the data entities needed by the service. Please refer to `internal/monitor/repository.go` as an example.
4. Wire up the above components together by injecting their dependencies in the main function. Please refer to 
   the `monitor.RegisterHandlers()` call in `cmd/server/main.go`.

### Working with DB Transactions

It is the responsibility of the service layer to determine whether DB operations should be enclosed in a transaction.
The DB operations implemented by the repository layer should work both with and without a transaction.

You can use `dbcontext.DB.Transactional()` in a service method to enclose multiple repository method calls in
a transaction. For example,

```go
func serviceMethod(ctx context.Context, repo Repository, transactional dbcontext.TransactionFunc) error {
    return transactional(ctx, func(ctx context.Context) error {
        repo.method1(...)
        repo.method2(...)
        return nil
    })
}
```

If needed, you can also enclose method calls of different repositories in a single transaction. The return value
of the function in `transactional` above determines if the transaction should be committed or rolled back.

You can also use `dbcontext.DB.TransactionHandler()` as a middleware to enclose a whole API handler in a transaction.
This is especially useful if an API handler needs to put method calls of multiple services in a transaction.


### Managing Configurations

The application configuration is represented in `internal/config/config.go`. When the application starts,
it loads the configuration from a configuration file as well as environment variables. The path to the configuration 
file is specified via the `-config` command line argument which defaults to `./config/local.yml`. Configurations
specified in environment variables should be named with the `APP_` prefix and in upper case. When a configuration
is specified in both a configuration file and an environment variable, the latter takes precedence. 

The `config` directory contains the configuration files named after different environments. For example,
`config/local.yml` corresponds to the local development environment and is used when running the application 
via `make run`.

Do not keep secrets in the configuration files. Provide them via environment variables instead. For example,
you should provide `Config.DSN` using the `APP_DSN` environment variable. Secrets can be populated from a secret
storage (e.g. HashiCorp Vault) into environment variables in a bootstrap script (e.g. `cmd/server/entryscript.sh`). 

## Deployment

The application can be run as a docker container. You can use `make build-docker` to build the application 
into a docker image. The docker container starts with the `cmd/server/entryscript.sh` script which reads 
the `APP_ENV` environment variable to determine which configuration file to use. For example,
if `APP_ENV` is `qa`, the application will be started with the `config/qa.yml` configuration file.

You can also run `make build` to build an executable binary named `server`. Then start the API server using the following
command,

```shell
./server -config=./config/prod.yml
```

```