# hosts-timer

Block websites on your system with `/etc/hosts`, and re-enable them on demand for a specified amount of time.

## Installation

### macOS via Homebrew

```shell
brew install cdzombak/oss/hosts-timer
```

### Debian via apt repository

[Install my Debian repository](https://www.dzombak.com/blog/2025/06/updated-instructions-for-installing-my-debian-package-repositories/) if you haven't already:

```shell
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://dist.cdzombak.net/keys/dist-cdzombak-net.gpg -o /etc/apt/keyrings/dist-cdzombak-net.gpg
sudo chmod 644 /etc/apt/keyrings/dist-cdzombak-net.gpg
sudo mkdir -p /etc/apt/sources.list.d
sudo curl -fsSL https://dist.cdzombak.net/cdzombak-oss.sources -o /etc/apt/sources.list.d/cdzombak-oss.sources
sudo chmod 644 /etc/apt/sources.list.d/cdzombak-oss.sources
sudo apt update
```

Then install `hosts-timer` via `apt-get`:

```shell
sudo apt-get install hosts-timer
```

### Manual installation from build artifacts

Pre-built binaries for Linux and macOS on various architectures are downloadable from each [GitHub Release](https://github.com/cdzombak/hosts-timer/releases). Debian packages for each release are available as well.

### From source

Requires Go to be installed:

```shell script
git clone https://github.com/cdzombak/hosts-timer.git
cd hosts-timer
make install
```

## Usage

### Setup for domain(s)

```shell script
sudo hosts-timer -install facebook.com twitter.com
```

(`-disable` is an alias for `-install`, which _disables_ the given sites.)

### Use a site for a specified duration of time

```shell script
sudo hosts-timer -time 5m twitter.com
```

Durations must be acceptable by [`time.ParseDuration`](https://golang.org/pkg/time/#ParseDuration).

### Remove hosts-timer block for domain(s)

```shell script
sudo hosts-timer -uninstall twitter.com facebook.com
```

(`-enable` is an alias for `-uninstall`, which _enables_ the given sites.)

## Setup Hints

### `.zshrc` helpers

```shell script
function use-facebook() {
    if [ $# -eq 0 ]; then
        echo "Usage: use-facebook DURATION"
        echo "   eg. use-facebook 5m"
        return 1
    fi
    sudo hosts-timer -time "$1" facebook.com
}

function use-twitter() {
    if [ $# -eq 0 ]; then
        echo "Usage: use-twitter DURATION"
        echo "   eg. use-twitter 5m"
        return 1
    fi
    sudo hosts-timer -time "$1" twitter.com
}
```

### Alfred helpers

Install the [included Alfred helpers](https://github.com/cdzombak/hosts-timer/tree/main/alfred), and then you can use Alfred to enable Twitter or Facebook for a given amount of time (default 5 minutes).

![Screenshot of Alfred helper for Twitter](https://github.com/cdzombak/hosts-timer/raw/main/alfred/screenshot.png)

### Create `/etc/sudoers.d/hosts-timer`

```
cdzombak ALL=NOPASSWD: /usr/local/bin/hosts-timer
```

And then:
```
sudo chown root:wheel /etc/sudoers.d/hosts-timer
sudo chmod 440 /etc/sudoers.d/hosts-timer
```

## Author

Chris Dzombak
- [dzombak.com](https://www.dzombak.com)
- [github/cdzombak](https://www.github.com/cdzombak)
