# GoSDK Installer

A command-line tool to install and show latest Go SDK version.

## Install

```sh
go install github.com/nopcoder/gosdk@latest
```

## Usage

This command-line tool supports two commands:

### remote

Fetches the latest Go version information and prints it to the console. Example:

```sh
gosdk remote
```

### install

Installs a specific or the latest Go SDK version. You can specify the version as an argument, like this:

```sh
# latest stable
gosdk install

# or specific version
gosdk install go1.23.3
```
