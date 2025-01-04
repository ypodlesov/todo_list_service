ARG BASE_IMAGE

FROM $BASE_IMAGE

# LABEL org.opencontainers.image.source=https://github.com/ypodlesov/todo_list_service

WORKDIR /build

COPY go.mod .
COPY cmd ./cmd
COPY pkg ./pkg

RUN go mod tidy
RUN go build -o /bin/todo_list_service cmd/todo_list_service/main.go

RUN rm -rf /build

# Environments
ARG TELEGRAM_API_TOKEN
ENV TELEGRAM_API_TOKEN=$TELEGRAM_API_TOKEN

ARG TELEGRAM_API_ID
ENV TELEGRAM_API_ID=$TELEGRAM_API_ID

ARG TELEGRAM_API_HASH
ENV TELEGRAM_API_HASH=$TELEGRAM_API_HASH

# Run app
WORKDIR /app

RUN telegram-bot-api --api-id=$TELEGRAM_API_ID --api-hash=$TELEGRAM_API_HASH --local &

CMD ["/bin/todo_list_service"]