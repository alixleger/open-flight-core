FROM golang:alpine as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .
RUN cp /build/.env .


FROM scratch
COPY --from=builder /dist /
ENTRYPOINT ["/main"]
