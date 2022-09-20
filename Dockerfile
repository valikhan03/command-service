FROM golang:1.18

WORKDIR /app
ADD /. /app

RUN go install

## needed to RUN <swagger generator>

RUN go build -o index
ENV PORT=8080

ENTRYPOINT ["/app/index"]