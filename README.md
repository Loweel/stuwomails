StuWoMails
==========

Installation
------------
Install [Go](https://golang.org/doc/install). Then, run this:

    go get github.com/srhnsn/stuwomails/...

Copy `resources/config.example.yaml` to `resources/config.development.yaml` and `resources/config.production.yaml` and adjust accordingly. Install the database schema which is in `resources/dump.sql`.

To create a standalone x64 Linux binary, use `go run build.go production`. In this case, `resources/config.production.yaml` will be used as configuration files. All other `build.go` commands use `resources/config.development.yaml` as configuration file. For example, it will be used when you start the application locally: `go run build.go run`. See `build.go` for more commands.


Code overview
-------------
`app/main.go` is the entry point for the program. The `assets` directory contains any needed files that will be bundled with the resulting binary.
