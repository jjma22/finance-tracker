FROM golang:1.26 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build -o /app/finance-tracker ./cmd/finance-api

FROM alpine:3.14
COPY --from=build /app/finance-tracker /bin/

EXPOSE 9090

CMD ["/bin/finance-tracker"]