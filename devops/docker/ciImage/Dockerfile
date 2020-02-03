FROM debian:10.2
RUN apt-get update && apt-get install -y openjdk-11-jdk \
    apt-transport-https \
    ca-certificates \
    curl \
    unzip \
    gnupg2 \
    software-properties-common \
    openssh-client \
    gettext-base \
    python
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/debian \
      $(lsb_release -cs) \
         stable"
RUN apt-get update && apt-get install -y docker-ce docker-ce-cli
RUN apt-get -y upgrade
RUN apt-get install -y zip

COPY docker-entrypoint.sh /usr/local/bin/
ENV DOCKER_TLS_CERTDIR=/certs
RUN mkdir /certs /certs/client && chmod 1777 /certs /certs/client

# Rely on the CI system to copy this file into this directory.
COPY bootstrap_bazel.sh /
WORKDIR /
RUN ["/bin/bash", "./bootstrap_bazel.sh"]
ENV PATH="/bazel/output:${PATH}"

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["bash"]
