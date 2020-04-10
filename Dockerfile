FROM golang:buster as builder
WORKDIR /hoverDdns
RUN go get -d -v github.com/taoofshawn/hoverDdns && \
    go build -o hoverDdns

FROM gcr.io/distroless/base-debian10 as runner
COPY --from=builder /hoverDdns /
CMD ["/hoverDdns"]
