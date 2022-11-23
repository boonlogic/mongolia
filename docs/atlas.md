Here is the state of the union of Boon Atlas. I have some recommendations on conventions to use going forward.

Tenants
-------

`boonlogic` is the name of our tenancy within Atlas. Our tenancy has a set of Project Spaces. boonlogic has one main Project Space, `amber-aws`.*

```
boonlogic          <- tenancy
   |--- amber-aws  <- project space
```

If we use Atlas for a non-Amber product in the future, that would go under a new Project Space. Different versions of the same product belong in the same Project Space.

Clusters
--------

A Project Space contains one or more clusters. `amber-aws` has one cluster `Amber0`.*

```
amber-aws
   |--- Amber0     <- cluster
```

A cluster is a set of three database replicates running on physically co-located instances. It is associated with a cloud provider and one of that provider's resource regions. `Amber0` is located in AWS region us-east-1. If we deployed Amber to us-west-1, it would be under a new cluster.

The naming convention for clusters is [Product][n] where n just increments each time we add another cluster. 

```
Amber0 <- Amber runs on this cluster right now.
Amber1 <- Amber v2 will run on this cluster.
```

Let's say a future version of AVIS uses Atlas...

```
AVIS2  <- AVIS v3
Buc3   <- Buc cloud DB or something
etc.
```

This allows our esteemed database janitor to see at a glance which project the cluster is associated with, but also an identifier within the growing list of Atlas cluster deployments.

I enabled Termination Protection which **prevents any user from accidentally deleting a cluster**. To delete a cluster, Termination Protection will have to be disabled by a superuser first.

Databases
---------

Each cluster has some number of databases. A database is the top-level container of all application data for one deployment of an application. `Amber0` has three databases.
```
Amber0
   |--- amber-prod   <- persists all Amber production data
   |--- amber-dev    <- sandbox for Amber development
   |--- aop-license  <- persists license information for amber-on-prem
```

The DB janitor is really just needs to keep `amber-dev` usable enough for development, and guard `amber-prod` with their life.

Jim created cluster `Amber1` today and the databases for the Amber v2 prod/dev deployments. Now we have:

```
Amber0
   |--- amber-prod
   |--- amber-dev
   |--- aop-license-server
Amber1
   |--- amber-v2-prod
   |--- amber-v2-dev
```

The naming convention for future databases is [productname]-[version]-[deployment]. Amber v3 would have databases `amber-v3-prod` and `amber-v3-dev`.

You can't change database names after the fact, so `amber-prod` and `amber-dev` are base cases, I guess. :)

Collections
-----------

A database comprises some number of collections. A collection is a set of documents identified by unique `_id`s. Documents are JSON objects with a reserved field `"_id"` used by mongo to index the documents in the collection.

It is conventional to use a separate collection for every kind of document your application needs to store. Collections are lightweight and there is no limit on how many you can make. Efficiency starts decreasing around 1000.

Regarding indexing, it's important to understand indexes thoroughly by reading the Mongo docs. Without an index, query time is proportional to document count. A well-constructed index makes your query run in constant time as the collection grows.

You can use any data type in the `"_id"` field (not just native `ObjectIDs`) but I don't recommend it. We do that in Amber v1: `"_id"` is a 16-digit string, the actual sensor ID in the Amber application domain. I think this conflated the storage-layer identifier with the application-domain identifier... v2 doesn't do that.

Amber v1 has ~10 collections...

```
sensors <- top-level sensor info
images  <- CommonState image for each sensor
tenants <- each of our tenants
pretrain <- chunked buffers of pretraining data in flight
...
```

Amber v2 will have the core collections

```
tenant       <- Amber tenant (a top-level collection of resources), associated with a business entity
license      <- licenses (a child of one or more tenants, and parent of one or more models)
model        <- (i) the serialized image for a model and (ii) the CommonState attributes that are actually in the blob
user         <- A portal user, associated with a person and email address
role         <- Named collection of permissions
permission   <- Document granting access/denial to perform some operation on some resource

```

------------

* There is the inactive `amber-azure` and a `dashboard-deployment` space as well.
* Actually, Jim added `Amber1` today as the cluster for Amber V2.
