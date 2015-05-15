# WARNING: Alpha stage, not functionnal

# Alpine Linux FlashPlugin installer

## Description

It simply download flashplugin and glibc latests versions from archlinux official repositories and install:

- `/usr/lib/mozilla/plugins/libflashplayer.so` from flashplugin
- `/usr/local/lib/ld-<version>.so` from glibc

that's all

## Install

```
$ go get github.com/WnP/InstallFlash
```

## Usage

```
$ InstallFlash
```

use `sudo` if you don't have write permission for `/usr/`
