ARG BASE_IMAGE

FROM $BASE_IMAGE

ARG APP_IMAGE_DIR
ENV APP_IMAGE_DIR=$APP_IMAGE_DIR

LABEL org.opencontainers.image.source=https://github.com/ypodlesov/todo_list_service

WORKDIR /src

RUN set -xe

COPY src .

RUN rm -rf go.mod && rm -rf go.sum
RUN go mod init todo_list_service
RUN go mod tidy

RUN go build -o /bin/todo_list_service

COPY $APP_IMAGE_DIR/users.sh .
RUN chmod +x /src/users.sh
RUN /src/users.sh
RUN rm -rf /src/users.sh

# Environments
ARG APP_DIR=/app
ENV APP_DIR=$APP_DIR

ARG TELEGRAM_API_TOKEN
ENV TELEGRAM_API_TOKEN=$TELEGRAM_API_TOKEN

ARG TELEGRAM_API_ID
ENV TELEGRAM_API_ID=$TELEGRAM_API_ID

ARG TELEGRAM_API_HASH
ENV TELEGRAM_API_HASH=$TELEGRAM_API_HASH

WORKDIR $APP_DIR

# Run app
RUN telegram-bot-api --api-id=$TELEGRAM_API_ID --api-hash=$TELEGRAM_API_HASH --local &

CMD ["/bin/todo_list_service"]