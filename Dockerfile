FROM python:3-buster as downloader
WORKDIR /hoverDdns
RUN git clone https://github.com/taoofshawn/hoverDdns.git .

FROM python:3-slim-buster as runner
COPY --from=downloader /hoverDdns /hoverDdns
WORKDIR /hoverDdns
RUN pip install --upgrade pip && \
	pip install -r requirements.txt

CMD ["python", "hoverDdns.py"]