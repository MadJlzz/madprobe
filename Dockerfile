FROM golang:1.13 as builder
WORKDIR /build/
COPY . .
RUN go get -v -t -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o madprobe .

FROM alpine:latest
RUN addgroup -S golang && adduser -S golang -G golang
USER golang:golang
WORKDIR /opt/madprobe/
COPY --from=builder /build/madprobe .
CMD ["./madprobe"]
