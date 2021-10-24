# bagccgop

Build for all `gccgo`-supported platforms by default, disable those which you don't want (bagop with CGo support).

[![hydrun CI](https://github.com/pojntfx/bagccgop/actions/workflows/hydrun.yaml/badge.svg)](https://github.com/pojntfx/bagccgop/actions/workflows/hydrun.yaml)
[![Matrix](https://img.shields.io/matrix/bagccgop:matrix.org)](https://matrix.to/#/#bagccgop:matrix.org?via=matrix.org)
[![Binary Downloads](https://img.shields.io/github/downloads/pojntfx/bagccgop/total?label=binary%20downloads)](https://github.com/pojntfx/bagccgop/releases)

## Overview

bagccgop is a simple build tool for Go which tries to build your app for all platforms supported by `gccgo` by default. Instead of manually adding specific `GOOS`es and `GOARCH`es, bagccgop builds for all valid targets by default, and gives you the choice to disable those which you don't want to support or which can't be supported. It is a variant of [bagop](https://github.com/pojntfx/bagop) which uses `gccgo` instead of `gc`, the default Go compiler, and as such it does not intend to replace bagop, but instead tries to give a similar user experience for Go apps which rely on CGo or target platforms which are only supported by `gccgo` (such as 32-bit PowerPC).

## Installation

Static binaries are available on [GitHub releases](https://github.com/pojntfx/bagccgop/releases).

You can install them like so:

```shell
$ curl -L -o /tmp/bagccgop "https://github.com/pojntfx/bagccgop/releases/latest/download/bagccgop.linux-$(uname -m)"
$ sudo install /tmp/bagccgop /usr/local/bin
```

## Usage

🚧 This project is a work-in-progress! Instructions will be added as soon as it is usable. 🚧

🚀 **That's it!** We've successfully added support for most Debian ports to this app.

If you're enjoying bagccgop, the following projects might also be of help to you too:

- Also want to test these cross-compiled binaries? Check out [hydrun](https://github.com/pojntfx/hydrun)!
- Need to cross-compile without CGo? Check out [bagop](https://github.com/pojntfx/bagop)!
- Want to build fully-featured desktop GUI for all these platforms without CGo? Check out [Lorca](https://github.com/zserge/lorca)!
- Want to use SQLite without CGo? Check out [cznic/sqlite](https://gitlab.com/cznic/sqlite)!

## Reference

```shell
$ bagccgop --help
Build for all gccgo-supported platforms by default, disable those which you don't want (bagop with CGo support).
	Example usage: bagccgop -b mybin -x '(linux/alpha|linux/ppc64el)' -j "$(nproc)" 'main.go'
	Example usage (with plain flag): bagccgop -b mybin -x '(linux/alpha|linux/ppc64el)' -j "$(nproc)" -p 'go build -o $DST main.go'
	See https://github.com/pojntfx/bagccgop for more information.
	Usage: bagccgop [OPTION...] '<INPUT>'
	  -b, --bin string               Prefix of resulting binary (default "mybin")
  -d, --dist string              Directory build into (default "out")
  -x, --exclude string           Regex of platforms not to build for, i.e. (linux/alpha|linux/ppc64el)
  -e, --extra-args string        Extra arguments to pass to the Go compiler
  -g, --goisms                   Use Go's conventions (i.e. amd64) instead of uname's conventions (i.e. x86_64)
  -s, --hostPackages strings     Comma-seperated list of Debian packages to install for the host architecture
  -j, --jobs int                 Maximum amount of parallel jobs (default 1)
  -m, --manualPackages strings   Comma-seperated list of Debian packages to manually install for the selected architectures (i.e. those which would break the dependency graph)
  -a, --packages strings         Comma-seperated list of Debian packages to install for the selected architectures
  -p, --plain                    Sets GOARCH, GOARCH, CC, GCCGO, GOFLAGS and DST and leaves the rest up to you (see example usage)
  -r, --prepare string           Command to run before running the main command; will have only CC and GCCGO set (i.e. for code generation)
```

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

To build bagccgop locally, run:

```shell
$ git clone https://github.com/pojntfx/bagccgop.git
$ cd bagccgop
$ go run main.go --help
```

To build the convenience images with pre-built Debian `chroot`s, run:

```shell
$ docker buildx build --allow security.insecure -t ghcr.io/pojntfx/bagccgop-base-sid --load -f Dockerfile.sid .
$ docker buildx build --allow security.insecure -t ghcr.io/pojntfx/bagccgop-base-bullseye --load -f Dockerfile.bullseye .
```

Have any questions or need help? Chat with us [on Matrix](https://matrix.to/#/#bagccgop:matrix.org?via=matrix.org)!

## License

bagccgop (c) 2021 Felix Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
