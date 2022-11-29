mongolia
========
Mongo object document mapper for golang.

![IBEX](docs/ibex.png)

Setup
-----
From top level of `builder`:
1. run `make mongolia-env` to install environment dependencies
2. run `make mongolia-pull` to pull the latest source

Usage
-----
See `test/smoke_test.go` for a minimal but complete usage example.

Environment
-----------
* `MONGOLIA_URI` connection URI for the underlying mongo instance.
* `MONGOLIA_DBNAME`: name of mongo database to use.

Organization
------------
```
├── cmd        <- server executable
├── docs       <- documentation
├── mongolia   <- package `mongolia` (core package)
├── restapi    <- package `restapi` (runs mongolia as a web service)
└── test       <- code for testing
```
