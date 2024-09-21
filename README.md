# wallet

## Description
This is a simple application meant to mimic a web-based banking application. It allows users to transfer 
funds between accounts, view their transaction history, and view their account balance. Its primary purpose is to prove
that I can write a full-stack application using a variety of technologies, from the platform level, to the back-end,
to the front-end in an idiomatic way. Its secondary purpose is to allow me to learn more about the technologies that
I chose to employ for the application, namely Go for the front-end, MongoDB for the database, Angular for
the front-end, and Docker (and eventually Kubernetes) for deployment.

## Installation

### Back-end
As previously mentioned, the back-end is written in Go. To install the back-end, you will need to have Go installed
on your machine. You can download Go from the official website [here](https://go.dev/doc/install). Currently, the
repository uses version 1.22.1. Once you have Go installed, you can clone the repository and navigate to the
directory ./go_webserver. From there you can either build the migrator app or the webserver app. The migrator
app used to create the database schema and seed the database with some initial data. The webserver app is the actual
webserver that serves the front-end and handles the API requests. To build the migrator app, you can run the following
command 
```bash
go build -o ${DESIRED_OUTPUT_NAME} ./cmd/migrator/
```
where ${DESIRED_OUTPUT_NAME} is the name of the output file. To build the webserver app, you can run the following
command
```bash
go build -o ${DESIRED_OUTPUT_NAME} ./cmd/webserver/
```
where ${DESIRED_OUTPUT_NAME} is the name of the output file. Once you have built the desired app, you can run it
by executing the output file. Both the webserver and migrator apps require a MongoDB instance to be 
running on the default port (27017).
