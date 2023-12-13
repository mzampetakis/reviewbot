# ReviewBot 🤖
ReviewBot is a review chatbot that gathers customers' reviews after a completed purchase.

## Assumptions 🙌🏻

We have made some assumption at this project in order to simplify and accelerate its development. These assumptions 
are: 
- The sentiment analysis and the text generator have been implemented as dummy but 3rd party integrations can be added
- For the chat implementation we have used a websocket communication channel
- Review chat support only one connection (client) at a time
- Database is pre-populated with some dummy data for testing purposes

## Requirements ✅
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

## Project Structure 📏

Here you can find details about how this project is structured. This will help you decide where to put anything new
added.

| Folder    | Description                                    |
|-----------|------------------------------------------------|
| **`app`** | Contains various helper application entities.  |

| Folder             | Description                                                               |
|--------------------|---------------------------------------------------------------------------|
| **`cmd`**          | Contains the main applications of the project.                            |
| `↳ cmd/reviewbot/` | Contains the applications of the project alongside the main function.     |
| `↳ cmd/reviewbot/api`   | Contains the api applications of the project alongside the main function. |

| Folder                     | Description                                                                         |
|----------------------------|-------------------------------------------------------------------------------------|
| **`internal`**             | Contains various helper packages used by the application.                           |
| `↳ internal/database/`     | Contains the application's database connection and migration logic                  |
| `↳ internal/domain/`       | Contains the application's specific packages.                                       |
| `↳ internal/domain/orders` | Contains the application's orders service.                                          |
| `↳ internal/env`           | Contains functionality to retrieve the application's configuration through EnvVars. |
| `↳ internal/version`       | Contains functionality to retrieve the application's version through Git.           |


| Folder                          | Description                                                                                                      |
|---------------------------------|------------------------------------------------------------------------------------------------------------------|
| **`pkg`**                       | Contains various packages used by the application but can also be used as standalone libraries by other applications. |
| `↳ pkg/responsegenerator/`          | Contains the Response Generator functionality through interface.                                                 |
| `↳ pkg/sentimentanalyzer` | Contains the Sentiment Analyzer functionality through interface.                                                 |


## Contribute 🙋

Open an issue for discussing any issue or bug.  
Open a PR for adding a new feature implementation or fix.  
Use appropriate commands from makefile to ensure application correctness.  
For a complete list of makefile commands run:
```
make help
```

### Further Improvements

There are several improvements that can be done to this project. Here are some indicative:
- Add more tests as we cover just a basic functionality among all stacks
- Add an Open API specification file for the served API

