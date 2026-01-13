FROM golang:1.25

WORKDIR /app

# cache go modules
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org && go mod download

# copy sources
COPY . .

EXPOSE 8080

ENV CONFIG_PATH=/app/example/config.dev.yaml
ENV GOPATH=/go

CMD ["sh", "-c", "export configPath=$CONFIG_PATH && go run ./cmd/app"]
