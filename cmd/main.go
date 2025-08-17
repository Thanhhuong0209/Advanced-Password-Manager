package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"password-manager/internal/crypto"
	"password-manager/internal/generator"
	"password-manager/internal/storage"

	"golang.org/x/term"
)

const (
	appName = "Advanced Password Manager"
	version = "1.0.0"
)

var (
	dbPath         string
	masterPassword string
	database       *storage.Database
)

func main() {
	// Set default database path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	dbPath = filepath.Join(homeDir, ".password-manager", "passwords.db")

	// Parse command line arguments
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	// Initialize database connection
	if err := initializeDatabase(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Handle commands
	command := os.Args[1]
	switch command {
	case "generate", "gen":
		handleGenerate()
	case "save":
		handleSave()
	case "get", "find":
		handleGet()
	case "list":
		handleList()
	case "delete", "del":
		handleDelete()
	case "search":
		handleSearch()
	case "stats":
		handleStats()
	case "analyze":
		handleAnalyze()
	case "help", "-h", "--help":
		showHelp()
	case "version", "-v", "--version":
		showVersion()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

// initializeDatabase initializes the database connection
func initializeDatabase() error {
	// Get master password
	fmt.Print("Enter master password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // New line after password input
	
	masterPassword = string(bytePassword)
	if masterPassword == "" {
		return fmt.Errorf("master password cannot be empty")
	}

	// Create database
	database, err = storage.NewDatabase(dbPath, masterPassword)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

// handleGenerate handles password generation
func handleGenerate() {
	config := generator.DefaultConfig()
	
	// Parse flags
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case strings.HasPrefix(arg, "--length="):
			if length, err := strconv.Atoi(strings.TrimPrefix(arg, "--length=")); err == nil {
				config.Length = length
			}
		case arg == "--uppercase":
			config.Uppercase = true
		case arg == "--lowercase":
			config.Lowercase = true
		case arg == "--numbers":
			config.Numbers = true
		case arg == "--symbols":
			config.Symbols = true
		case arg == "--no-repeating":
			config.NoRepeating = true
		case strings.HasPrefix(arg, "--exclude="):
			config.Exclude = strings.TrimPrefix(arg, "--exclude=")
		}
	}

	// Generate password
	password, err := generator.GeneratePassword(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated password: %s\n", password)
	
	// Analyze strength
	analysis := generator.AnalyzePasswordStrength(password)
	fmt.Printf("Strength: %s (Score: %d/7)\n", 
		analysis["strength_level"], analysis["strength_score"])
}

// handleSave handles saving a password
func handleSave() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s save <name> [--username <username>] [--password <password>] [--url <url>] [--notes <notes>] [--tags <tag1,tag2>]\n", os.Args[0])
		os.Exit(1)
	}

	entry := &storage.PasswordEntry{
		Name: os.Args[2],
	}

	// Parse optional flags
	for i := 3; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--username" && i+1 < len(os.Args):
			entry.Username = os.Args[i+1]
			i++
		case arg == "--password" && i+1 < len(os.Args):
			entry.Password = os.Args[i+1]
			i++
		case arg == "--url" && i+1 < len(os.Args):
			entry.URL = os.Args[i+1]
			i++
		case arg == "--notes" && i+1 < len(os.Args):
			entry.Notes = os.Args[i+1]
			i++
		case arg == "--tags" && i+1 < len(os.Args):
			tags := strings.Split(os.Args[i+1], ",")
			for j, tag := range tags {
				tags[j] = strings.TrimSpace(tag)
			}
			entry.Tags = tags
			i++
		}
	}

	// If password not provided, prompt for it
	if entry.Password == "" {
		fmt.Print("Enter password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
		entry.Password = string(bytePassword)
	}

	// Save to database
	if err := database.SavePassword(entry); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password '%s' saved successfully!\n", entry.Name)
}

// handleGet handles retrieving a password
func handleGet() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s get <name>\n", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[2]
	entry, err := database.GetPassword(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	displayPasswordEntry(entry)
}

// handleList handles listing all passwords
func handleList() {
	entries, err := database.ListPasswords()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing passwords: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("No passwords found.")
		return
	}

	fmt.Printf("Found %d passwords:\n\n", len(entries))
	for _, entry := range entries {
		fmt.Printf("Name: %s\n", entry.Name)
		if entry.Username != "" {
			fmt.Printf("Username: %s\n", entry.Username)
		}
		if entry.URL != "" {
			fmt.Printf("URL: %s\n", entry.URL)
		}
		if len(entry.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(entry.Tags, ", "))
		}
		fmt.Printf("Updated: %s\n", entry.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("---")
	}
}

// handleDelete handles deleting a password
func handleDelete() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s delete <name>\n", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[2]
	
	// Confirm deletion
	fmt.Printf("Are you sure you want to delete password '%s'? (y/N): ", name)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	if err := database.DeletePassword(name); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password '%s' deleted successfully!\n", name)
}

// handleSearch handles searching passwords
func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s search <query>\n", os.Args[0])
		os.Exit(1)
	}

	query := os.Args[2]
	entries, err := database.SearchPasswords(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching passwords: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Printf("No passwords found matching '%s'.\n", query)
		return
	}

	fmt.Printf("Found %d passwords matching '%s':\n\n", len(entries), query)
	for _, entry := range entries {
		fmt.Printf("Name: %s\n", entry.Name)
		if entry.Username != "" {
			fmt.Printf("Username: %s\n", entry.Username)
		}
		if entry.URL != "" {
			fmt.Printf("URL: %s\n", entry.URL)
		}
		fmt.Println("---")
	}
}

// handleStats handles displaying database statistics
func handleStats() {
	stats, err := database.GetStats()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting stats: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database Statistics:")
	fmt.Printf("Total passwords: %d\n", stats["total_passwords"])
	fmt.Printf("Database size: %d bytes\n", stats["database_size"])
	fmt.Printf("Created: %s\n", stats["created_at"])
}

// handleAnalyze handles password strength analysis
func handleAnalyze() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s analyze <password>\n", os.Args[0])
		os.Exit(1)
	}

	password := os.Args[2]
	analysis := generator.AnalyzePasswordStrength(password)

	fmt.Println("Password Strength Analysis:")
	fmt.Printf("Length: %d characters\n", analysis["length"])
	fmt.Printf("Has uppercase: %t\n", analysis["has_uppercase"])
	fmt.Printf("Has lowercase: %t\n", analysis["has_lowercase"])
	fmt.Printf("Has numbers: %t\n", analysis["has_numbers"])
	fmt.Printf("Has symbols: %t\n", analysis["has_symbols"])
	fmt.Printf("Unique characters: %d\n", analysis["unique_chars"])
	fmt.Printf("Strength score: %d/7\n", analysis["strength_score"])
	fmt.Printf("Strength level: %s\n", analysis["strength_level"])
}

// displayPasswordEntry displays a password entry
func displayPasswordEntry(entry *storage.PasswordEntry) {
	fmt.Printf("Name: %s\n", entry.Name)
	if entry.Username != "" {
		fmt.Printf("Username: %s\n", entry.Username)
	}
	fmt.Printf("Password: %s\n", entry.Password)
	if entry.URL != "" {
		fmt.Printf("URL: %s\n", entry.URL)
	}
	if entry.Notes != "" {
		fmt.Printf("Notes: %s\n", entry.Notes)
	}
	if len(entry.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(entry.Tags, ", "))
	}
	fmt.Printf("Created: %s\n", entry.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", entry.UpdatedAt.Format("2006-01-02 15:04:05"))
}

// showHelp displays help information
func showHelp() {
	fmt.Printf("%s v%s\n\n", appName, version)
	fmt.Println("Usage:")
	fmt.Printf("  %s <command> [options]\n\n", os.Args[0])
	fmt.Println("Commands:")
	fmt.Println("  generate, gen     Generate a new password")
	fmt.Println("  save              Save a password")
	fmt.Println("  get, find         Retrieve a password")
	fmt.Println("  list              List all passwords")
	fmt.Println("  delete, del       Delete a password")
	fmt.Println("  search            Search passwords")
	fmt.Println("  stats             Show database statistics")
	fmt.Println("  analyze           Analyze password strength")
	fmt.Println("  help              Show this help message")
	fmt.Println("  version           Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Printf("  %s generate --length 20 --uppercase --numbers --symbols\n", os.Args[0])
	fmt.Printf("  %s save gmail --username user@example.com --password mypass\n", os.Args[0])
	fmt.Printf("  %s get gmail\n", os.Args[0])
	fmt.Printf("  %s analyze mypassword123\n", os.Args[0])
}

// showVersion displays version information
func showVersion() {
	fmt.Printf("%s v%s\n", appName, version)
}
