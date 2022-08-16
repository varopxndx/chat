# Chat

## Architecture

[Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) consists on separate the logic into different packages. The main proposal is the _Usecase_ package.
The Usecase package is used to handle the _bussines logic_, and call other packages according to the different flows.

Clean Architecture proposes a lot of packages, the Chat app contains the principal ones:

* Model
* Service
* Usecase
* Router
* Handler

Also, it uses other packages to handle specific logic, like _Rabbit_ and _Token_

### Model

Model package contains all the structures to handle the JSON data and responses.
Is the only package that could be used in all the packages.

### Service

Service package (also called Repository) contains all the logic to handle the operations. This logic is separated into several functions, which allows to test in an easier way, have decoupled logic, and maintain the code in a better way.
If some logic is added in the future, a new function will be created and the Usecase will call it.

### Usecase

Usecase package contains the business logic, according to the flows, it will call the specific function in the Service package to get or store data.

### Handler

Handler package contains the controllers, it call the usecase flows, return the responses and handle the http errors.

### Router

Router package contains the endpoints, each enpoint is linked to a handler.
It also interacts with the Middleware to handle the Authorization process.

### Rabbit

Rabbit package contains the RabbitMQ logic.

## Requirements

* [Go](https://golang.org/doc/install) - Version _go1.15.2_ or above.
* [Docker](https://docs.docker.com/get-docker/)

## Setup

1. Clone the [repository](https://github.com/varopxndx/chat)

1. Open the terminal and go to the project path

1. In the root level, run the following command:

    ```sh
        docker-compose up
    ```

    > Note: The Chat app will start, it may take some time to be up since RabbitMQ takes much time to start.

1. The Chat app will be running on:

    ```sh
        http://localhost:8080/v1
    ```

## Running the tests

Because of the time, I only added tests for the _Handler_ and _Usecase_ packages.

* Run all the tests with the following command:

    ```go
        go test ./...  -v
    ```

## Achievements

* All Mandatory Features finished.
* All Bonus finished.
  * Two chat rooms supported.

## Future Improvements

* Add more unit tests.
* Allow users to create more chat rooms.
* Fix the commands messages (`/stock`), they appear in both rooms if both are open.
* Add an error message view when creating duplicated users.
* Add integration testing.
* Create a Makefile.
