FROM debian:buster
RUN apt-get update && apt-get install -y openssh-server

COPY entrypoint.sh /entrypoint.sh
COPY sshd_config /etc/ssh/sshd_config
RUN mkdir /root/.ssh
COPY id_rsa.pub /root/.ssh/authorized_keys

EXPOSE 22
ENTRYPOINT ["/entrypoint.sh"]
