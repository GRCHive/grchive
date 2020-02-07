FROM debian:buster-slim
RUN mkdir -p /usr/share/man/man1mkdir -p /usr/share/man/man1
RUN apt-get update && apt-get install -y --no-install-recommends --no-install-suggests \
    gnupg2 \
    curl \
    ca-certificates \
    lsb-release \
    software-properties-common && \
    rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/debian \
      $(lsb_release -cs) \
         stable"
RUN export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)" && \
    echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - 

RUN apt-get update && apt-get install -y --no-install-recommends --no-install-suggests\
    openjdk-11-jdk \
    apt-transport-https \
    ca-certificates \
    curl \
    unzip \
    openssh-client \
    openssl \
    gettext-base \
    python \
    python3-pip \
    python3-distutils \
    imagemagick \ 
    jpegoptim \
    docker-ce \
    docker-ce-cli \
    google-cloud-sdk \
    build-essential \
    git \
    zip \
    zlib1g \
    zlib1g-dev && rm -rf /var/lib/apt/lists/*

COPY docker-entrypoint.sh /usr/local/bin/
ENV DOCKER_TLS_CERTDIR=/certs
RUN mkdir /certs /certs/client && chmod 1777 /certs /certs/client

# Rely on the CI system to copy these files into this directory.
COPY bootstrap_bazel.sh /
COPY download_libffi.sh /
COPY download_python.sh /
COPY python-helper/ /python-helper/
COPY download_terraform.sh /
COPY download_flyway.sh /
COPY download_kubectl.sh /

WORKDIR /
RUN /bootstrap_bazel.sh && mv /bazel/output /bazel-bin && rm -rf /bazel
RUN /download_libffi.sh
RUN /download_python.sh && mv /python/python-3.8.1/bin /python-bin && rm -rf /python && mkdir -p /python/python-3.8.1 && mv /python-bin /python/python-3.8.1/bin
RUN /download_terraform.sh
RUN /download_flyway.sh
RUN /download_kubectl.sh

ENV PATH="/bazel-bin:${PATH}"
ENV PATH="/terraform:${PATH}"
ENV PATH="/flyway:${PATH}"
ENV PATH="/kubectl:${PATH}"

RUN mkdir -p /gcloud
RUN curl -o /gcloud/cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64
RUN chmod +x /gcloud/cloud_sql_proxy
ENV PATH="/gcloud:${PATH}"

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["bash"]
