# hosts-timer

Block websites on your system with `/etc/hosts`, and re-enable them on demand for a specified amount of time.

## Installation

Requires Go to be installed (tested with Go 1.15.x).

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
