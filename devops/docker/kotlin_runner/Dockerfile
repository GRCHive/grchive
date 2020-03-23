FROM debian:buster-slim

RUN apt-get update && apt-get install -y software-properties-common \
        curl \
        unzip \
        wget \
        gnupg2 \
        gosu && \
        rm -rf /var/lib/apt/lists/*
COPY download_install_kotlin.sh /
RUN /download_install_kotlin.sh

RUN mkdir -p /usr/share/man/man1
RUN wget -qO - https://adoptopenjdk.jfrog.io/adoptopenjdk/api/gpg/key/public | apt-key add - && \
    add-apt-repository --yes https://adoptopenjdk.jfrog.io/adoptopenjdk/deb/ && \
    apt-get update && apt-get install -y adoptopenjdk-8-hotspot && \
    rm -rf /var/lib/apt/lists/*

ENV PATH="/kotlinc/bin:${PATH}"

COPY entrypoint.sh /
COPY run.py /

RUN groupadd -r kotlin && useradd -r -g kotlin kotlin
RUN mkdir /data && chown -R kotlin /data
RUN mkdir /output && chown -R kotlin /output
WORKDIR /data

VOLUME ["/data"]
ENTRYPOINT ["/entrypoint.sh"]
