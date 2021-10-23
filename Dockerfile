FROM debian:sid

COPY prepare_chroots.sh /usr/local/bin/prepare_chroots.sh

RUN prepare_chroots.sh \
    "sid|http://deb.debian.org/debian|powerpc|powerpc|-powerpc-linux-gnu|powerpc-linux-gnu "
