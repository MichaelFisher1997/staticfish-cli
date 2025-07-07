# Gemini Workspace

This file outlines the project structure and conventions for the `staticfish-cli` tool, intended for use by the Gemini agent.

## Project Structure

- `main.go`: The entry point for the CLI application. It uses the `cobra` library to define commands and flags.
- `google-search/`: This directory contains the logic for the `google-search` command.
  - `google-search.go`:  Defines the `Search` function, which encapsulates the core functionality of the `google-search` command.
- `go.mod`, `go.sum`: Go module files that manage project dependencies.

## Conventions

- **Commands:** New commands should be added to the `main.go` file, following the `cobra` library's conventions.
- **Functionality:** Core logic for new commands should be encapsulated in their own packages, similar to the `google-search` package.
- **Dependencies:** Project dependencies are managed using Go modules. Any new dependencies should be added using `go get`.
