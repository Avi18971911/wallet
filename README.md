# Wallet

## Description
This is a simple application meant to mimic a web-based banking application. It allows users to transfer 
funds between accounts, view their transaction history, and view their account balance. Its primary purpose is to prove
that I can write a full-stack application using a variety of technologies, from the platform level, to the back-end,
to the front-end in an idiomatic way. Its secondary purpose is to allow me to learn more about the technologies that
I chose to employ for the application, namely Go for the front-end, MongoDB for the database, Angular for
the front-end, and Docker (and eventually Kubernetes) for deployment.

## Installation

### Back-end

#### Prerequisites
Make sure the following are installed on your machine:
- **Go**: Download from [Go website](https://go.dev/doc/install) (This repo uses `Go v1.22.1`)
- **Docker**: Download from [Docker website](https://docs.docker.com/get-docker/) (This repo uses `Docker 25.0.3`)
- **Docker Compose**: Download from [Docker Compose website](https://docs.docker.com/compose/install/) (This repo uses `Docker Compose v2.24.6-desktop.1`)

Optional:
- **MongoDB CLI**: Download from [MongoDB website](https://www.mongodb.com/try/download/shell) (This repo uses `MongoDB Shell 2.2.3`)

#### Steps
1. Install the dependencies. Note that all commands are run from the root of the repository.
    ```bash
    cd ./go_webserver
    go mod tidy
    ```
2. Build either the migrator or webserver app (Optional).

The migrator
app used to create the database schema and seed the database with some initial data. The webserver app is the actual
webserver that serves the front-end and handles the API requests. To build the migrator app, you can run the following
command.
```bash
cd ./go_webserver
go build -o ${DESIRED_OUTPUT_NAME} ./cmd/migrator/
```
where `DESIRED_OUTPUT_NAME` is the name of the output file. To build the webserver app, you can run the following
command
```bash
cd ./go_webserver
go build -o ${DESIRED_OUTPUT_NAME} ./cmd/webserver/
```
Once you have built the desired app, you can run it
by executing the output file. Both the webserver and migrator apps require a MongoDB instance to be 
running on the default port (`30001`), or you can specify a different port by setting the `MONGO_URL` environment 
variable accordingly: `mongodb://localhost:${PORT}`, where `PORT` is the desired port. 

3. Run the DB, webserver, and migrator using Docker Compose
```bash
cd ./go_webserver
docker-compose up --build
```
Once this has been executed, you can peek into the database using the MongoDB CLI using the following command
```bash
mongosh "mongodb://mongo:30001/?replicaSet=rs0"
```
Note that `docker-compose` creates a MongoDB instance running on port `30001` and a webserver running on port `8080`.
The database is named wallet and has two collections: `account` and `transaction`. The `account` collection contains
the account information for each user, and the `transaction` collection contains the transaction history for each user.
There will be no transactions in the `transaction` collection until you have made some transactions using the front-end.

4. Interacting with the API (Optional)

You can also send POST and GET requests to the API using endpoints detailed in `go_webserver/docs/swagger.json`.
Note that docker-compose exposes the webserver on port `8080`, so you can send requests to the API using the following
base URL: http://localhost:8080.

5. Creating the Swagger JSON (Optional)

To generate the swagger.json and swagger.yaml files, you can run the following command
```bash
cd ./go_webserver
swag init -g ./cmd/webserver/main.go
```
These will be consumed by the front-end to generate the API functions.

### Front-end

#### Prequisites
Make sure the following are installed on your machine:
- **Node.js**: Download from [Node.js website](https://nodejs.org/en/download/) (This repo uses `Node v20.13.1`)
- **npm**: Comes with Node.js (This repo uses `npm v10.5.2`)
- **Angular CLI**: Install globally with npm
  ```bash
  npm install -g @angular/cli
- **OpenAPI Generator CLI**: Install globally with npm
  ```bash
  npm install -g @openapitools/openapi-generator-cli
  ```

#### Steps

1. Install the dependencies
    ```bash
    cd ./angular_frontend
    npm install
    ```
2. Generate the API functions

    ```bash
    cd ./
    openapi-generator-cli generate -i go_webserver/docs/swagger.json -g typescript-angular -o angular_frontend/projects/account-app/src/app/backend-api
    ```
   
    This will generate the API functions in the `angular_frontend/projects/account-app/src/app/backend-api` directory.


3. Start the Angular development server
    ```bash
    cd ./angular_frontend
    ng serve
    ```
    This will start the Angular development server on port `4200`. You can navigate to http://localhost:4200 to view the
    front-end. The front-end will be able to communicate with the back-end running on port `8080`. Note that this is
    hardcoded in the `proxy.conf.json` file in the `angular_frontend` directory. You can login using the credentials
    `username: Hilda`, `password: Hilda`, or you can probe the DB to find other users.

Note that the front-end only has three working components: the login page, the dashboard (landing page after login)
and the transfer to other walletbank accounts page. The transaction history page is not yet implemented.