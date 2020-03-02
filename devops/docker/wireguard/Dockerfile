FROM debian:buster-slim
RUN echo "deb http://deb.debian.org/debian/ unstable main" > /etc/apt/sources.list.d/unstable.list  && \
    printf 'Package: *\nPin: release a=unstable\nPin-Priority: 90\n' > /etc/apt/preferences.d/limit-unstable && \
    apt-get update && \
    apt-get install -y  --no-install-recommends --no-install-suggests wireguard iproute2 && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY etc/wireguard/wg0.conf /etc/wireguard/wg0.conf
COPY run.sh /run.sh
RUN chown -v root:root /etc/wireguard/wg0.conf && chmod -v 600 /etc/wireguard/wg0.conf
CMD ["bash", "-c", "/run.sh"]
