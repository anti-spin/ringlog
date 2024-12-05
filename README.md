# Ringlog

**Ringlog** is a simple command-line utility that captures output from standard input and writes it to a log file, automatically capping the file size or number of lines. When the log file exceeds the defined limit, it truncates the oldest entries to maintain a circular logging behavior.

## Example Usage

```
my_cron_script.sh | ringlog -l 1000 -f /path/to/logfile.log
```
