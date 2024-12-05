package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	maxSize := flag.Int64("s", 0, "Maximum size of the log file in bytes (use either -s or -l)")
	maxLines := flag.Int("l", 0, "Maximum number of lines in the log file (use either -s or -l)")
	logFile := flag.String("f", "", "Log file path")
	verbose := flag.Bool("v", false, "Enable verbose output")

	flag.Parse()

	if *verbose {
		fmt.Fprintf(os.Stderr, "Verbose mode enabled\n")
		fmt.Fprintf(os.Stderr, "Log file: %s\n", *logFile)
		if *maxSize > 0 {
			fmt.Fprintf(os.Stderr, "Max size: %d bytes\n", *maxSize)
		}
		if *maxLines > 0 {
			fmt.Fprintf(os.Stderr, "Max lines: %d\n", *maxLines)
		}
	}

	if *logFile == "" || (*maxSize == 0 && *maxLines == 0) {
		fmt.Println("Usage: ringlog -s <max_size_bytes> or -l <max_lines> -f <log_file>\n\nPipe-friendly utility to manage log files by capping size or line count.")
		os.Exit(1)
	}

	file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if *verbose {
			fmt.Fprintf(os.Stderr, "Writing line to log: %s\n", line)
		}
		if _, err := file.WriteString(line + "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
			os.Exit(1)
		}

		// Check size or lines and truncate if needed
		if *maxSize > 0 {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Checking if truncation by size is needed\n")
			}
			checkAndTruncateBySize(file, *maxSize)
		} else if *maxLines > 0 {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Checking if truncation by lines is needed\n")
			}
			checkAndTruncateByLines(file, *maxLines)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
	}
}

func checkAndTruncateBySize(file *os.File, maxSize int64) {
	info, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stating file: %v\n", err)
		return
	}

	if info.Size() <= maxSize {
		return
	}

	// Truncate the file in place to keep only the last maxSize bytes
	remainingContent := make([]byte, maxSize)
	_, err = file.ReadAt(remainingContent, info.Size()-maxSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file for truncation: %v\n", err)
		return
	}

	err = file.Truncate(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error truncating file: %v\n", err)
		return
	}

	_, err = file.WriteAt(remainingContent, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing truncated content: %v\n", err)
	}
}

func checkAndTruncateByLines(file *os.File, maxLines int) {
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) <= maxLines {
		return
	}

	lines = lines[len(lines)-maxLines:]

	err := file.Truncate(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error truncating file: %v\n", err)
		return
	}

	file.Seek(0, 0)
	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
			return
		}
	}
}
