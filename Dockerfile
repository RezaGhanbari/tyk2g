# iron/go:dev is the alpine image with the go tools added
FROM iron/go:dev
WORKDIR /app
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV TYK_MONGO_HOST=""
ENV TYK_MONGO_DB="tyk_analytics"
ENV TYK_MONGO_COL="tyk_analytics"
ENV TYK_MONGO_PORT="27017"
ENV TYK_TIME="10000"
ENV TYK_GRAYLOG_HOST=""
ENV TYK_GRAYLOG_PORT=""

ENV SRC_DIR=$GOPATH/tyk2g
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR; go build -o myapp; cp myapp /app/
ENTRYPOINT ["./myapp"]
