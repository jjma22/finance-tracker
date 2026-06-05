FROM golang:1.26 AS build
WORKDIR /app

COPY go.mod go.sum ./
#Download  dependencies
RUN go mod download

COPY ./ ./

#Compiles and creates binary
RUN CGO_ENABLED=0 go build -o /app/finance-tracker ./cmd/finance-api

#Start second stage
FROM alpine:3.14
#Copy only compiled binary to new image
COPY --from=build /app/finance-tracker /bin/

EXPOSE 9090

#Switch to non-root user
RUN adduser -S finance-app
USER finance-app


CMD ["/bin/finance-tracker"]