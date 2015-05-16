# Alpine Linux FlashPlugin installer

## Description

It simply download `x86_64` flashplugin latest version from archlinux official repositories and install:

- `/usr/lib/mozilla/plugins/libflashplayer.so`

and create an empty `.so` to fake glibc dependency:

- `/usr/local/lib/ld-linux-x86-64.so.2`

that's all

## Dependencies

- [go](https://golang.org/)
- [xz](http://tukaani.org/xz/): to extract archives
- [gcc](https://gcc.gnu.org/): to create the fake glibc

all are available in Alpine main repo, so

```
$ apk add go xz gcc
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

## Contributor

- [dalias](http://www.musl-libc.org/)
