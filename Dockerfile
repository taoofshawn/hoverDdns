FROM golang:buster as builder
WORKDIR /hoverDdns
RUN go get -d -v  github.com/taoofshawn/hoverDdns

FROM gcr.io/distroless/base-debian10 as runner
COPY --from=downloader /hoverDdns /hoverDdns

CMD ["/hoverDdns"]