# Cheatsheat of useful stuff

Dependencies:
* Language:   [Golang](https://go.dev)
* ORM:        [GORM](https://gorm.io)
* Web:        [Mux](https://github.com/gorilla/mux)
* Hot reload: [Fresh](https://www.github.com/pilu/fresh) for hot reload
* Formatter:  [gofmt](https://pkg.go.dev/cmd/gofmt)
* Linter:     [golangci-lint](https://golangci-lint.run)

First install dependencies then you can proceed to run the rest api.
To run the project with hot-reload:
```terminal
make -f cmd/Makefile
```


## Precommit

* Formatter:  [gofmt](https://pkg.go.dev/cmd/gofmt)
* Linter:     [golangci-lint](https://golangci-lint.run)


## Commands

Format the whole project

```terminal
./fmt
```

Connect to the docker development database

```terminal
./db
```
