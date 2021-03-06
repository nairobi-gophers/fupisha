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
- postgresql (database)
- redis (cache layer)
- vuejs (web UI library)
- vuex (state management)
- vuetify (vuejs material design framework)
- Circle CI (Continous integration)

# Setup

You need to have Golang installed at the very least `Go 1.11`. The easiest way to get started with fupisha in your local environment is to clone it using git:

```
 git clone https://github.com/nairobi-gophers/fupisha.git
```  

Rename the `example.env` file to `.env` file and fill in all the config sections, these will be used by the server container to set up the necessary resources.For the smtp section the values must be valid and existing smtp account credentials.   

# Run
To run the application, you will need to ensure that you have the `make` utility installed and running in your local computer.If you have the `make` utility, 
then you can follow along with the below instructions.  

- Run tests in the container  

    `make tests`

- Start the api server container  

    `make up`

- Stop the api server container 

    `make down`

- Watch the api server container logs
        
    `make logs`  
          
## Implemented Features
- [x] Signup 
- [x] Login 
- [x] Shorten URL
- [x] URL Redirection

## Sample HTTP Requests
- Signup 

```
curl -X POST -H "Api:v1" -d '{"email":"admin@fupisha.io","password":"iamjustauser"}' http://localhost:8888/auth/signup
```

- Login 
```
curl -X POST -H "Api:v1" -d '{"email":"admin@fupisha.io","password":"iamjustauser"}' http://localhost:8888/auth/login
```

- Shorten URL
```
curl -X POST -H "Api:v1"  -d '{"url":"https://pkg.go.dev/github.com/mailgun/groupcache"}' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjIzOTkzMjYsImlhdCI6MTYyMjM5ODk2NiwiaXNzIjoiZnVwaXNoYSIsIl91aWQiOiJkNzkwOTVmYy0zYjFlLTQ4MTUtOTNmZS1lMTk4MTI0ZTcxZDMifQ.EGD2up_7iJP5FS_OhlR6UYvEqu9orRgU1iR65u-3Hrg" http://localhost:8888/url/shorten
```

- URL Redirection
```
curl -X GET http://localhost:8888/a3UdbL
```
# Why build this

It will involve the community and awesome technologies like:

- Golang unit testing
- VueJS
- Makefiles
- Circle CI
- Relational Database
- In Memory caching
- Docker Compose and Docker

Let's learn together.

## License

Copyright 2020 Nairobi Gophers

Licensed under [the MIT License](LICENSE)
