FROM debian:buster-slim
RUN mkdir -p /usr/share/man/man1mkdir -p /usr/share/man/man1
RUN apt-get update && apt-get install -y --no-install-recommends --no-install-suggests libmagic1 \
    libreoffice \
    ca-certificates && rm -rf /var/lib/apt/lists/*
