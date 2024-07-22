# chirpy
Chirpy is a social network similar to Twitter.

## Prerequisites

Ensure you have the following tools installed:
- [Air](https://github.com/cosmtrek/air) - Live reload for Go applications.
- `gofmt` - Go source code formatter (comes with the Go toolchain).
- [Docker Compose](https://docs.docker.com/compose/) - Define and run multi-container Docker applications.

## Makefile Commands

### `make help`

Displays the help message with the list of available commands.

```sh
make help
```

### `make watch`

Uses Air to watch for file changes and rebuild the Go project.

```sh
make watch
```

### `make lint`

Runs `gofmt` on all Go files to check and correct formatting.

```sh
make lint
```

### `make docker-watch`

Runs `docker-compose up --build` to build and start the Docker containers.

```sh
make docker-watch
```

### `make docker-build`

Runs `docker-compose build` to build the Docker containers.

```sh
make docker-build
```

### `make check-tools`

Ensures that `air` and `gofmt` are installed before running any commands.

```sh
make check-tools
```

## Usage

Clone the repository and navigate to the project directory:

```sh
git clone <repository-url>
cd <project-directory>
```

Run any of the make commands listed above to perform the respective actions. For example, to watch for file changes and rebuild the Go project:

```sh
make watch
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.