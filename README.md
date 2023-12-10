# ReviewBot 🤖
ReviewBot is a review chatbot that gathers customers' reviews after a completed purchase.

## Requirements✅
- Go v1.20
- make
- Docker
- Docker-compose


## Getting started 💻

ReviewBot offers an API and chatbot functionality for gathering user review after a purchase. 

### Configuration

The application uses configuration through Environment variables. Here is a list with the details and the default
value for each one of them:

| EnvVar           | Description                          | Default Value      |
|------------------|--------------------------------------|--------------------|
| `BASE_URL`       | Base URL of the API server.          | "http://localhost" |
| `HTTP_PORT`      | Port used byt the API server.        | "4444"             |
| `DB_HOST`        | Database host to use for connection. | "myreviewbotdb"    |
| `DB_USER`        | User of the database.                | "user"             |
| `DB_PASSWORD`    | User's password in the database      | "pass"             |
| `DB_NAME`        | Database name.                       | "myreviewbot"      |
| `DB_PORT`        | Database port to use for connection. | "3306"             |
| `DB_AUTOMIGRATE` | Enable auto DB schema migration      | true               |

### Running the application through containers

In order to run the application using docker images use the following commands:
```
make start
```

When all containers have started you will be able to access the API through:
```
http://localhost:4444
```
unless the default configuration has been updated. 

To review the logs of all containers execute:
```
make logs
```

In order to stop all running containers execute:
```
make stop
```

In order to cleanup all created containers execute:
```
make clean
```

### Running the application locally

In order to build the `reviewbot` run:
```
make build
```

In order to run the `reviewbot` run:
```
make run
```
Logs stream is the stdout.


## Contribute 🙋

Open an issue for discussing any issue or bug.  
Open a PR for adding a new feature implementation or fix.  
Use appropriate commands from makefile to ensure application correctness.  
For a complete list of makefile commands run:
```
make help
```
