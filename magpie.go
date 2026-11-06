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

const VERSION = "v0.8.0"

const (
	reset  = "\x1b[0m"
	dim    = "\x1b[2m"
	red    = "\x1b[91m"
	yellow = "\x1b[93m"
	cyan   = "\x1b[96m"
)

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
		fmt.Fprintf(out, "magpie - snatches your shiny thoughts into a daily nest.\n\n")
		fmt.Fprintf(out, "Usage: magpie [-h|--help] [-v|--version] <note>|pluck|nest|renest|nests|abandon\n\n")
		fmt.Fprintf(out, "Commands:\n")
		fmt.Fprintf(out, "\t<note>\t\tSave a note, thought or idea to today's nest file\n")
		fmt.Fprintf(out, "\tpluck\t\tDelete the last note from today's nest\n")
		fmt.Fprintf(out, "\tnest\t\tPrint notes from today's nest\n")
		fmt.Fprintf(out, "\trenest\t\tDelete today's nest file and all its contents\n")
		fmt.Fprintf(out, "\tnests\t\tPrint notes from all existing nest files\n")
		fmt.Fprintf(out, "\tabandon\t\tDelete all nest files and all their contents\n\n")
		fmt.Fprintf(out, "Options:\n")
		fmt.Fprintf(out, "\t-h, --help\tPrint this help message\n")
		fmt.Fprintf(out, "\t-v, --version\tPrint version\n")
	}

	flag.Parse()
	args.posArgs = flag.Args()

	return args
}

func paint(colour string, text string) string {
	// Build string with ANSI colour codes
	return colour + text + reset
}

func confirm(prompt string) bool {
	// Setup reader from STDIN
	reader := bufio.NewReader(os.Stdin)

	// Prompt user for input
	fmt.Printf("%s %s ", paint(yellow, prompt), paint(yellow, "[y/N]:"))
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	// Evaluate input -> only explicit yes evaluates to true
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "y", "yes", "Yes":
		return true
	default:
		return false
	}
}

func createOrRetrieveNest() (string, error) {
	// Create or retrieve magpie dir under HOME dir, all nest files are stored in the magpie dir
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ERROR: Unable to get user's home directory: %w", err)
	}
	magpieDir := home + "/.magpie/"
	if err := os.MkdirAll(magpieDir, 0700); err != nil {
		return "", fmt.Errorf("ERROR: Unable to create magpie directory: %w", err)
	}

	// Create or retrieve today's nest file from the magpie dir
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
		fmt.Println(paint(dim, "No shiny notes in the nest yet."))
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
		fmt.Printf("======== %s ========\n%s\n", paint(cyan, strings.TrimSuffix(filepath.Base(path), "-nest.txt")), content)
	}

	// If no nests are found, print message
	if !found {
		fmt.Println(paint(dim, "No shiny notes in any nest yet."))
	}
	return nil
}

func renestCommand(nestFilePath string) error {
	// Ask for confirmation
	if !confirm("This action will delete current nest file. Proceed?") {
		fmt.Println(paint(dim, "Aborted. Nest kept."))
		return nil
	}

	// Delete the nest file
	err := os.Remove(nestFilePath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("ERROR: Unable to delete nest file: %w", err)
	}

	fmt.Println(paint(dim, "Renested."))
	return nil
}

func abandonCommand(nestFilePath string) error {
	// Ask for confirmation
	if !confirm("This action will delete all nest files. Proceed?") {
		fmt.Println(paint(dim, "Aborted. All nests kept."))
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

	fmt.Println(paint(dim, "All nests abandoned."))
	return nil
}

func pluckCommand(nestFilePath string) error {
	// Read nest file
	content, err := os.ReadFile(nestFilePath)
	if err != nil {
		return fmt.Errorf("ERROR: Unable to read nest file: %w", err)
	}

	if len(content) == 0 {
		fmt.Println(paint(dim, "No shiny notes in the nest to pluck."))
		return nil
	}

	// Keep everything up to and including the newline before the last note, i.e. remove only last line/note from file
	trimmed := strings.TrimSuffix(string(content), "\n")
	lastNewline := strings.LastIndex(trimmed, "\n")
	plucked := trimmed[lastNewline+1:]
	kept := trimmed[:lastNewline+1]
	if err := os.WriteFile(nestFilePath, []byte(kept), 0600); err != nil {
		return fmt.Errorf("ERROR: Unable to write nest file: %w", err)
	}

	fmt.Printf("%s\n%s\n", paint(yellow, "Plucked:"), plucked)
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
	savedNote := timestamp + "  " + note + "\n"
	if _, err := nestFile.WriteString(savedNote); err != nil {
		return fmt.Errorf("ERROR: Unable to append the nest file: %w", err)
	}

	fmt.Printf("%s\n%s\n", paint(cyan, "Snatched:"), savedNote)
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
		fmt.Fprintln(os.Stderr, paint(red, err.Error()))
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
		fmt.Fprintln(os.Stderr, paint(red, err.Error()))
		os.Exit(1)
	}
}
