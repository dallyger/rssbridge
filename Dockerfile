FROM golang AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY rssbridge.go ./
RUN go build .

FROM debian

COPY --from=build /app/rssbridge /usr/local/bin

CMD ["rssbridge"]

EXPOSE 3000

