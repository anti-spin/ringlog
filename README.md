# Ringlog

[![Go](https://github.com/anti-spin/ringlog/actions/workflows/go.yml/badge.svg)](https://github.com/anti-spin/ringlog/actions/workflows/go.yml)
[![Latest Version](https://img.shields.io/github/v/tag/anti-spin/ringlog?sort=semver)](https://github.com/anti-spin/ringlog/releases)
[![License](https://img.shields.io/github/license/anti-spin/ringlog)](https://github.com/anti-spin/ringlog/blob/main/LICENSE)


**Ringlog** is is a pipe-friendly utility to manage log files by capping their size or line count. It reads input from stdin and writes to a specified log file, ensuring that the log does not grow beyond the defined size or number of lines. When the log file exceeds the defined limit, it truncates the oldest entries to maintain a circular logging behavior.

## Usage
```sh
ringlog [options]
```

## Options
- `-s <max_size_bytes>`: Maximum size of the log file in bytes.
- `-l <max_lines>`: Maximum number of lines in the log file.
- `-f <log_file>`: Path to the log file where the input will be appended.
- `-v`: Enable verbose output. Prints information about log actions to stderr.

## Examples

### When to Use ringlog vs. Traditional Logging Solutions

**ringlog** is a lightweight solution for scenarios where you need quick and simple log capping, especially for small scripts or cronjobs. However, for more complex logging needs, it might be better to use established logging systems like `logrotate`, `syslog`, `journalctl`, or centralized logging services (e.g., ELK Stack, Splunk).

- Use **logrotate** when you need robust log rotation with advanced features like compression, scheduled rotations, and retention policies.
- Use **syslog** or **journalctl** for system-level logging, especially when logs from multiple services need to be aggregated and stored with consistent formatting.
- Use **centralized logging services** when you need to aggregate logs from multiple servers or containers, provide analytics, and support auditing.

**ringlog** fits well for quick, local use cases where you want to avoid complex setup or dependencies, keeping the solution lean and script-friendly.

### Typical Cronjob Scenarios
ringlog can also be helpful for a variety of other real-world DevOps scenarios where logs grow rapidly and need to be kept under control.

#### Docker Container Logs
When running Docker containers, logs are often redirected to files. If these logs are not managed, they can quickly consume disk space. Use ringlog to limit log file growth:
```sh
docker run --rm my_app 2>&1 | ringlog -s 100000000 -f /var/log/my_app/docker.log
```

#### Backup Script Logs
Automated backup scripts can generate logs that grow in size over time. Use ringlog to keep only the latest output:
```sh
0 2 * * * /usr/local/bin/backup.sh 2>&1 | ringlog -s 200000 -f /var/log/backup/backup.log
```

#### System Monitoring Logs
System monitoring scripts like custom health checks can generate frequent logs. Use ringlog to keep these logs concise:
```sh
*/5 * * * * /usr/local/bin/system_health_check.sh 2>&1 | ringlog -l 500 -f /var/log/monitoring/health_check.log
```

#### rsync Synchronization Logs
Logs generated by `rsync` during scheduled file synchronization can grow very large, especially when there are many changes. Use ringlog to cap these logs:
```sh
rsync -avz /source/ /destination/ 2>&1 | ringlog -l 1000 -f /var/log/rsync/sync.log
```

#### Matomo Analytics
Matomo's housekeeping cron job can generate verbose logs that grow indefinitely if not managed. Use ringlog to cap the log size:
```sh
*/15 * * * * /usr/bin/php /var/www/matomo/console core:archive 2>&1 | ringlog -s 50000000 -f /var/log/matomo/housekeeping.log
```

#### Laravel Artisan
Laravel's task scheduler (artisan) often runs every minute, and its logs can accumulate quickly. Use ringlog to cap the log lines:
```sh
* * * * * /usr/bin/php /var/www/html/artisan schedule:run >> /dev/null 2>&1 | ringlog -l 1000 -f /var/log/laravel/artisan_schedule.log
```

### General Usage
ringlog is primarily used to cap log output from scripts or cron jobs to ensure log files remain manageable:
```sh
./my_script.sh | ringlog -s 10000000 -f /var/log/my_log.log
```

### Capturing Both `stdout` and `stderr`
To capture both `stdout` and `stderr` from a script:
```sh
./my_script.sh 2>&1 | ringlog -l 100 -f /tmp/test.log
```

## License
MIT License.
