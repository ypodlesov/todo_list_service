#!/bin/bash

# telegram bot api set local
apt-get update && apt install -y \
    openssl \
    zlib1g-dev \
    gcc \
    g++ \
    gperf \
    cmake

git clone --recursive https://github.com/tdlib/telegram-bot-api.git
cd telegram-bot-api
mkdir build
cd build
cmake -DCMAKE_BUILD_TYPE=Release ..
cmake --build . --target install

cd $APP_DIR
rm -rf telegram-bot-api