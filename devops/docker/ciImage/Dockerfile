FROM debian:10.2
RUN apt-get update && apt-get install -y openjdk-11-jdk \
    apt-transport-https \
    ca-certificates \
    curl \
    unzip \
    gnupg2 \
    software-properties-common \
    openssh-client \
    openssl \
    gettext-base \
    python \
    python3-pip \
    python3-distutils
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/debian \
      $(lsb_release -cs) \
         stable"
RUN apt-get update && apt-get install -y docker-ce docker-ce-cli
RUN apt-get -y upgrade
RUN apt-get install -y zip
RUN apt-get install -y zlib1g zlib1g-dev

COPY docker-entrypoint.sh /usr/local/bin/
ENV DOCKER_TLS_CERTDIR=/certs
RUN mkdir /certs /certs/client && chmod 1777 /certs /certs/client

# Rely on the CI system to copy these files into this directory.
COPY bootstrap_bazel.sh /
COPY download_libffi.sh /
COPY download_python.sh /
WORKDIR /
RUN ["/bin/bash", "./bootstrap_bazel.sh"]
RUN ["/bin/bash", "./download_libffi.sh"]

COPY python-helper/ /python-helper/
RUN ["/bin/bash", "./download_python.sh"]
ENV PATH="/bazel/output:${PATH}"

COPY download_terraform.sh /
RUN ["/bin/bash", "./download_terraform.sh"]
ENV PATH="/terraform:${PATH}"

RUN export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)" && \
    echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    apt-get update -y && apt-get install google-cloud-sdk -y

RUN mkdir -p /gcloud
RUN curl -o /gcloud/cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64
RUN chmod +x /gcloud/cloud_sql_proxy
ENV PATH="/gcloud:${PATH}"

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["bash"]
