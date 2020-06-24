# Fupisha Project Layout

Fupisha uses the `/internal` layout pattern.You put all of the private code in the `/internal` directory the application code in the `cmd` directory will be limited to small files that define the `main` function for the corresponding application binaries. Everything else will be imported from the internal or pkg subdirectory.

We expose restful application interface code outside internal directory in the `/api` directory to limit the size of internal packages.

Below is the current project layout in use with `Fupisha`.

> This is not a permanent layout. It is subject to change based on community view. However the current project layout puts into consideration future growth of the project.

```
fupisha
│   README.md
|   LICENSE
|   docker-compose.yml
|   go.mod
|   go.sum
│   init-mongo.sh
|   .gitignore
|
└───templates
|   |   verify.html
|   |   reset.html
|   |   ...
│
└───internal
│   └───pkg
│   |   └───v1
|   |       └───auth
|   |       |   | key.go
|   |       |   | token.go
|   |       |   | token_test.go
|   |       |   | ...
|   |       |
|   |       └───config
|   |       |   | config.go
|   |       |   | ...
|   |       |
|   |       └───email
|   |       |   | content.go
|   |       |   | email.go
|   |       |   | email_test.go
|   |       |   | ...
|   |       |
|   |       └───encoding
|   |       |   | uuid_v4.go
|   |       |   | ...
|   |       |
|   |       └───logging
|   |           | logger.go
|   |           | ...
|   |
|   └───store
|       | store.go
|       |
|       └───mongodb
|       |   | store.go
|       |   | store_test.go
|       |   | url.go
|       |   | url_test.go
|       |   | user.go
|       |   | user_test.go
|       |   | ...
|       |
|       └───model
|           | url.go
|           | user.go
|           | ...
|
└───api
|   │ api.go
|   │ server.go
|   | ...
|   |
|   └───v1
|       └───auth
|       |   | auth.go
|       |   | errors.go
|       |   | handlers.go
|       |   | handlers_test.go
|       |   | middlewares.go
|       |   | routes.go
|       |   | ...
|       |
|       └───links
|       |   | links.go
|       |   | errors.go
|       |   | handlers.go
|       |   | handlers_test.go
|       |   | routes.go
|       |   | ...
|       |
|       └───users
|       |   | users.go
|       |   | errors.go
|       |   | handlers.go
|       |   | handlers_test.go
|       |   | routes.go
|       |   | ...
|       | ...
|
└───docs
|   | api.md
|   | project_layout.md
|   | ...
|
└───cmd
    | main.go

```

# Directory Structure

Below is a detailed explantion of each directory and its purpose.

## /cmd

This directory houses our `main` binary. **What does it mean?** It means you can use the go get command to fetch (and install) fupisha project, its applications and its libraries (e.g., go get github.com/nairobi-gophers/fupisha/cmd/main). You don’t have to separate the application files.

## /docs

This directory houses the project's documentation files.

## /templates

This directory houses html templates to be used by the internal email package. For example when a user signsup Fupisha needs to send that user a verification email. This does not restrict files housed int here to the email package alone. Any other packages that might need to work with html files can also have their html files placed here probably in a subdirectory to explicitly define their use.

## /internal

This directory houses private code, that which Fupisha does not wish to be imported or used in third party libraries or people's projects. Most of the code in here describes the core of Fupisha project.

### /internal/pkg/v1

This directory `pkg/v1` housed within internal provides physical versioning boundaries. It enforces code to conform to the described version e.g. v1 .This enables to support older versions of the API without much effort as otherwise needed.

### /internal/store

This directory is houses the database logic code. It is placed outside `/pkg/v1` because we don't need to explicitly version it. **Why?** Versioning in this layer can be done using migrations, which is the normal for majority of projects of this kind.

## /api

This houses the restful api logic. Most of the code here works with the HTTP/1 or HTTP/2 protocol. This layer defines the Fupisha restful API for use by third party applications and web clients. In turn this api imports internal code from the packages within `/internal` directory.

## /api/v1

This directory `v1` helps provide explicit versioning of code by enforcing code to conform to the described version. It helps prevent mixing of different versions of code since it will be a part of the import path.

## /api/v1/...

Directories inside of `v1` follow a domain based design, where each directory is named after the domain in which it is addressing. e.g. auth - addresses authentication domain. users - addresses users' domain ,e.t.c. You will notice that each domain-based directory has its own set of the following files:

- errors.go
- handlers.go
- {domain-name}.go
- routes.go

This makes it easier for when debugging / troubleshooting is called for. It becomes easier to pinpoint where the error is coming from or which package is responsible.

> You may also have note that `middlewares.go` is only in the auth package. This is because all other domain require to be authenticated before allowing access to their logic.Thus it makes sense for `users` domain to use the `auth` middleware to authorize requests before further processing.
