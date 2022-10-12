mongolia
========

Golang object document mapper for mongo.

Setup
-----
From top level of `builder`:
1. run `make mongolia-env` to installs environment dependencies
2. run `make mongolia-pull` to pull the latest source

Usage
-----
From top level of repository:
* run `make test` to run all tests
* run `make run` to run server
* run `make generate` to generate all generated files
* run `make tidy` to tidy go modules

Environment
-----------
* `AMBER_MONGO_URI`: MongoDB [connection URI](https://www.mongodb.com/docs/manual/reference/connection-string/) for the database instance to use

Organization
------------
```
├── api
├── pkg                 <- Go package source
└── cmd                 <- runnable server implementing the OpenAPI 3 spec
```

Roadmap
-------
