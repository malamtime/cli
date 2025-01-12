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

### Command Execution Performance

By default, the CLI performs synchronization directly which may impact shell responsiveness in certain scenarios:

- Standard command saving: <8ms (local file I/O only)
- Network synchronization:
  - Southeast Asia (Singapore servers): ~100ms
  - Other regions: Can vary significantly based on location

### Recommended: Daemon Mode

If you experience latency issues, we strongly recommend using daemon mode for better performance:

```bash
sudo ~/.shelltime/bin/shelltime daemon install
```

Benefits of daemon mode:
- Asynchronous command tracking (shell blocking time <8ms)
- Background synchronization handling
- No impact on shell responsiveness
- Reliable data delivery even during network issues

The daemon service:
1. Runs in the background as a system service
2. Handles all network synchronization operations
3. Buffers commands during connectivity issues
4. Automatically retries failed synchronizations

For users experiencing high latency, daemon mode is the recommended configuration. You can also adjust `FlushCount` in the config for additional optimization:

```toml
FlushCount = 100  # Increased buffer size for less frequent syncs
```

Note: Even without the daemon, all commands are still preserved locally first, ensuring no data loss during network issues.

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
