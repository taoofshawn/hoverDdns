FROM gradle:jdk11 as builder

WORKDIR /app
RUN git clone --branch java https://github.com/taoofshawn/hoverDdns.git . && \
    gradle build


FROM gcr.io/distroless/java:11 as runner

COPY --from=builder /build/libs/*.jar /app
WORKDIR /app

ENTRYPOINT ["hoverDdns.jar"]

