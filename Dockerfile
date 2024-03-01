FROM debian

COPY ./bin/rssbridge /usr/local/bin

CMD ["rssbridge"]

EXPOSE 3000

