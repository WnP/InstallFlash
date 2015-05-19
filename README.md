# Alpine Linux FlashPlugin installer

## Description

It simply download `x86_64` flashplugin latest version from archlinux official repositories and install:

- `/usr/lib/mozilla/plugins/libflashplayer.so`

and create an empty `.so` to fake glibc dependency:

- `/usr/local/lib/ld-linux-x86-64.so.2`

that's all

## Dependencies

If you plan to install it from sources:

- [go](https://golang.org/)
- [xz-dev](http://tukaani.org/xz/): to extract archives
- [gcc](https://gcc.gnu.org/): to create the fake glibc
- [wgo](https://github.com/skelterjohn/wgo): to manage dependencies versions

all, except `wgo`, are available in Alpine main repo, so

```
$ apk add go xz-dev gcc
```

to install `wgo`:

```
$ go get github.com/skelterjohn/wgo
```

## Install

### From source

```
$ git clone https://github.com/WnP/InstallFlash
$ cd InstallFlash
$ wgo restore
$ wgo install InstallFlash
```

`InstallFlash` binary is in `./bin` folder

### Using Alpine Package Manager

`InstallFlash` is actualy available from testing repository, so:

```
$ apk add installflash
```

will install it with no dependencies, thanks to [Carlo Landmeter](https://github.com/clandmeter)

## Usage

```
$ InstallFlash
```

use `sudo` if you don't have write permission for `/usr/`

## License

See [LICENSE file](https://github.com/WnP/InstallFlash/blob/master/LICENSE)

## Development

Pull requests are welcome

- fork
- make some changes
- commit you changes
- create a pull request

### TODO list

- Implement tests suite

## Contributors

- [dalias](http://www.musl-libc.org/)
- [Carlo Landmeter](https://github.com/clandmeter)
