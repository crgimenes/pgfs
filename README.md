# pgfs

**pgfs** mounts a Postgresql database as a read-only file system. In the current version is just a proof of concept, keep in mind that large tables will degrade the performance very quickly.

The purpose is to be able to use standard UNIX tools to access data such as grep, sort and others instead of SQL and so be more comfortable for the experienced command line user to inspect the database. In a development environment with very little data it should work without problems.

## how to use

- Create a pgfs.toml file with database credentials. you can sweat the pgfs.toml.sample file as an example

- Mount the file system with the following command

```console
go run main.go -m=mountpoint
```

Where *mountpoint* is the directory where you want to mount the database
