# dont-panic-api [WIP]
Work in progress - the current version doesn't have the basic features yet and shouldn't be used

This is an api written to support your studies in competitive programming.
It intends to be a way to track the problems you've solve, the subjects you've learned and suggests new subjects to learn and when its time to review any of them by suggesting problems with these subjects.
It also intends to help coachs and teammates to track the development of each other.

## pre-requisites
This api uses azure cosmosdb graph database. To use it you will need to configure these environment variables for accessing your db:

```
$env:CDB_USERNAME="/dbs/<collection name>/<graph name>"
$env:CDB_HOST="wss://<resource name>.gremlin.cosmos.azure.com"
$env:CDB_KEY="<resource access key>"
```

## install and run
```
go build .
.\dont-panic-api
```