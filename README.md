# dont-panic-api [WIP]
Work in progress - the current version doesn't have the basic features yet and shouldn't be used

This is an API written to support your studies in Competitive Programming.

Its is a tool to:

- List subjects with difficulty and dependencies
- List problems for each subject
- Suggest new subjects to learn
- Suggest new problems to solve
- Track subjects you've learned
- Track problems you've solved
- Suggest when its time to review some subject by recommending problems of that subject (based on spaced repetition method)
- Give permissions for coaches and teammates to track your progress

## Pre-requisites
This API uses Azure CosmosDB graph database. To use it you will need to configure these environment variables for accessing your db:

```
$env:CDB_USERNAME="/dbs/<collection name>/<graph name>"
$env:CDB_HOST="wss://<resource name>.gremlin.cosmos.azure.com"
$env:CDB_KEY="<resource access key>"
```

## Install and run
```
go build .
.\dont-panic-api
```

## Examples to test the API
Subjects CRUD

- `CREATE`
 It's a `POST` request with URL `localhost:8080/api/problems`. Remember to put a valid body. Example:
 ```json
 {
   "name": "xavi-and-the-cookies",
   "subjects": [
     "dynamic-programming",
     "math"
   ],
   "difficulty": 2
 }
 ```
  The expected result of this request is that the Problem will be created if there isn't already any other Problem with the name you specified.
 
- `READ`
 To fetch all problems, make a `GET` request with URL `localhost:8080/api/problems`. To get a specific Problem, make a GET request with URL `localhost:8080/api/problems/{name}`, being `{name}` the name of the Problem.
 The expected result is a `json` with the information from the problems you requested.
 
- `UPDATE`
 It's a `PUT` request with URL `localhost:8080/api/problems/{name}`, being `{name}` the name of the Problem. Put a valid body that represents the updated Problem.
 The expected result is that the Problem information will be updated if there exists a Problem in the database with the name you specified.
 
- `DELETE`
 It's a `DELETE` request with URL `localhost:8080/api/problems/{name}`, being `{name}` the name of the Problem.
 The expected result is that the Problem will be delete if there is a Problem in the database with the name you specified, or an error if there isn't.