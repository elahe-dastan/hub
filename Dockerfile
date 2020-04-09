FROM golang:alpine As builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /applifier

FROM alpine:latest

LABEL maintainer="Elahe Dastan <elahe.dstn@gmail.com>"

WORKDIR /root/

COPY --from=builder /applifier .

ENTRYPOINT ["./hub"]