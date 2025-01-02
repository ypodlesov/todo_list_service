FROM ubuntu:22.04

LABEL org.opencontainers.image.source=https://github.com/ypodlesov/todo_list_service

# install packages with apt
RUN apt-get update
RUN apt-get install -y \
    wget \
    systemd \
    sudo \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update

# installing go
ENV GO_VERSION=1.23.1

RUN wget https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
    rm go$GO_VERSION.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

ARG BASE_IMAGE_DIR
ENV BASE_IMAGE_DIR=$BASE_IMAGE_DIR

COPY $BASE_IMAGE_DIR/tg_bot_api_install.sh /
RUN chmod +x /tg_bot_api_install.sh
RUN /tg_bot_api_install.sh
RUN rm -rf /tg_bot_api_install.sh