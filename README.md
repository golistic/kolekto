kolekto - Manage JSON Collections Efficiently
=============================================

Copyright (c) 2022, Geert JM Vanderkelen

Manages data as collections in data stores which support JSON
functionality using Go.


Disclaimer
----------

This package is created based on year long experience with similar
implementations in Go. However, Kolekto is still in early stages
of development, and offers solutions for the author's own projects.

Ideas are welcome. Pull Requests are great, but there is currently
no capacity to process them.

Use in production is not advised just yet.


Overview
--------

Kolekto offers functionality to store collections of JSON documents
using data stores provided they support JSON functionality.

Any struct instances of which their types embed a type and implement
an interface, can be relatively easy stored into a table.

For example, given the following struct:

    type Band struct {
        Name    string   `json:"name"`
        Members []string `json:"members"`
        Active  bool     `json:"active"`
    }

Can be made compatible for Kolekto as:

    type Band struct {
        kolektor.Model
        Name    string   `json:"name"`
        Members []string `json:"members"`
        Active  bool     `json:"active"`
    }

    func (b Band) CollectionName() string {
        return "bands"
    }

The above is part of the `example_test.go` file which shows a full example
on how to use Kolekto.

### Example

Please see the `example_test.go` for a full example.


Models, Collections, and Sessions
---------------------------------

Each struct that implements the `kolektor.Modeler` interface, and 
embeds `kolektor.Meta`, is considered a Model. An instance of a model
is called an Object.

A group of objects instantiated from a certain model is called a
Collection. When stored using a relational database management system
(RDBMS), this corresponds to an SQL table, which is created automatically.

A Session is a wrapper around a database pool, for example, Go's `*sql.DB`.
It is used to initialize the Collections based on Models.  
A Session manages all collections it can access.

There is usually one Session possibly globally available with which
multiple Collections are created to store object based on the collection's
Model.


Supported Data Stores
---------------------

Kolekto can be extended with more data stores.

Currently, Kolekto supports the following data stores:

| Data Store           | Pool           | Driver                                  |
|----------------------|----------------|-----------------------------------------|
| **MySQL** v8.0.29    | `sql.DB`       | https://github.com/go-sql-driver/mysql  |
| **PostgreSQL** v14.3 | `pgxpool.Pool` | https://github.com/jackc/pgx/v4/pgxpool |

Note: the version denotes what we use for testing. Previous versions might not work.

All the above data stores are available when compiling Kolekto. However,
if you need only MySQL or only PostgreSQL you can use build tags:

| Build Tag   | Effect                     |
|-------------|----------------------------|
| **nomysql** | disable MySQL support      |
| **nopgsql** | disable PostgreSQL support |


License
-------

Distributed under the MIT license. See LICENSE.txt for more information.
