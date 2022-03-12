FROM golang:alpine

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /output/barcode-generator main.go

###############################################################################

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=0 /src/data/ /data/
COPY --from=0 /output/barcode-generator /
CMD ["/barcode-generator"]
