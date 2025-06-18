package main

import (
	_ "log/slog" // unused
)

var versions = []string{
	"5.4.8", // 5.4
	"5.3.6", // 5.3
	"5.2.4", // 5.2
	"5.1.5", // 5.1
}

// TODO: implement whole logic

func main() {
	for _, version := range versions {
		_ = version
	}
}
