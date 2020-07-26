FROM golang:1.13-alpine

WORKDIR /go/src/inreach2marinetraffic
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV MAPSHARE_ID=""
ENV MAPSHARE_PASSWORD=""
ENV MAPSHARE_INTERVAL=60
ENV SMTP_HOST=""
ENV SMTP_PORT=""
ENV SMTP_USER=""
ENV SMTP_PASSWORD=""
ENV EMAIL_ADDRESS=""
ENV MMSI=""

CMD ["inreach2marinetraffic"]