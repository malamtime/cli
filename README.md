# ShellTime CLI [![codecov](https://codecov.io/gh/malamtime/cli/graph/badge.svg?token=N09WIJHNI2)](https://codecov.io/gh/malamtime/cli)

The foundation CLI tool for shelltime.xyz - a platform for tracking DevOps work.

AnnatarHe: [![shelltime](https://api.shelltime.xyz/badge/AnnatarHe/count)](https://shelltime.xyz/users/AnnatarHe)

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/malamtime/installation/master/install.bash | bash
```

## Configuration

The CLI stores its configuration in `$HOME/.shelltime/config.toml`.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `token` | string | `""` | Your authentication token for shelltime.xyz |
| `apiEndpoint` | string | `"https://api.shelltime.xyz"` | The API endpoint URL for shelltime.xyz |
| `webEndpoint` | string | `"https://shelltime.xyz"` | The web interface URL for shelltime.xyz |
| `flushCount` | integer | `10` | Number of records to accumulate before syncing to server |
| `gcTime` | integer | `14` | Number of days to keep tracked data before garbage collection |
| `dataMasking` | boolean | `true` | Enable/disable masking of sensitive data in tracked commands |
| `enableMetrics` | boolean | `false` | Enable detailed command metrics tracking (WARNING: May impact performance) |
| `endpoints` | array | `[]` | Additional API endpoints for development or testing |

Example configuration:
```toml
token = "your-token-here"
apiEndpoint = "https://api.shelltime.xyz"
webEndpoint = "https://shelltime.xyz"
flushCount = 10
gcTime = 14
dataMasking = true
enableMetrics = false
```

⚠️ Note: Setting `enableMetrics` to `true` will track detailed metrics for every command execution. Only enable this when requested by developers for debugging purposes, as it may impact shell performance.

## Commands

### Authentication

```bash
shelltime auth [--token <your-token>]
```

Initializes the CLI with your shelltime.xyz authentication token. This command needs to be run before using other features.

Options:
- `--token`: Your personal access token from shelltime.xyz. if omit, you can also redirect to website to auth

Example:
```bash
shelltime auth --token abc123xyz
```

### Track

```bash
shelltime track [options]
```

Tracks your shells activities and sends them to shelltime.xyz.

Options:
- TODO: List track command options

Example:
```bash
shelltime track # TODO: Add example
```

### Sync

```bash
shelltime sync
```

Manually triggers synchronization of locally tracked commands to the shelltime.xyz server. This command can be useful when:
- You want to force an immediate sync without waiting for the automatic sync threshold
- You're troubleshooting data synchronization issues
- You need to ensure all local data is uploaded before system maintenance

Example:
```bash
shelltime sync
```

There are no additional options for this command as it simply processes and uploads any pending tracked commands according to your configuration settings.

### GC (Garbage Collection)

```bash
shelltime gc [options]
```

Performs cleanup of old tracking data and temporary files.

Options:
- TODO: List GC command options

Example:
```bash
shelltime gc # TODO: Add example
```

## Performance

### Local Storage Performance
The standard command saving process performs efficiently as it only involves file I/O operations. Testing on a MacBook Pro M1-Pro (14-inch) with 1TB storage shows consistent write latencies under 8ms, which should not impact your daily operations.

### Network Synchronization
Server synchronization times can vary significantly based on your geographical location:

- Users in Southeast Asia (near Singapore servers): ~100ms
- Users in other regions: May experience longer latency

If you experience slower synchronization times due to your location, we recommend:

1. Increasing the `FlushCount` value in `~/.shelltime/config.toml` to accumulate more commands before syncing
2. Manually running `shelltime sync` during off-peak hours

Example configuration for users in regions far from Singapore:
```toml
FlushCount = 100  # Increased from default 10
```

This configuration reduces the frequency of automatic syncs while ensuring your command history is still preserved locally.

## Version Information

Use `shelltime --version` or `shelltime -v` to display the current version of the CLI.


# Daemon(From 0.1.0)

a client daemon service that could process request from cli and sync to server

To install the service:

For Linux (systemd):
1. Copy the binary to `/usr/local/bin/shelltime-daemon`
2. Copy `shelltime.service` to `/etc/systemd/system/`
3. Run:
```bash
sudo systemctl daemon-reload
sudo systemctl enable shelltime
sudo systemctl start shelltime
```

For macOS:
1. Copy the binary to `/usr/local/bin/shelltime-daemon`
2. Copy `xyz.shelltime.daemon.plist` to `/Library/LaunchDaemons/`
3. Run:
```bash
sudo launchctl load /Library/LaunchDaemons/xyz.shelltime.daemon.plist
```

This implementation provides:
- A daemon service that listens on a Unix domain socket
- Handling of status and track messages
- Proper shutdown handling
- Service description files for both Linux and macOS
- Basic logging
- JSON message format for communication

You can extend the `handleStatus` and `handleTrack` functions to implement the specific functionality you need for each command.


## Support

For support, please contact: annatar.he+shelltime.xyz@gmail.com

## License

Copyright (c) 2024 shelltime.xyz Team
