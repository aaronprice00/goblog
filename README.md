# GoBlog API

## Goal

Investigate code structure in golang by building a sample CRUD API with very few 3rd party libs

Includes: postgres, gorm, authentication, testing, docker, and docker-compose

## Install

copy .env.example to .env, open and edit to give your Postgres port/user/password/dbname

Warning: I recommend creating a new database as there are some tables drops / seeding during early execution.

(or you can run as docker-compose)

## Execute

```markdown
go run main.go

(or docker-compose) $ docker-compose up
```

## Testing

```markdown
cd ./tests && go test -v ./...
```

## Todo

Let me know if you think of something, this is just a throw away education project.

I plan to make 2 more projects looking at:

Domain Driven Design

Microservices

## Credit

Victor Steven - <https://levelup.gitconnected.com/crud-restful-api-with-go-gorm-jwt-postgres-mysql-and-testing-460a85ab7121>
little is changed from his project, again this is for my newbie research
