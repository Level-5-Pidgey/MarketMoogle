# Builder for go binaries
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

#Copy all .go files into container
COPY *.go ./

#Copy other package contents as well
COPY ./business ./business/
COPY ./core ./core/
COPY ./infrastructure ./infrastructure/

RUN go build -o /sanctuary-docker

EXPOSE 8080

CMD [ "/sanctuary-docker" ]