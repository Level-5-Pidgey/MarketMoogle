FROM golang:1.21-alpine
LABEL authors="Carl Alexander"

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN go build -v -o /app/MarketMoogle

CMD ["/app/MarketMoogle"]