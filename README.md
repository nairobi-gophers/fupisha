# fupisha

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

Fupisha is a modern url shortening service.

:construction: Work In Progress! :construction: Contributions and bug reports are welcome. :tada:

# Features

- Free and open source.
- URL Shortening.
- Visitor Counting.
- Expirable Links.
- URL deletion.
- Restful API.

# Stack

- Go (language)
- go-chi/chi (http routing)
- mongo (database)
- redis (cache layer)
- vuejs (web UI library)
- vuex (state management)
- vuetify (vuejs material design framework)
- Circle CI (Continous integration)

# Setup

You need to have Golang installed at the very least `Go 1.11`. The easiest way to set up fupisha in your local environment is to install it using go get:

```Go
 go get github.com/nairobi-gophers/fupisha
```

Run the mongo container which the app will use. Later we will also dockerize the app.

```Go
 docker-compose up
```

# Why build this

It will involve the community and cool techniques like:

- Golang unit testing
- VueJS
- Makefiles
- Circle CI
- Document Databases
- In Memory caching
- Dockerfile, Docker Compose and Docker Image/Container Creation

## License

Copyright 2020 Nairobi Gophers

Licensed under [the MIT License](LICENSE)
