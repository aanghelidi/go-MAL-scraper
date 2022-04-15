FROM golang:1.17 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create a scraper user to avoid root user
ENV SCRAPER_USER=scraper-user
RUN groupadd --gid 1000 $SCRAPER_USER \
  && useradd --uid 1000 --gid $SCRAPER_USER --create-home $SCRAPER_USER

ENV SCRAPER_HOME=/home/$SCRAPER_USER

# Switch to non-root user
USER $SCRAPER_USER
WORKDIR $SCRAPER_HOME

# Copy file and buid app
COPY . .
RUN go build -o /go/bin main.go cleanUtils.go

FROM alpine:latest

# Alpine uses adduser and addgroup...
ENV SCRAPER_USER=scraper-user
RUN addgroup -g 1000 $SCRAPER_USER \
  && adduser -u 1000 -G $SCRAPER_USER -h /home/$SCRAPER_USER -D $SCRAPER_USER
ENV SCRAPER_HOME=/home/$SCRAPER_USER

# Create tmp data volume
RUN mkdir -p /tmp/data
RUN chown -R $SCRAPER_USER:users /tmp

# Switch to non-root user
USER $SCRAPER_USER
WORKDIR $SCRAPER_HOME
COPY --from=builder /go/bin ./
ENTRYPOINT [ "./main" ]
