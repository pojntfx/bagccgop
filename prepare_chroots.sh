#!/bin/bash

set -e

prepare_chroot() {
    local DEBIAN_DIST="$1"
    local DEBIAN_MIRROR="$2"
    local GOLANG_ARCH="$3"
    local DEBIAN_ARCH="$4"
    local APT_PKG_SUFFIX="$5"
    local GCC_ARCH="$6"

    apt update
    apt install -y debootstrap
    mkdir -p "/var/lib/bagccgop/${DEBIAN_ARCH}-chroot/data"
    debootstrap "${DEBIAN_DIST}" "/var/lib/bagccgop/${DEBIAN_ARCH}-chroot" "${DEBIAN_MIRROR}"

    chroot "/var/lib/bagccgop/${DEBIAN_ARCH}-chroot" /bin/bash -c 'apt install -y ca-certificates debian-ports-archive-keyring'
    chroot "/var/lib/bagccgop/${DEBIAN_ARCH}-chroot" /bin/bash -c 'printf "deb http://ftp.ports.debian.org/debian-ports unstable main\ndeb http://ftp.ports.debian.org/debian-ports unreleased main\ndeb http://ftp.ports.debian.org/debian-ports experimental main" >>/etc/apt/sources.list'
    chroot "/var/lib/bagccgop/${DEBIAN_ARCH}-chroot" /bin/bash -c "dpkg --add-architecture \"${DEBIAN_ARCH}\""
    chroot "/var/lib/bagccgop/${DEBIAN_ARCH}-chroot" /bin/bash -c 'apt update'
    chroot "/var/lib/bagccgop/${DEBIAN_ARCH}-chroot" /bin/bash -c "apt install -y git golang \"gccgo${APT_PKG_SUFFIX}\" \"gcc${APT_PKG_SUFFIX}\""
    chroot "/var/lib/bagccgop/${DEBIAN_ARCH}-chroot" /bin/bash -c "echo 'export PATH=\"\${PATH}:\${HOME}/go/bin\"' >> ~/.bashrc"
}

prepare_chroots() {
    local TARGETS="$1"

    for TARGET in ${TARGETS}; do
        DEBIAN_DIST="$(cut -d'|' -f1 <<<"${TARGET}")"
        DEBIAN_MIRROR="$(cut -d'|' -f2 <<<"${TARGET}")"
        GOLANG_ARCH="$(cut -d'|' -f3 <<<"${TARGET}")"
        DEBIAN_ARCH="$(cut -d'|' -f4 <<<"${TARGET}")"
        APT_PKG_SUFFIX="$(cut -d'|' -f5 <<<"${TARGET}")"
        GCC_ARCH="$(cut -d'|' -f6 <<<"${TARGET}")"

        export DEBIAN_DIST
        export DEBIAN_MIRROR
        export GOLANG_ARCH
        export DEBIAN_ARCH
        export APT_PKG_SUFFIX
        export GCC_ARCH

        prepare_chroot "${DEBIAN_DIST}" "${DEBIAN_MIRROR}" "${GOLANG_ARCH}" "${DEBIAN_ARCH}" "${APT_PKG_SUFFIX}" "${GCC_ARCH}"
    done
}

prepare_chroots "$1"
