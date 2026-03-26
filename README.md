# Gin CRUD backend challenge

You have to build a microservice that exposes a REST api with two different
tables, users and companies. Both tables should be open to creation, deletion,
or update. Every request must only accept this `Content-type: application/json`.

## Table of contents

- [Features](#Features)
- [Tables](#Tables)
- [Badges](#Badges)
- [Technologies](#Technologies)
- [PreRequisites](#Pre-requisites)
- [Run APP](#Run-APP)
- [Routes](#Routes)
- [Run tests](#Run-tests)
- [CI](#CI)
- [Deployment](#Deployment)
- [Decisions made](#Decisions-made) 
- [Areas to improve](#Areas-to-improve)
- [Author](#Author)

## Features

- Create `users` and `companies` tables.
- CRUD REST API to interface with `users` table.
- CRUD REST API to interface with `companies` table.

## Tables

### Users


| Column     | Type                |
| ---------- | ------------------- |
| id         | int sequence ( PK ) |
| name       | string              |
| age        | int                 |
| company    | int (FK companies)  |
| updated_at | datetime            |
| created_at | datetime            |


### Companies


| Column     | Type                |
| ---------- | ------------------- |
| id         | int sequence ( PK ) |
| name       | string              |
| updated_at | datetime            |
| created_at | datetime            |


## Badges

[CircleCI](https://dl.circleci.com/status-badge/redirect/gh/DarioJejer/gin-proyect/tree/main)
[Coverage Status](https://coveralls.io/github/DarioJejer/gin-proyect?branch=main)

## Technologies

- Programming languaje: Go
- APP Framework: Gin
- DBMS: Postgres
- Containers: Docker-compose
- Deployment: Heroku
- CI: CircleCI

## Pre-requisites

- Docker and docker compose installed.
- Linux/Mac terminal (Or emulated linux on Windows)
- No services running on localhost port 5432 or 3000.

## Run APP

1. Create an .env file on base directory with values set per .env.example
2. Execute docker compose to create the db, run migrations, connect to it and run the app.

```
docker-compose up
```

1. Go to the swagger endpoint and test the app or consume api through postman.
2. To stop the app.

```
docker-compose down
```

## Routes

### Deployed on Heroku

- Swagger: [http://gin-proyect-39779df05d77.herokuapp.com/docs/index.html](http://gin-proyect-39779df05d77.herokuapp.com/docs/index.html)
- API users: [http://gin-proyect-39779df05d77.herokuapp.com/users/](http://gin-proyect-39779df05d77.herokuapp.com/users/)
- API companies: [http://gin-proyect-39779df05d77.herokuapp.com/companies/](http://gin-proyect-39779df05d77.herokuapp.com/companies/)

### Local

- Swagger: [http://localhost:3000/docs/index.html](http://localhost:3000/docs/index.html)
- API users: [http://localhost:3000/users/](http://localhost:3000/users/)
- API companies: [http://localhost:3000/companies/](http://localhost:3000/companies/)

## Run tests

On the terminal, from the main project directory, run:

```
docker compose up
go test ./... -v
```

## CI

We check in a remote environment on CircleCi that
the build is valid and that none of the tests are failing. If everything is
okay, then the code coverage is sent to coveralls and in that site the test
coverage can be reviewed in detail.

## Deployment

Done on Heroku because of it's ease of use and low cost and connecting to a DB hosted on Neon. 

We can access the app on the url: [http://gin-proyect-39779df05d77.herokuapp.com](http://gin-proyect-39779df05d77.herokuapp.com)

## Decisions made

- Clean architecture: to be able to make future changes cleanly.
- Gorm: It is the most popular Go ORM so it properly maintained by the comunity.
- Docker: to make the app portable and easy to deploy.
- CI: to automate the coverage of tests and build pre-requisite for deployment.
- E2E testing: to cover the use of repository and router settup.

## Areas to improve

- Error handling could be improved 
- Add seed to populate db
- Apply migration scheme to set db

## Author

Dario Jejer

- GitHub: [https://github.com/DarioJejer](https://github.com/DarioJejer)
- LinkedIn: [https://www.linkedin.com/in/dariojejer](https://www.linkedin.com/in/dariojejer)

