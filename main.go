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

	"github.com/fatih/color"
)

const (
	containerName = "my-mongodb"
	imageName     = "mongo:latest"
)

func main() {
	for {
		printMenu()
		handleUserInput()
	}
}

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

func printMenu() {
	clearScreen()
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println(" ██████╗ ██████╗ ███╗   ██╗ ██████╗  ██████╗")
	cyan.Println("██╔════╝██╔═══██╗████╗  ██║██╔════╝ ██╔═══██╗")
	cyan.Println("██║     ██║   ██║██╔██╗ ██║██║  ███╗██║   ██║")
	cyan.Println("██║     ██║   ██║██║╚██╗██║██║   ██║██║   ██║")
	cyan.Println("╚██████╗╚██████╔╝██║ ╚████║╚██████╔╝╚██████╔╝")
	cyan.Println(" ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝  ╚═════╝ ")
	color.New(color.FgHiBlack).Println("             MongoDB Docker Manager by Shubham")
	fmt.Println()

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
	fmt.Printf("│ [%s] %-34s │\n", color.RedString("7"), "Exit")
	yellow.Println("└───────────────────────────────────────────┘")
}

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
		printSuccess("Exiting. Goodbye!")
		os.Exit(0)
	default:
		printError("Invalid choice. Please try again.")
	}

	color.New(color.FgYellow).Println("\nPress Enter to return to the menu...")
	reader.ReadString('\n')
}

func printSuccess(msg string) {
	color.New(color.FgGreen, color.Bold).Printf("\n✔ %s\n", msg)
}

func printError(msg string) {
	color.New(color.FgRed, color.Bold).Fprintf(os.Stderr, "\n✖ Error: %s\n", msg)
}

func printInfo(msg string) {
	color.New(color.FgBlue, color.Bold).Printf("\nℹ %s\n", msg)
}

func executeCommand(name string, args ...string) error {
	color.New(color.FgMagenta).Printf("▶ Executing: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeAndCaptureCommand(name string, args ...string) (string, error) {
	color.New(color.FgMagenta).Printf("▶ Executing: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func isContainerRunning() bool {
	cmd := exec.Command("docker", "ps", "-q", "-f", fmt.Sprintf("name=%s", containerName))
	output, err := cmd.Output()
	if err != nil {
		printError(fmt.Sprintf("Failed to check container status: %v", err))
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

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

func listDatabasesForUser() {
	printInfo("Fetching available databases...")
	listScript := `db.adminCommand({ listDatabases: 1 }).databases.forEach(db => { if(db.name !== 'local' && db.name !== 'config') print(db.name) });`
	dbList, err := executeAndCaptureCommand("docker", "exec", containerName, "mongosh", "--quiet", "--eval", listScript)
	if err != nil {
		printError(fmt.Sprintf("Could not fetch database list: %v", err))
		return
	}
	if strings.TrimSpace(dbList) == "" {
		color.New(color.FgYellow).Println("No user-created databases found yet. 'admin' is always available.")
	} else {
		color.New(color.FgYellow).Println("Available databases to add a user to:")
		fmt.Println(dbList)
	}
}

func addNewUser() {
	if !isContainerRunning() {
		printError("MongoDB container is not running. Please start it first.")
		return
	}

	listDatabasesForUser()
	reader := bufio.NewReader(os.Stdin)

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

	mongoCommand := fmt.Sprintf(`db.getSiblingDB('%s').createUser({ user: '%s', pwd: '%s', roles: [{ role: 'readWrite', db: '%s' }] })`, db, username, password, db)

	printInfo("Adding new user...")
	err := executeCommand("docker", "exec", "-i", containerName, "mongosh", "--quiet", "--eval", mongoCommand)
	if err != nil {
		printError(fmt.Sprintf("Failed to add user: %v", err))
		return
	}
	printSuccess(fmt.Sprintf("User '%s' added to database '%s' successfully.", username, db))
}

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

	mongoCommand := fmt.Sprintf(`db.getSiblingDB('%s').createCollection('initial_collection')`, dbName)

	printInfo(fmt.Sprintf("Creating database '%s'...", dbName))
	err := executeCommand("docker", "exec", "-i", containerName, "mongosh", "--quiet", "--eval", mongoCommand)
	if err != nil {
		printError(fmt.Sprintf("Failed to create database: %v", err))
		return
	}
	printSuccess(fmt.Sprintf("Database '%s' created successfully.", dbName))
}

func getDatabaseInfo() {
	if !isContainerRunning() {
		printError("MongoDB container is not running. Please start it first.")
		return
	}
	printInfo("Fetching information for all databases...")
	color.New(color.FgYellow).Println("Note: Passwords are encrypted and cannot be displayed.")

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
			if (users.length > 0) {
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
