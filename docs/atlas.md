Hi team,

Here are some comments on our current Atlas resources. I guess it is also a mini tutorial in mongo. I have some recommendations for conventions moving forward which you may follow to the extent they prove helpful.

Our tenancy
-----------

`boonlogic` is the name of our tenancy within Atlas. Within Atlas, a tenancy has a set of Project Spaces. boonlogic has one Project Space, `amber-aws`.*

```
boonlogic          <- tenancy
   |--- amber-aws  <- project space
```

If we use Atlas for a non-Amber product in the future, that should go under a new Project Space. Different versions of the same product should go in the same Project Space.

Clusters
--------

A Project Space contains one or more clusters. `amber-aws` has two clusters, `Amber0` and `Amber1`.

```
amber-aws
   |--- Amber0     <- Amber runs on this cluster right now.
   |--- Amber1     <- Amber v2 runs on this cluster.
```

A cluster is a set of synchronized database replicates running on physically co-located compute. It is associated with a cloud provider and one of that provider's resource regions. `Amber0` is located in AWS region `us-east-1`. If we deployed Amber to `us-west-1`, it would need to be on a new cluster in that region.

The naming convention for clusters should be [Product][n] where n simply increments each time we add another cluster. Let's say future versions of Buc and AVIS use Atlas.

```
Amber0 <- Amber runs on this cluster right now.
Amber1 <- Amber v2 runs on this cluster
AVIS2  <- AVIS v10 places data on this cluster
Buc3   <- Buc data goes here instead of a physical safe
etc.
```

This allows the esteemed database janitor to recognize at a glance both the associated project and unique identifier for the cluster within growing list of Atlas cluster deployments.

I enabled Termination Protection which **prevents any user from accidentally deleting a cluster**. To delete a cluster, Termination Protection will have to be disabled by a superuser first.

Databases
---------

A cluster contains a number of databases. A database is the top-level container of all application data for one deployment of an application. `Amber0` has three databases.
```
Amber0
   |--- amber-prod   <- Amber production
   |--- amber-dev    <- Amber development environment
   |--- aop-license  <- license information for amber-on-prem
```

The esteemed janitor really just needs to keep `amber-dev` usable enough for development and guard `amber-prod` ferociously.

Jim created cluster `Amber1` earlier this week and the databases for the Amber v2 prod/dev deployments. This makes:

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

You can't change database names after the fact, so `amber-prod` and `amber-dev` are just going to be exceptions (maybe viewed more favorably as base cases). I think you could migrate `amber-prod` to `amber-v2-prod` by getting to where you have two databases with the two names and which stay synchronized somehow, then throwing a switch that directs Lambda instances to the new database instead.

Collections
-----------

A database comprises some number of collections. A collection is a set of documents identified by unique `_id`s. Documents are JSON objects with a reserved field `"_id"` used by mongo. the `"_id"` field is indexed by default and unique with respect to all documents within the collection.

It is conventional to use a separate collection for each kind of document your application needs to store. Collections are lightweight and there is no limit on how many you can make. I have read on StackOverflow that you can reach about ~1000 collections before it starts impacting performance.

Regarding indexing, it's important to understand indexes thoroughly by reading the Mongo docs. Generally speaking, an unindexed query is a linear search whose cost becomes proportional to collection size. A proper index will make your query run in constant time as the collection grows. A second utility of indexes is to enforce uniqueness on a field.

You can use any data type in the `"_id"` field (not just native `ObjectIDs`) but I don't recommend it. I did that in Amber v1: `"_id"` is a 16-digit string and is the actual sensor ID in the Amber application domain. This conflated a storage-layer identifier with an application-domain identifier.

---

It has been a great honor to be your Database Custodian.

-LA
