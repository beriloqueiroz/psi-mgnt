FROM  golang:1.22.4 as build

USER root
ENV TZ="America/Fortaleza"
WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o /handle /cmd/main.go

# FROM scratch

# COPY --from=build /handle /handle
# COPY .env .env

EXPOSE 8000

ENTRYPOINT ["/handle"]