FROM gradle:jdk14 as builder

WORKDIR /app

RUN git clone --branch java https://github.com/taoofshawn/hoverDdns.git .

ENTRYPOINT ["tail", "-f", "/dev/null"]

