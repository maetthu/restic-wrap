# restic-wrap

[![release](https://github.com/maetthu/restic-wrap/actions/workflows/release.yml/badge.svg)](https://github.com/maetthu/restic-wrap/actions/workflows/release.yml)

Just a thin [restic](https://restic.net/) wrapper to ease backup/restore when using multiple backends.

## Usage

``` 
Restic wrapper tool with profile support

Usage:
  restic-wrap [command]

Available Commands:
  backup      Executes configured backup stages for all backends
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run adhoc restic command with all the necessary environment variables set for a specific backend

Flags:
  -b, --backend string   Backend to use (depending on the command, either the first one or all are used by default)
  -h, --help             help for restic-wrap
  -p, --profile string   Path to profile.yaml
  -v, --version          version for restic-wrap

Use "restic-wrap [command] --help" for more information about a command.
```

## Profile

```yaml
# Environment variables used for all restic calls
env:
  - name: RESTIC_CACHE_DIR
    value: /var/cache/restic
# List of restic backends (name, repository and password are mandatory)
backends:
  - name: rest
    repository: rest:https://restic.example.org/repo
    password: this-is-a-password
    # List of environment variables used for this backend only
    env:
      - name: RESTIC_TLS_CLIENT_CERT
        value: /etc/restic/cert.pem
  - name: aws
    repository: s3:s3.amazonaws.com/bucket-name
    password: this-is-also-a-password
    env:
      - name: AWS_ACCESS_KEY_ID
        value: ...
      - name: AWS_SECRET_ACCESS_KEY
        value: ...
      - name: AWS_DEFAULT_REGION
        value: ...
# These are the individual restic commands run in order for each backend when the "backup" command is invoked
stages:
  - command: backup
    args:
      - --exclude-if-present=.nobackup
      - --exclude-file=/etc/restic/excludes.txt
      - --exclude-caches
      - -o
      - s3.storage-class=STANDARD_IA
      - -x
      - /home
  - command: check
    args: []
  - command: forget
    args:
      - --keep-last
      - "7"
      - --keep-daily
      - "30"
      - --keep-weekly
      - "8"
      - --keep-monthly
      - "24"
      - --keep-yearly
      - "5"
      - --prune
# Commands to be run for notification. For each backend, it is called with
# /path/to/command <backend-name> <stage-name> <success|error> <message-if-available>
# (including all environment variables for the current backend)
notify:
  - /etc/restic/notify.sh
```