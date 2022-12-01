mongolia
========
Mongo object document mapper for golang.

![IBEX](docs/ibex.png)

This ODM does the following:
1) accept a jsonschema for each type of struct (i.e. model) that will be used by the client application
2) establish a one-to-one mapping from each application object to a corresponding document in the database
3) restrict access to the underlying document; it can only be updated through operations on the application object
4) deny any operation which would cause a document to violate the schema for its application object

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
