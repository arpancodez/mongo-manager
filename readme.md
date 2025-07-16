# MongoDB Docker Manager

A command-line interface (CLI) tool written in Go for managing a MongoDB container via Docker. This application provides an interactive menu to perform common operations such as starting/stopping the container, viewing logs, creating databases and users, and retrieving database information.
I built this application to make my life easy.

## Description

This tool simplifies the management of a MongoDB instance running in a Docker container. It utilizes the official `mongo:latest` Docker image and offers a user-friendly, colorized terminal menu. Key operations are handled programmatically, with built-in error checking and cross-platform compatibility.

**Prerequisites:**
- Go (version 1.16 or higher).
- Docker installed and running on your system.

## Features

- Interactive menu with ASCII art and colored output for better usability.
- Start, stop, and remove the MongoDB container (`my-mongodb`).
- View real-time container logs (interruptible with Ctrl+C).
- Create new databases.
- Add users with read/write roles to specific databases.
- Retrieve database information, including users and statistics (passwords are not displayed for security).
- Cross-platform support (Windows, macOS, Linux), including screen clearing.
- Graceful handling of errors and container states.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/inlovewithgo/mongo-manager.git
cd mongo-manager
```
2. Install dependencies:
```bash
go mod tidy
```
3. Build the executable:
```bash
go build -o mongodb-manager
```


This will produce an executable file named `mongodb-manager` (or `mongodb-manager.exe` on Windows).

## Usage

Ensure Docker is running, then execute the program:
```bash
./mongodb-manager
```

### The menu will display options:
- **1**: Start the MongoDB container (exposes port 27017).
- **2**: Stop and remove the container.
- **3**: View live logs (press Ctrl+C to return).
- **4**: Create a new database.
- **5**: Add a new user to an existing database.
- **6**: Retrieve information for all databases (users and stats).
- **7**: Exit.

After each operation, press Enter to return to the menu.


**Notes:**
- The container name is fixed as `my-mongodb`.
- System databases (`local`, `config`) are excluded from listings.
- Users are assigned the `readWrite` role on the specified database.

## Dependencies

- Go standard library.
- External: `github.com/fatih/color` (for colored terminal output).

Install via `go mod tidy`.

## Troubleshooting

- **Docker not running**: Start the Docker daemon and try again.
- **Container already exists**: Use option 2 to stop and remove it, or modify `containerName` in the source code.
- **Permission issues**: Run the executable with elevated privileges (e.g., `sudo ./mongodb-manager` on Linux/macOS).
- **No output or errors**: Ensure Docker is accessible from the command line and the MongoDB image is available.
- **GIF not displaying**: Upload `demo.gif` to your repository root and ensure the path is correct.

## Contributing

Contributions are welcome. Please fork the repository, create a feature branch, commit your changes, and submit a pull request.

## Credits

- Developed by Shubham.
- Built with Go, Docker, and the `fatih/color` library.

