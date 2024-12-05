package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

	file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
			checkAndTruncateBySize(*logFile, *maxSize)
		} else if *maxLines > 0 {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Checking if truncation by lines is needed\n")
			}
			checkAndTruncateByLines(*logFile, *maxLines)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
	}
}

func checkAndTruncateBySize(filePath string, maxSize int64) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file for truncation: %v\n", err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stating file: %v\n", err)
		return
	}

	if info.Size() <= maxSize {
		return
	}

	// Seek to the point where we want to start reading
	offset := info.Size() - maxSize
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error seeking file for truncation: %v\n", err)
		return
	}

	// Read the remaining content
	remainingContent := make([]byte, maxSize)
	bytesRead, err := file.Read(remainingContent)
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading file for truncation: %v\n", err)
		return
	}

	// Truncate and rewrite the file
	err = file.Truncate(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error truncating file: %v\n", err)
		return
	}

	file.Seek(0, 0)
	_, err = file.Write(remainingContent[:bytesRead])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing truncated content: %v\n", err)
	}
}

func checkAndTruncateByLines(filePath string, maxLines int) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file for truncation: %v\n", err)
		return
	}
	defer file.Close()

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

	err = file.Truncate(0)
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
