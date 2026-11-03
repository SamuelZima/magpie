package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

const VERSION = "v0.1.1"

type arguments struct {
	help    bool
	version bool
	posArgs []string
}

func parseArgs() arguments {
	var args arguments
	flag.BoolVar(&args.help, "help", false, "show this help message")
	flag.BoolVar(&args.help, "h", false, "shorthand for --help")
	flag.BoolVar(&args.version, "version", false, "show version of app")
	flag.BoolVar(&args.version, "v", false, "shorthand for --version")

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "magpie - collects your shiny thoughts, ideas or notes into a daily nest file.\n\n")
		fmt.Fprintf(out, "Usage: magpie [-h|--help] [-v|--version] <string>|nest|renest\n\n")
		fmt.Fprintf(out, "Options:\n")
		fmt.Fprintf(out, "\t-h, --help\tShow this help message\n")
		fmt.Fprintf(out, "\t-v, --version\tShow version\n")
		fmt.Fprintf(out, "\t<string>\tThought, note or idea which you want to save\n")
		fmt.Fprintf(out, "\tnest\t\tPrints the daily nest containing all your thoughts, notes or ideas\n")
		fmt.Fprintf(out, "\trenest\t\tDeletes current nest file and all its contents\n")
	}

	flag.Parse()
	args.posArgs = flag.Args()

	return args
}

func createNest() (string, error) {
	// Create magpie dir, where the nest files will be stored
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ERROR: Unable to get user's home directory: %w", err)
	}
	magpieDir := home + "/.magpie/"
	if err := os.MkdirAll(magpieDir, 0755); err != nil {
		return "", fmt.Errorf("ERROR: Unable to create magpie directory: %w", err)
	}

	// Create today's nest file in the magpie dir if it doesn't exist already
	nestFilePath := magpieDir + time.Now().Format(time.DateOnly) + "-nest.md"
	nestFile, err := os.OpenFile(nestFilePath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("ERROR: Unable to create/open nest file: %w", err)
	}
	nestFile.Close()

	return nestFilePath, nil
}

func nestCommand(nestFilePath string) error {
	// Print the contents of the nest file
	content, err := os.ReadFile(nestFilePath)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to read nest file: %w", err)
	}
	fmt.Print(string(content))
	return nil
}

func renestCommand(nestFilePath string) error {
	// Delete the nest file
	err := os.Remove(nestFilePath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("ERROR: Unable to delete nest file: %w", err)
	}
	return nil
}

func addNoteToNestCommand(nestFilePath string, note string) error {
	// Open the nest file
	nestFile, err := os.OpenFile(nestFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to open the nest file: %w", err)
	}
	defer nestFile.Close()

	// Append timestamped note to the nest file
	timestamp := time.Now().Format(time.TimeOnly)
	if _, err := nestFile.WriteString(timestamp + "  " + note + "\n"); err != nil {
		return fmt.Errorf("ERROR: Unable to append the nest file: %w", err)
	}

	return nil
}

func main() {
	args := parseArgs()

	if args.help {
		flag.Usage()
		os.Exit(0)
	}

	if args.version {
		fmt.Println("magpie", VERSION)
		os.Exit(0)
	}

	if len(args.posArgs) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Create or retrieve the path to the nest file
	nestFilePath, err := createNest()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Unquoted words in positional arguments form a single note
	noteOrCommand := strings.Join(args.posArgs, " ")

	switch noteOrCommand {
	case "nest":
		err = nestCommand(nestFilePath)
	case "renest":
		err = renestCommand(nestFilePath)
	default:
		err = addNoteToNestCommand(nestFilePath, noteOrCommand)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
