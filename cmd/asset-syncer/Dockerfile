FROM golang:1.13 as builder
COPY . /go/src/github.com/suomitek/suomitek-appboard
WORKDIR /go/src/github.com/suomitek/suomitek-appboard
RUN CGO_ENABLED=0 go build -a -installsuffix cgo ./cmd/asset-syncer

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/suomitek/suomitek-appboard/asset-syncer /asset-syncer
CMD ["/asset-syncer"]
