package main

import (
	"log/slog"

	"github.com/oklog/ulid/v2"
)

var version = "unknown"

func main() {
	slog.Info(`Hello world!`, `version`, version)
	defer slog.Info(`Good bye!`)

	slog.Info(``, `ulid`, ulid.Make())
	slog.Info(``, `ulid`, ulid.Make())
	slog.Info(``, `ulid`, ulid.Make())
}
