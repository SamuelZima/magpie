package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const VERSION = "v0.6.0"

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
		fmt.Fprintf(out, "Usage: magpie [-h|--help] [-v|--version] <string>|nest|nests|renest|abandon|pluck\n\n")
		fmt.Fprintf(out, "Options:\n")
		fmt.Fprintf(out, "\t-h, --help\tShow this help message\n")
		fmt.Fprintf(out, "\t-v, --version\tShow version\n")
		fmt.Fprintf(out, "\t<string>\tThought, note or idea which you want to save\n")
		fmt.Fprintf(out, "\tnest\t\tPrints the daily nest containing all your thoughts, notes or ideas\n")
		fmt.Fprintf(out, "\tnests\t\tPrints all notes from all existing nest files\n")
		fmt.Fprintf(out, "\trenest\t\tDeletes current nest file and all its contents\n")
		fmt.Fprintf(out, "\tabandon\t\tDeletes all nest files and all their contents\n")
		fmt.Fprintf(out, "\tpluck\t\tRemoves the last note from the daily nest\n")
	}

	flag.Parse()
	args.posArgs = flag.Args()

	return args
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	// Prompt user for input
	fmt.Printf("%s [y/N]: ", prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	// Evaluate input -> only explicit yes evaluates to true
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

func createOrRetrieveNest() (string, error) {
	// Create magpie dir, where the nest files will be stored
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ERROR: Unable to get user's home directory: %w", err)
	}
	magpieDir := home + "/.magpie/"
	if err := os.MkdirAll(magpieDir, 0700); err != nil {
		return "", fmt.Errorf("ERROR: Unable to create magpie directory: %w", err)
	}

	// Create today's nest file in the magpie dir if it doesn't exist already
	nestFilePath := magpieDir + time.Now().Format(time.DateOnly) + "-nest.txt"
	nestFile, err := os.OpenFile(nestFilePath, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return "", fmt.Errorf("ERROR: Unable to create/open nest file: %w", err)
	}
	nestFile.Close()

	return nestFilePath, nil
}

func nestCommand(nestFilePath string) error {
	// Read the nest file and print its contents, if it's not empty
	content, err := os.ReadFile(nestFilePath)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to read nest file: %w", err)
	}

	if len(content) == 0 {
		fmt.Println("No shiny notes in the nest yet. Add a note with: magpie <note>")
		return nil
	}

	fmt.Print(string(content))
	return nil
}

func nestsCommand(nestFilePath string) error {
	// Get all nest files
	paths, err := filepath.Glob(filepath.Join(filepath.Dir(nestFilePath), "*-nest.txt"))
	if err != nil {
		return fmt.Errorf("ERROR: Unable to list nest files: %w", err)
	}

	found := false

	// For every nest file print its date and its contents
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("ERROR: Unable to read nest file: %w", err)
		}
		if len(content) == 0 {
			continue
		}
		found = true
		fmt.Printf("======== %s ========\n%s\n", strings.TrimSuffix(filepath.Base(path), "-nest.txt"), content)
	}

	// If no nests are found, print message
	if !found {
		fmt.Println("No shiny notes in any nest yet. Add a note with: magpie <note>")
	}
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

func abandonCommand(nestFilePath string) error {
	if !confirm("This action will delete all nest files. Proceed?") {
		fmt.Println("Aborted. All nests kept.")
		return nil
	}

	// Get all nest files
	paths, err := filepath.Glob(filepath.Join(filepath.Dir(nestFilePath), "*-nest.txt"))
	if err != nil {
		return fmt.Errorf("ERROR: Unable to list nest files: %w", err)
	}

	// Delete every nest file
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("ERROR: Unable to delete nest file: %w", err)
		}
	}

	fmt.Println("All nests abandoned.")
	return nil
}

func pluckCommand(nestFilePath string) error {
	// Read nest file
	content, err := os.ReadFile(nestFilePath)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to read nest file: %w", err)
	}

	if len(content) == 0 {
		fmt.Println("No shiny notes in the nest to pluck. Add a note with: magpie <note>")
		return nil
	}

	// Keep everything up to and including the newline before the last note, i.e. remove only last line/note from file
	trimmed := strings.TrimSuffix(string(content), "\n")
	kept := trimmed[:strings.LastIndex(trimmed, "\n")+1]
	if err := os.WriteFile(nestFilePath, []byte(kept), 0600); err != nil {
		return fmt.Errorf("ERROR: Unable to write nest file: %w", err)
	}
	return nil
}

func addNoteToNestCommand(nestFilePath string, note string) error {
	// Open the nest file
	nestFile, err := os.OpenFile(nestFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
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
	nestFilePath, err := createOrRetrieveNest()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Unquoted words in positional arguments form a single note
	noteOrCommand := strings.Join(args.posArgs, " ")

	switch noteOrCommand {
	case "nest":
		err = nestCommand(nestFilePath)
	case "nests":
		err = nestsCommand(nestFilePath)
	case "renest":
		err = renestCommand(nestFilePath)
	case "abandon":
		err = abandonCommand(nestFilePath)
	case "pluck":
		err = pluckCommand(nestFilePath)
	default:
		err = addNoteToNestCommand(nestFilePath, noteOrCommand)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
