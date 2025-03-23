FROM mirror.gcr.io/golang:1.22.12 AS build

RUN apt-get update \
    && apt install unzip \
    && wget https://github.com/protocolbuffers/protobuf/releases/download/v27.2/protoc-27.2-linux-x86_64.zip \
    && unzip protoc-27.2-linux-x86_64.zip


WORKDIR /src

COPY . .

RUN make prepare
RUN make addplugins
RUN make build \
    && make plugin

FROM mirror.gcr.io/alpine:3.19

LABEL org.opencontainers.image.source="https://github.com/k1nky/tookhook"
LABEL org.opencontainers.image.description="TookHook is a webhook server with pluggable handlers."
LABEL org.opencontainers.image.licenses="Apache 2.0"

RUN apk add --no-cache tzdata

COPY --from=build /src/build/* /app/

EXPOSE 8080

CMD ["/app/tookhook", "run"]