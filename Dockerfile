FROM golang:1.13-alpine3.11 AS build-env
RUN apk --no-cache add build-base git make
COPY . /src
WORKDIR /src
RUN make

FROM alpine:3.11
RUN apk --no-cache add su-exec
COPY --from=build-env /src/dist/peruse /bin/
CMD su-exec nobody /bin/peruse serv
