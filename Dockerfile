FROM python:3.7-buster as downloader
WORKDIR /hoverDdns
RUN git clone https://github.com/taoofshawn/hoverDdns.git . && \
	pip install -r requirements.txt

FROM gcr.io/distroless/python3-debian10 as runner
COPY --from=downloader /hoverDdns /hoverDdns
COPY --from=downloader /usr/local/lib/python3.7/site-packages /usr/local/lib/python3.7/site-packages
ENV PYTHONPATH=/usr/local/lib/python3.7/site-packages
WORKDIR /hoverDdns

CMD ["hoverDdns.py"]