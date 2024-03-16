FROM debian

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    update-ca-certificates

COPY ./bin/rssbridge /usr/local/bin

CMD ["rssbridge"]

EXPOSE 3000

