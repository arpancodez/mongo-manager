package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/fatih/color" // Used for colorful CLI output
)

const (
	containerName = "my-mongodb"  // Default container name
	imageName     = "mongo:latest" // MongoDB Docker image
)

func main() {
	// Infinite loop to show menu until user exits
	for {
		printMenu()
		handleUserInput()
	}
}

// clearScreen clears the terminal screen depending on the OS.
func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// printMenu displays the main menu with options.
func printMenu() {
	clearScreen()

	// ASCII banner
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println(" ██████╗ ██████╗ ███╗   ██╗ ██████╗  ██████╗")
	cyan.Println("██╔════╝██╔═══██╗████╗  ██║██╔════╝ ██╔═══██╗")
	cyan.Println("██║     ██║   ██║██╔██╗ ██║██║  ███╗██║   ██║")
	cyan.Println("██║     ██║   ██║██║╚██╗██║██║   ██║██║   ██║")
	cyan.Println("╚██████╗╚██████╔╝██║ ╚████║╚██████╔╝╚██████╔╝")
	cyan.Println(" ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝  ╚═════╝ ")
	color.New(color.FgHiBlack).Println("             MongoDB Docker Manager by Shubham")
	fmt.Println()

	// Menu box
	yellow := color.New(color.FgYellow)
	yellow.Println("┌───────────────────────────────────────────┐")
	yellow.Println("│ Please choose an option:                  │")
	yellow.Println("├───────────────────────────────────────────┤")
	fmt.Printf("│ [%s] %-34s │\n", color.GreenString("1"), "Start MongoDB Container")
	fmt.Printf("│ [%s] %-34s │\n", color.GreenString("2"), "Stop MongoDB Container")
	fmt.Printf("│ [%s] %-34s │\n", color.GreenString("3"), "View Live Logs")
	fmt.Printf("│ [%s] %-34s │\n", color.GreenString("4"), "Add New Database")
	fmt.Printf("│ [%s] %-34s │\n", color.GreenString("5"), "Add New User")
	fmt.Printf("│ [%s] %-34s │\n", color.GreenString("6"), "Get Database Info")
	fmt.Printf("│ [%s] %-34s │\n", color.GreenString("7"), "Get Connection URI")
	fmt.Printf("│ [%s] %-34s │\n", color.RedString("8"), "Exit")
	yellow.Println("└───────────────────────────────────────────┘")
}

// handleUserInput reads user input and performs corresponding action.
func handleUserInput() {
	color.New(color.FgWhite, color.Bold).Print("\n▶ Enter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		startMongoDB()
	case "2":
		stopMongoDB()
	case "3":
		viewLogs()
	case "4":
		addNewDatabase()
	case "5":
		addNewUser()
	case "6":
		getDatabaseInfo()
	case "7":
		getConnectionURI()
	case "8":
		printSuccess("Exiting. Goodbye!")
		os.Exit(0)
	default:
		printError("Invalid choice. Please try again.")
	}

	// Pause before returning to menu
	color.New(color.FgYellow).Println("\nPress Enter to return to the menu...")
	reader.ReadString('\n')
}

// Utility functions for consistent colored messages
func printSuccess(msg string) {
	color.New(color.FgGreen, color.Bold).Printf("\n✔ %s\n", msg)
}
func printError(msg string) {
	color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "\n✖ Error: %s\n", msg)
}
func printInfo(msg string) {
	color.New(color.FgBlue, color.Bold).Printf("\nℹ %s\n", msg)
}

// executeCommand runs a shell command and streams output live.
func executeCommand(name string, args ...string) error {
	color.New(color.FgMagenta).Printf("▶ Executing: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// executeAndCaptureCommand runs a shell command and returns output.
func executeAndCaptureCommand(name string, args ...string) (string, error) {
	color.New(color.FgMagenta).Printf("▶ Executing: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

// isContainerRunning checks if the MongoDB container is running.
func isContainerRunning() bool {
	cmd := exec.Command("docker", "ps", "-q", "-f", fmt.Sprintf("name=%s", containerName))
	output, err := cmd.Output()
	if err != nil {
		printError(fmt.Sprintf("Failed to check container status: %v", err))
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

// startMongoDB starts the MongoDB Docker container.
func startMongoDB() {
	if isContainerRunning() {
		printInfo(fmt.Sprintf("Container '%s' is already running.", containerName))
		return
	}
	printInfo(fmt.Sprintf("Starting MongoDB container named '%s'...", containerName))
	err := executeCommand("docker", "run", "-d", "--name", containerName, "-p", "27017:27017", imageName)
	if err != nil {
		printError(fmt.Sprintf("Failed to start MongoDB container: %v", err))
		return
	}
	printSuccess("MongoDB container started successfully.")
}

// stopMongoDB stops and removes the MongoDB container.
func stopMongoDB() {
	cmd := exec.Command("docker", "ps", "-a", "-q", "-f", fmt.Sprintf("name=%s", containerName))
	output, err := cmd.Output()
	if err != nil || len(strings.TrimSpace(string(output))) == 0 {
		printInfo(fmt.Sprintf("Container '%s' not found.", containerName))
		return
	}

	printInfo(fmt.Sprintf("Stopping and removing container '%s'...", containerName))
	if err := executeCommand("docker", "stop", containerName); err != nil {
		color.New(color.FgYellow).Printf("! Warning: Failed to stop container, it might already be stopped. %v\n", err)
	}
	if err := executeCommand("docker", "rm", containerName); err != nil {
		printError(fmt.Sprintf("Failed to remove container: %v", err))
		return
	}
	printSuccess("MongoDB container stopped and removed successfully.")
}

// viewLogs shows live container logs until user presses Ctrl+C.
func viewLogs() {
	if !isContainerRunning() {
		printError("MongoDB container is not running. Please start it first.")
		return
	}
	clearScreen()
	printInfo("Showing live logs... Press Ctrl+C to stop and return to the menu.")

	cmd := exec.Command("docker", "logs", "-f", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Handle Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	cmdDone := make(chan bool, 1)

	go func() {
		cmd.Run()
		cmdDone <- true
	}()

	select {
	case <-sigs:
		if err := cmd.Process.Kill(); err != nil {
			printError(fmt.Sprintf("Failed to stop log process: %v", err))
		}
		printInfo("\nLog streaming stopped.")
	case <-cmdDone:
		printInfo("\nLog stream ended because the container stopped.")
	}
}

// listDatabasesForUser lists all user-created databases.
func listDatabasesForUser() (string, error) {
	printInfo("Fetching available databases...")
	listScript := `db.adminCommand({ listDatabases: 1 }).databases.forEach(db => { if(db.name !== 'local' && db.name !== 'config') print(db.name) });`
	dbList, err := executeAndCaptureCommand("docker", "exec", containerName, "mongosh", "--quiet", "--eval", listScript)
	if err != nil {
		printError(fmt.Sprintf("Could not fetch database list: %v", err))
		return "", err
	}
	if strings.TrimSpace(dbList) == "" {
		color.New(color.FgYellow).Println("No user-created databases found yet. 'admin' is always available.")
	} else {
		color.New(color.FgYellow).Println("Available databases:")
		fmt.Println(dbList)
	}
	return dbList, nil
}

// addNewUser adds a new MongoDB user with readWrite role on a given DB.
func addNewUser() {
	if !isContainerRunning() {
		printError("MongoDB container is not running. Please start it first.")
		return
	}

	listDatabasesForUser()
	reader := bufio.NewReader(os.Stdin)

	// Ask for username, password, database
	fmt.Println()
	color.New(color.FgHiWhite).Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	color.New(color.FgHiWhite).Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	color.New(color.FgHiWhite).Print("Enter database from the list above to assign the user to: ")
	db, _ := reader.ReadString('\n')
	db = strings.TrimSpace(db)

	if username == "" || password == "" || db == "" {
		printError("Username, password, and database cannot be empty.")
		return
	}

	// MongoDB command for user creation
	mongoCommand := fmt.Sprintf(`db.createUser({ user: '%s', pwd: '%s', roles: [{ role: 'readWrite', db: '%s' }] })`, username, password, db)

	printInfo("Adding new user...")
	err := executeCommand("docker", "exec", "-i", containerName, "mongosh", db, "--quiet", "--eval", mongoCommand)
	if err != nil {
		printError(fmt.Sprintf("Failed to add user: %v", err))
		return
	}
	printSuccess(fmt.Sprintf("User '%s' added to database '%s' successfully.", username, db))
}

// addNewDatabase creates a new MongoDB database by making an initial collection.
func addNewDatabase() {
	if !isContainerRunning() {
		printError("MongoDB container is not running. Please start it first.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	color.New(color.FgHiWhite).Print("Enter new database name: ")
	dbName, _ := reader.ReadString('\n')
	dbName = strings.TrimSpace(dbName)

	if dbName == "" {
		printError("Database name cannot be empty.")
		return
	}

	mongoCommand := `db.createCollection('initial_collection')`

	printInfo(fmt.Sprintf("Creating database '%s'...", dbName))
	err := executeCommand("docker", "exec", "-i", containerName, "mongosh", dbName, "--quiet", "--eval", mongoCommand)
	if err != nil {
		printError(fmt.Sprintf("Failed to create database: %v", err))
		return
	}
	printSuccess(fmt.Sprintf("Database '%s' created successfully.", dbName))
}

// getDatabaseInfo shows stats and users of all non-system databases.
func getDatabaseInfo() {
	if !isContainerRunning() {
		printError("MongoDB container is not running. Please start it first.")
		return
	}
	printInfo("Fetching information for all databases...")

	// Script to print DB stats and users
	mongoScript := `
		const dbs = db.adminCommand({ listDatabases: 1 }).databases;
		dbs.forEach(dbInfo => {
			if (dbInfo.name === 'local' || dbInfo.name === 'config') return;
			const currentDb = db.getSiblingDB(dbInfo.name);
			
			print("\n========================================================");
			print("DATABASE: " + dbInfo.name);
			print("========================================================");

			print("\n--- USERS ---");
			const users = currentDb.getUsers();
			if (users && users.length > 0) {
				printjson(users);
			} else {
				print("No users found in this database.");
			}

			print("\n--- STATS ---");
			const stats = currentDb.stats();
			printjson(stats);
		});
	`

	err := executeCommand("docker", "exec", containerName, "mongosh", "--quiet", "--eval", mongoScript)
	if err != nil {
		printError(fmt.Sprintf("Failed to retrieve database info: %v", err))
		return
	}
	printSuccess("Finished retrieving database information.")
}

// getConnectionURI builds a MongoDB connection URI from user input.
func getConnectionURI() {
	if !isContainerRunning() {
		printError("MongoDB container is not running. Please start it first.")
		return
	}

	printInfo("Fetching all users...")
	listUsersScript := `db.getSiblingDB('admin').system.users.find().forEach(u => print(u.user));`
	userList, err := executeAndCaptureCommand("docker", "exec", containerName, "mongosh", "--quiet", "--eval", listUsersScript)
	if err != nil {
		printError(fmt.Sprintf("Could not fetch user list: %v", err))
		return
	}
	if strings.TrimSpace(userList) == "" {
		printError("No users found. Please add a user first.")
		return
	}
	color.New(color.FgYellow).Println("Available users:")
	fmt.Println(userList)

	// Show available DBs
	listDatabasesForUser()

	reader := bufio.NewReader(os.Stdin)

	// Ask user details
	color.New(color.FgHiWhite).Print("\nEnter user from the list above: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)

	color.New(color.FgHiWhite).Print("Enter password for user '" + user + "': ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	color.New(color.FgHiWhite).Print("Enter database to connect to: ")
	db, _ := reader.ReadString('\n')
	db = strings.TrimSpace(db)

	if user == "" || password == "" || db == "" {
		printError("User, password, and database cannot be empty.")
		return
	}

	// Build connection string
	uri := fmt.Sprintf("mongodb://%s:%s@localhost:27017/%s", user, password, db)
	printSuccess("Your connection URI is:")
	color.New(color.FgCyan, color.Bold).Println(uri)
}
