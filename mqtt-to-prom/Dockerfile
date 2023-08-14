FROM golang:1.19.3-alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build

# Cache the dependencies: https://github.com/montanaflynn/golang-docker-cache
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

# Compile statically
COPY . .
RUN go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

# Build the final image
FROM scratch as final

COPY --from=builder /build/main /app/

# Add the config
COPY config.ini config.ini

WORKDIR /app

ENTRYPOINT ["./main"]