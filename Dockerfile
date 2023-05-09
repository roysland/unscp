FROM golang:1.17.2-alpine3.14
LABEL maintainer="Øyvind Røysland <roysland@gmail.com>"
RUN apk add --no-cache gcc libc-dev

WORKDIR /app
COPY . .
RUN go mod download && go build -o main .
EXPOSE 8090
CMD ["./main"]