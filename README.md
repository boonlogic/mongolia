mongolia
========
Mongo object document mapper for golang.

![IBEX](docs/ibex.png)


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

todo: update below as code evolves

Environment
-----------
* `AMBER_MONGO_URI`: MongoDB [connection URI](https://www.mongodb.com/docs/manual/reference/connection-string/) for the mongo instance to use.
* `AMBER_DB_NAME`: Name of mongo database to connect to.

Organization
------------
```
├── api
├── pkg                 <- Go package source
└── cmd                 <- runnable server implementing the OpenAPI 3 spec
```
