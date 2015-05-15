# Alpine Linux FlashPlugin installer

## Description

It simply download flashplugin and glibc latests versions from archlinux official repositories and install:

- `/usr/lib/mozilla/plugins/libflashplayer.so` from flashplugin
- `/usr/local/lib/ld-<version>.so` from glibc

that's all

## Dependecies

- [go](https://golang.org/)
- [xz](http://tukaani.org/xz/)

both are available in Alpine main repo, so

```
$ apk add go xz
```

and you're done

## Install

```
$ go get github.com/WnP/InstallFlash
```

## Usage

```
$ InstallFlash
```

use `sudo` if you don't have write permission for `/usr/`

## License

See [LICENSE file](https://github.com/WnP/InstallFlash/blob/master/LICENSE)
