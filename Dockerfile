# Builder for go binaries
FROM golang:alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod ./
COPY go.sum ./

RUN go mod download

#Copy all .go files into container
COPY *.go ./

#Copy other package contents as well
COPY ./business ./business/
COPY ./core ./core/
COPY ./infrastructure ./infrastructure/

RUN go build -o /marketmoogle-docker

EXPOSE 3000

CMD [ "/marketmoogle-docker" ]