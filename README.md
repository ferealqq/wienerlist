# Cheatsheat of useful stuff

Dependencies:
* Language:   [golang](https://go.dev)
* ORM:        [GORM](https://gorm.io)
* Web:        [gin](https://gin-gonic.com/docs/)
* Hot reload: [fresh](https://www.github.com/pilu/fresh) 
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
