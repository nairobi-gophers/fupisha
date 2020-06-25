# Fupisha Project Layout Ideas

## **Layout I**

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
└───cmd
|   | main.go
|   | ...
|
└───docs
|   | api.md
|   | project_layout.md
|   | ...
|
└───store
|       | store.go
|       |
|       └───mongodb
|       |   └───model
|       |   |   | url.go
|       |   |   | user.go
|       |   |   |...
|       |   |
|       |   |   store.go
|       |   |   store_test.go
|       |   |   url.go
|       |   |   url_test.go
|       |   |   user.go
|       |   |   user_test.go
|       |   |   ...
|       |
|       | ...
|
|   config.go
|   apikey.go
|   apitoken.go
|   apitoken_test.go
|   emailcontent.go
|   email.go
|   logger.go
|   uuidv4.go
|   middlewares.go
|   auth.go
|   auth_handlers.go
|   auth_errors.go
|   auth_routes.go
|   link.go
|   link_handlers.go
|   link_errors.go
|   link_routes.go
|   user.go
|   user_handlers.go
|   user_errors.go
|   user_routes.go
|   visit.go
|   visit_errors.go
|   visit_handlers.go
|   visit_routes.go
|   api.go
|   server.go

```

## **Layout II**

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
└───cmd
|   | main.go
|   | ...
|
└───docs
|   | api.md
|   | project_layout.md
|   | ...
|
└───store
|       | store.go
|       |
|       └───mongodb
|       |   └───schema
|       |   |   | url.go
|       |   |   | user.go
|       |   |   |...
|       |   |
|       |   |   store.go
|       |   |   store_test.go
|       |   |   url.go
|       |   |   url_test.go
|       |   |   user.go
|       |   |   user_test.go
|       |   |   ...
|       |
|       └───postgres
|           └───schema
|           |   | schema.sql
|           |   |...
|           |
|           |   store.go
|           |   store_test.go
|           |   url.go
|           |   url_test.go
|           |   user.go
|           |   user_test.go
|           |   ...
|
|   config.go
|   apikey.go
|   apitoken.go
|   apitoken_test.go
|   emailcontent.go
|   email.go
|   logger.go
|   uuidv4.go
|   middlewares.go
|   auth.go
|   errors.go
|   routes.go
|   link.go
|   handlers.go
|   user.go
|   visit.go
|   api.go
|   server.go
```
