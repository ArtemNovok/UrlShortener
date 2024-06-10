FROM golang:1.22 as build

RUN mkdir /app

WORKDIR /app

COPY .  /app

RUN CGO_ENABLED=0 go build -o ShortApp ./cmd/url-shortener

RUN chmod +x /app/ShortApp

FROM alpine

RUN mkdir /app

WORKDIR /app

COPY --from=build /app/ShortApp  /app

COPY --from=build /app/config /app

EXPOSE 8000

CMD [ "/app/ShortApp" ]