FROM debian:sid

ENV DEBIAN_ARCHITECTURES="alpha powerpc ppc64 sparc64 riscv64 amd64 arm64 armel armhf i386 mipsel mips64el ppc64el s390x"
ENV APT_PACKAGES="gccgo-alpha-linux-gnu gcc-alpha-linux-gnu gccgo-powerpc-linux-gnu gcc-powerpc-linux-gnu gccgo-powerpc64-linux-gnu gcc-powerpc64-linux-gnu gccgo-sparc64-linux-gnu gcc-sparc64-linux-gnu gccgo-riscv64-linux-gnu gcc-riscv64-linux-gnu gccgo gcc gccgo-aarch64-linux-gnu gcc-aarch64-linux-gnu gccgo-arm-linux-gnueabi gcc-arm-linux-gnueabi gccgo-arm-linux-gnueabihf gcc-arm-linux-gnueabihf gccgo-i686-linux-gnu gcc-i686-linux-gnu gccgo-mipsel-linux-gnu gcc-mipsel-linux-gnu gccgo-mips64el-linux-gnuabi64 gcc-mips64el-linux-gnuabi64 gccgo-powerpc64le-linux-gnu gcc-powerpc64le-linux-gnu gccgo-s390x-linux-gnu gcc-s390x-linux-gnu"

RUN apt update
RUN apt install -y ca-certificates debian-ports-archive-keyring
RUN printf "deb http://ftp.ports.debian.org/debian-ports unstable main\ndeb http://ftp.ports.debian.org/debian-ports unreleased main\ndeb http://ftp.ports.debian.org/debian-ports experimental main" >>/etc/apt/sources.list
RUN for arch in ${DEBIAN_ARCHITECTURES}; do dpkg --add-architecture "${arch}"; done
ENV PATH="$PATH:/${HOME}/go/bin"
RUN apt update

RUN apt install -y golang ${APT_PACKAGES}
