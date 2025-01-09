ARG BASE_IMAGE

FROM $BASE_IMAGE

WORKDIR /build

COPY go.mod .
COPY cmd ./cmd
COPY pkg ./pkg

RUN go mod tidy
RUN go build -o /bin/todo_list_service cmd/todo_list_service/main.go

RUN rm -rf /build

# Environments
ARG CONFIG_PATH
ENV CONFIG_PATH=$CONFIG_PATH

ARG PG_HOST
ENV PG_HOST=$PG_HOST

ENV PG_PORT=5432

ARG PG_USER
ENV PG_USER=$PG_USER

ARG PG_PASSWORD
ENV PG_PASSWORD=$PG_PASSWORD

ARG PG_DB_NAME
ENV PG_DB_NAME=$PG_DB_NAME

ARG PG_MIGRATIONS_DIR
ENV PG_MIGRATIONS_DIR=$PG_MIGRATIONS_DIR

# Run app
WORKDIR /app

COPY pkg/storage/postgres/migrations $PG_MIGRATIONS_DIR
COPY configs/prod.yaml $CONFIG_PATH

CMD ["/bin/todo_list_service"]