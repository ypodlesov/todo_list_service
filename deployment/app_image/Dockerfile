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

# Run app
WORKDIR /app

COPY pkg/storage/postgres/migrations ./migrations
COPY configs/prod.yaml .

# Environments
ARG CONFIG_PATH="/app/prod.yaml"
ENV CONFIG_PATH=$CONFIG_PATH

CMD ["/bin/todo_list_service"]