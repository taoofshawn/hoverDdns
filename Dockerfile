FROM golang:buster as builder
RUN git clone https://github.com/taoofshawn/hoverDdns.git /hoverDdns && \
    cd /hoverDdns && \
    go build

FROM gcr.io/distroless/base as runner
COPY --from=builder /hoverDdns/hoverDdns /
CMD ["/hoverDdns"]