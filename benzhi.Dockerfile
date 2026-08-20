FROM golang:1.23
WORKDIR /app
ENV GOPROXY=off GOSUMDB=off
COPY vendor ./vendor
COPY . .
RUN go build -mod=vendor ./...
CMD ["go", "run", "-mod=vendor", "./cmd/tracelink-demo"]
