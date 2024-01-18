# Golang Web App for CF

A very basic web application that is Cloud Foundry service binding aware. This app is primarily for testing out functionality in CF and is faster to deploy than an equivalent Java application.

## Usage

```bash
$ GOOS=linux GOARCH=amd64 go build .
$ cf push
```
