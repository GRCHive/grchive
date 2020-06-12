FROM nginx:1.19.0

RUN apt-get update && \
    apt-get install -y certbot \
        python3-certbot-nginx

COPY entrypoint.sh /
COPY nginx.conf /etc/nginx/nginx.conf

ENTRYPOINT ["/entrypoint.sh"]
