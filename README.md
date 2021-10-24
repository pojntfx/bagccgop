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

Let's assume we have a Go app called `hello-world` and we want to build it for as many platforms as possible using bagccgop. This is the `main.go`:

```go
package main

// #include <stdio.h>
//
// void print_using_c(char* s) {
//   printf("%s\n", s);
// }
import "C"

func main() {
	C.print_using_c(C.CString("Hello, world!"))
}
```

As you can see, it includes C code. It can be compiled like so:

```shell
$ go build -o out/hello-world main.go
$ file out/*
out/hello-world: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, BuildID[sha1]=bd2259005fcb565f71b26181762eeecaefe0bc31, for GNU/Linux 3.2.0, not stripped
```

But as we try to cross-compile, we start to get errors:

```shell
$ GOARCH=riscv64 go build -o out/hello-world main.go
package command-line-arguments: build constraints exclude all Go files in /home/pojntfx/Projects/hello-world
$ CGO_ENABLED=1 GOARCH=riscv64 go build -o out/hello-world main.go
# runtime/cgo
gcc_riscv64.S: Assembler messages:
gcc_riscv64.S:15: Error: no such instruction: `sd x1,-200(sp)'
gcc_riscv64.S:16: Error: no such instruction: `addi sp,sp,-200'
gcc_riscv64.S:17: Error: no such instruction: `sd x8,8(sp)'
gcc_riscv64.S:18: Error: no such instruction: `sd x9,16(sp)'
# ...
```

This is because in order to use C code in Go, we have to use CGo, which requires a C compiler for the specified `GOOS` and `GOARCH`. `bagccgo` simplifies this process. To use it, we first create an interactive shell in of the base images using [hydrun](https://github.com/pojntfx/hydrun):

```shell
$ hydrun -o ghcr.io/pojntfx/bagccgop-base-sid -e '--privileged' -i bash
2021/10/24 19:03:37 /usr/bin/docker inspect ghcr.io/pojntfx/bagccgop-base-sid-amd64
2021/10/24 19:03:37 /usr/bin/docker run -it -v /home/pojntfx/Projects/hello-world:/data:z --platform linux/amd64 --privileged ghcr.io/pojntfx/bagccgop-base-bullseye /bin/sh -c cd /data && bash
root@bed49fc87e96:/data#
```

We now have a temporary shell which we can use to cross-compile. There are currently two convenience images available; `ghcr.io/pojntfx/bagccgop-base-bullseye`, which is based on Debian 11 (Bullseye), and `ghcr.io/pojntfx/bagccgop-base-sid`, which is based on Debian Unstable. It is recommended to always use the former; the latter might have broken multi-architecture packages, but is more up-to-date. Note that for architectures which are only available as Debian ports (`linux/alpha`, `linux/ppc`, `linux/ppc64`, `linux/sparc64` and `linux/riscv64`), Debian Unstable is the only option.

In the next step, let's install bagccgop inside the shell:

```shell
$ curl -L -o /tmp/bagccgop "https://github.com/pojntfx/bagccgop/releases/latest/download/bagccgop.linux-$(uname -m)"
$ install /tmp/bagccgop /usr/local/bin
```

You can now start the cross-compilation process:

```shell
$ bagccgop -j "$(nproc)" -b hello-world main.go
2021/10/24 17:15:01 building linux/arm (out/hello-world.linux-armv6l)
2021/10/24 17:15:01 building linux/alpha (out/hello-world.linux-alpha)
2021/10/24 17:15:01 building linux/sparc64 (out/hello-world.linux-sparc64)
2021/10/24 17:15:01 building linux/ppc (out/hello-world.linux-powerpc)
2021/10/24 17:15:01 building linux/amd64 (out/hello-world.linux-x86_64)
2021/10/24 17:15:01 building linux/ppc64 (out/hello-world.linux-ppc64)
2021/10/24 17:15:01 building linux/arm64 (out/hello-world.linux-aarch64)
2021/10/24 17:15:01 building linux/riscv64 (out/hello-world.linux-riscv64)
2021/10/24 17:15:02 building linux/arm (out/hello-world.linux-armv7l)
2021/10/24 17:15:03 building linux/386 (out/hello-world.linux-i686)
2021/10/24 17:15:03 building linux/mipsle (out/hello-world.linux-mips)
2021/10/24 17:15:03 building linux/mips64le (out/hello-world.linux-mips64)
2021/10/24 17:15:03 building linux/ppc64le (out/hello-world.linux-ppc64le)
2021/10/24 17:15:03 building linux/s390x (out/hello-world.linux-s390x)
2021/10/24 17:15:04 could not build for platform linux/alpha: err=could not install packages: err=exit status 2, stdout=, stderr=# command-line-arguments
cgo: cannot load DWARF output from $WORK/b001//_cgo_.o: applyRelocations: not implemented
```

As you can see, we get an error for the `linux/alpha` platform. We decide we don't want to support `linux/alpha`, so let's re-run the command with these platforms disabled:

```shell
$ bagccgop -j "$(nproc)" -b hello-world -x '(linux/alpha)' main.go
2021/10/24 17:16:06 building linux/arm (out/hello-world.linux-armv6l)
2021/10/24 17:16:06 building linux/ppc64 (out/hello-world.linux-ppc64)
2021/10/24 17:16:06 building linux/sparc64 (out/hello-world.linux-sparc64)
2021/10/24 17:16:06 building linux/ppc (out/hello-world.linux-powerpc)
2021/10/24 17:16:06 building linux/amd64 (out/hello-world.linux-x86_64)
2021/10/24 17:16:06 skipping linux/alpha (platform matched the provided regex)
2021/10/24 17:16:06 building linux/arm64 (out/hello-world.linux-aarch64)
2021/10/24 17:16:06 building linux/arm (out/hello-world.linux-armv7l)
2021/10/24 17:16:06 building linux/riscv64 (out/hello-world.linux-riscv64)
2021/10/24 17:16:07 building linux/386 (out/hello-world.linux-i686)
2021/10/24 17:16:07 building linux/mipsle (out/hello-world.linux-mips)
2021/10/24 17:16:07 building linux/mips64le (out/hello-world.linux-mips64)
2021/10/24 17:16:08 building linux/ppc64le (out/hello-world.linux-ppc64le)
2021/10/24 17:16:08 building linux/s390x (out/hello-world.linux-s390x)
```

If we now check the `out` directory, we can see that we now have successfully built binaries for all supported platforms:

```shell
$ file out/*
out/hello-world.linux-aarch64: ELF 64-bit LSB pie executable, ARM aarch64, version 1 (SYSV), dynamically linked, interpreter /lib/ld-linux-aarch64.so.1, for GNU/Linux 3.7.0, with debug_info, not stripped
out/hello-world.linux-armv6l:  ELF 32-bit LSB pie executable, ARM, EABI5 version 1 (SYSV), dynamically linked, interpreter /lib/ld-linux.so.3, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-armv7l:  ELF 32-bit LSB pie executable, ARM, EABI5 version 1 (SYSV), dynamically linked, interpreter /lib/ld-linux-armhf.so.3, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-i686:    ELF 32-bit LSB pie executable, Intel 80386, version 1 (SYSV), dynamically linked, interpreter /lib/ld-linux.so.2, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-mips:    ELF 32-bit LSB pie executable, MIPS, MIPS32 rel2 version 1 (SYSV), dynamically linked, interpreter /lib/ld.so.1, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-mips64:  ELF 64-bit LSB pie executable, MIPS, MIPS64 rel2 version 1 (SYSV), dynamically linked, interpreter /lib64/ld.so.1, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-powerpc: ELF 32-bit MSB pie executable, PowerPC or cisco 4500, version 1 (SYSV), dynamically linked, interpreter /lib/ld.so.1, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-ppc64:   ELF 64-bit MSB pie executable, 64-bit PowerPC or cisco 7500, version 1 (SYSV), dynamically linked, interpreter /lib64/ld64.so.1, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-ppc64le: ELF 64-bit LSB pie executable, 64-bit PowerPC or cisco 7500, version 1 (SYSV), dynamically linked, interpreter /lib64/ld64.so.2, for GNU/Linux 3.10.0, with debug_info, not stripped
out/hello-world.linux-riscv64: ELF 64-bit LSB pie executable, UCB RISC-V, version 1 (SYSV), dynamically linked, interpreter /lib/ld-linux-riscv64-lp64d.so.1, for GNU/Linux 4.15.0, with debug_info, not stripped
out/hello-world.linux-s390x:   ELF 64-bit MSB pie executable, IBM S/390, version 1 (SYSV), dynamically linked, interpreter /lib/ld64.so.1, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-sparc64: ELF 64-bit MSB pie executable, SPARC V9, relaxed memory ordering, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux.so.2, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-x86_64:  ELF 64-bit LSB pie executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, for GNU/Linux 3.2.0, with debug_info, not stripped
```

Now, let's a few compiler flags to make the build binaries fully static; we can do this by setting the `GOFLAGS` env variable:

```shell
$ GOFLAGS='-gccgoflags=-static' bagccgop -j "$(nproc)" -b hello-world -x '(linux/alpha)' main.go
2021/10/24 17:18:00 building linux/arm (out/hello-world.linux-armv6l)
2021/10/24 17:18:00 building linux/sparc64 (out/hello-world.linux-sparc64)
2021/10/24 17:18:00 building linux/ppc (out/hello-world.linux-powerpc)
2021/10/24 17:18:00 skipping linux/alpha (platform matched the provided regex)
2021/10/24 17:18:00 building linux/arm (out/hello-world.linux-armv7l)
2021/10/24 17:18:00 building linux/ppc64 (out/hello-world.linux-ppc64)
2021/10/24 17:18:00 building linux/amd64 (out/hello-world.linux-x86_64)
2021/10/24 17:18:00 building linux/riscv64 (out/hello-world.linux-riscv64)
2021/10/24 17:18:00 building linux/arm64 (out/hello-world.linux-aarch64)
2021/10/24 17:18:02 building linux/386 (out/hello-world.linux-i686)
2021/10/24 17:18:03 building linux/mipsle (out/hello-world.linux-mips)
2021/10/24 17:18:03 building linux/mips64le (out/hello-world.linux-mips64)
2021/10/24 17:18:03 building linux/ppc64le (out/hello-world.linux-ppc64le)
2021/10/24 17:18:03 building linux/s390x (out/hello-world.linux-s390x)
```

If we now check the output again, you can see that the binaries are now fully static:

```shell
$ file out/*
out/hello-world.linux-aarch64: ELF 64-bit LSB executable, ARM aarch64, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.7.0, with debug_info, not stripped
out/hello-world.linux-armv6l:  ELF 32-bit LSB executable, ARM, EABI5 version 1 (GNU/Linux), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-armv7l:  ELF 32-bit LSB executable, ARM, EABI5 version 1 (GNU/Linux), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-i686:    ELF 32-bit LSB executable, Intel 80386, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-mips:    ELF 32-bit LSB executable, MIPS, MIPS32 rel2 version 1 (SYSV), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-mips64:  ELF 64-bit LSB executable, MIPS, MIPS64 rel2 version 1 (SYSV), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-powerpc: ELF 32-bit MSB executable, PowerPC or cisco 4500, version 1 (SYSV), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-ppc64:   ELF 64-bit MSB executable, 64-bit PowerPC or cisco 7500, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-ppc64le: ELF 64-bit LSB executable, 64-bit PowerPC or cisco 7500, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.10.0, with debug_info, not stripped
out/hello-world.linux-riscv64: ELF 64-bit LSB executable, UCB RISC-V, version 1 (SYSV), statically linked, for GNU/Linux 4.15.0, with debug_info, not stripped
out/hello-world.linux-s390x:   ELF 64-bit MSB executable, IBM S/390, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-sparc64: ELF 64-bit MSB executable, SPARC V9, Sun UltraSPARC1 Extensions Required, relaxed memory ordering, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
out/hello-world.linux-x86_64:  ELF 64-bit LSB executable, x86-64, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
```

🚀 **That's it!** We've successfully added support for most Debian ports to this app.

If you're enjoying bagccgop, the following projects might also be of help to you too:

- Need to cross-compile with additional packages, such as OpenSSL, SDL2 or SQLite? Check out the [Reference](#reference)!
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
  -v, --verbose                  Enable logging of executed commands
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

bagccgop (c) 2021 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
