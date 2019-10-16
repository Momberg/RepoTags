# RepoTags

## Environment variables ##

For the application to run correctly you need to configure some environment variables in your OS.

Variables needed for the application (This example is been used in a MacOS):

```
export GITTOKEN={{Add here your github token}}
export MYSQLUSER=root
export MYSQLPASS={{Add here your mysql password}}
export MYSQLHOST=localhost:3306
```

## Usage ##

Imports needed for the application:

```
go get github.com/gorilla/mux
go get golang.org/x/oauth2
go get github.com/google/go-github/github
go get github.com/go-sql-driver/mysql
go get github.com/jmoiron/sqlx
```
## MYSQL Database ##

You need to import the gittags.sql file located in the project root folder.
