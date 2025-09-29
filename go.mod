// Module path (import path for your project)
module github.com/inlovewithgo/mongo-manager

// Go version used for this project
go 1.23.5

require (
	// Library for colored terminal output
	github.com/fatih/color v1.18.0 // indirect

	// Cross-platform library to handle colored output (needed by fatih/color on Windows)
	github.com/mattn/go-colorable v0.1.13 // indirect

	// Detects whether output is a terminal or not (needed for colored output)
	github.com/mattn/go-isatty v0.0.20 // indirect

	// Low-level system calls (used by other libraries like isatty/colorable)
	golang.org/x/sys v0.25.0 // indirect
)
