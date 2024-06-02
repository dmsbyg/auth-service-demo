RESTful API service for the sign up & login. Written using Golang / Typescript.

- Must include any kind tests (unit test / postman test / integration test)
- Code must be uploaded to Github along with a README that has
- Details on the structure of the service
- Instructions on how to run the service
- [Bonus] Include deployment, linting etc to your code. Mention what you added in

## Overview

RESTful API service for authentication.
It only consist of two endpoint just for signup and login.
It provides JSON over HTTP protocol for communication from client.

In this demo it uses SQLite as database for the ease of deployment.
It might be replaced with other database by changing the implementation of database connection and repository layer.

Unit test is not written for logic part of the code. While the end-to-end test is done through Postman. 

## Code Structure

The code is structured as follow
.
├── cmd //cmd is entry points for application.
│   ├── api // our main api service entry point.
│   └── migrate // our migrator entry point.
├── config // configuration parsing and definition of this app goes here.
├── database
│   ├── database.go // database connections related codes goes here.
│   └── migrations // migrations script goes here, it uses golang-migrate to migrate the database.
├── Dockerfile // define what will be run to build this app's container.
├── internal // everything that is internal and specific to this application.
│  ├── auth // our auth module, if there is any other service / module that want to be added to this application, we can add in the same level as this 'auth' folder.
│  │   ├── interfaces.go // interfaces is the contract for this service.
│  │   ├── model.go // our business domain types goes here.
│  │   ├── repository.go // all about store/database layer should go here.
│  │   ├── reqresp.go // all types that define request and response of this API.
│  │   ├── service.go // the business logic layer of this application.
│  │   ├── transport.go // transport layer, in this case we use echo as HTTP router, goes here.
│  │   └── token // token related codes, it is separated from auth package since it we might have different token generation and verification method.
│  └── common // contains all common code regardless of which module/service it will be used.
├── Makefile  // all the scripts goes here.
├── pkg // all external packages that is depended by its interfaces, so we can change the implementation details in the future.
│  ├── logger 
│  └── validator 
└── utils // general utility functions goes here


## Configuration

This application is configurable through environment variable. You can use [example env file](./example.env) starting.

## How to Use
You can use [makefile](/Makefile) as your guide on what commands are available for this application.
First, you need an `.env` file located in the root of the project. To use the example environment you can run:
```bash
cp example.env .env
```

If you want to run this application on docker, you can build it by running:
```bash
make docker_build
```
and run it with:
```bash
make docker_run
```


## Tools

- [golang-migrate](https://github.com/golang-migrate/migrate): golang-migrate is used to generate migration script and do migration and rollback to the database schema.
- [air](https://github.com/cosmtrek/air): air is used to hot reloading the application. It is a nice to have package, but it makes the development experience better.
