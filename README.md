# Ender

Ender is a safe and secure, temporary secret storage system for applications and systems that need a secure temporary storage solution for keys, and other security credentials.

## Developing

Build and test with `make`.

## Information

### Why?
Cross platform secret storage is difficult, and annoying, esecially on headless sytems. Ender sets out to solve this problem by creating a solution on-top of existing keychain systems that make securely fetching data easier, and more reliable across platforms.

Specifically, ender provides a cross platform secure secret store, independent of secret-server, dbus on Linux - enabling headless/terminal applications to securely store and persist secret data through a single binary.

### How it works?
The ender binary is cross platform, and provides two main commands: a daemon, and a CLI.

The daemon provides an zero-knowledge, in-memory key value store accessible over gPRC over a Unix socket. It is intended to store encrypted secrets provided by the CLI. On termination, the daemon wipes all secret.

The CLI provides an interface to connect to the daemon, sending it encrypted key/value pairs.

On MacOS, the initial key material is stored in Keychain and on Linux, key material is stored in Linux Kernel Keyring. All encrypted pieces of data are encrypted using XSalsa2020-Poly1305 via libsodium.

### Security Limitations
Ender is a tool to provide encrypted secret storage within the context of a secure user session. While the underlying key material stored in the user's keychain is secure, Ender is intended to provide secrets to ANY client that requests it within the user namespace. That means:

- Any processes running within the user space can access data stored in any ender chest.
- Ender depends and relies upon the security of the user namespace, and is vulernable to root users who can access the socket, or su into the user. Consider using session namespaces rather than user namespaces when starting `ender daemon`.

## Daemon
The ender daemon is the backend storage system. On termination it will delete all data stored in all chests, as well as the secret encryption key.

Ender can be launched by running `ender daemon`. It is recommended to use the provided `ender.service` systemd script and install it to your user's systemd folder and start it on login to ensure you have a keychain available across multiple terminal sessions.

```
mkdir -p ~/.config/systemd/user
cp ender.service ~/.config/systemd/user
systemctl --user start ender
systemctl --user enable ender
```

### Session
Ender provides an optional `daemon-helper`, which can be used to create _session_ keychains. If you need a per-session secret store you can add the following to your .profile or similar file:

```bash
eval `ender daemon-helper`
```

`daemon-helper` outputs the generated SOCKET FILE and the default chest to your environment variables.

Note that if you're using session daemons, make sure you add the appropriate kill to your bash .logout or similar file, as otherwise you'll just have multiple ender daemons running.

> Sessions are available across all open terminals by design, but will be terminated when the spawning terminal closes if using a bash .logout

## CLI
Ender ships with a dead simple CLI that lets you get, set, delete, and check if a key exists. The CLI will output the results to stdout, and will return an exit status of 0 if successful, and 1 otherwise. Individual applications may define specific chests to use via the `--chest`. View the associated `--help` documentation for each option prior to running.

```bash
$ ender cli set foo bar // 0
$ ender cli get foo // 0
   bar
$ ender cli exists foo // 0
$ ender cli exists bar // 1
$ ender cli delete foo // 0
```

The CLI is useful for testing and working with keys, however applications and scripting should work directly with the provided gRPC API.

## API
As a gRPC application Ender provides it's .proto file in the `protobuf` folder for use in other applications. Applications interfacing with the API over gRPC in golang can utilize the provided API

#### Set a key
```go
import client "github.com/kaidyth/ender/client"
if result, err := client.Set(socketAddressPath, chestName, key, value); err == nil && result.Ok {
    // Key was added
}
```

#### Get a key
```go
if result, err := client.Get(socketAddressPath, chestName, key); err == nil {
   fmt.Printf("%s", result) // result is the returned key
}
```

#### Delete a key
```go
if result, err := client.Del(socketAddressPath, chestName, key); err == nil && result.Ok {
   // Key was deleted
}
```

#### Check if a key exists
```go
if result, err := client.Exists(socketAddressPath, chestName, key); err == nil && result.Exists {
   // Key exists
}
```