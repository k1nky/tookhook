FROM golang:1.20 AS build

RUN apt-get update \
    && apt install unzip \
    && wget https://github.com/protocolbuffers/protobuf/releases/download/v27.2/protoc-27.2-linux-x86_64.zip \
    && unzip protoc-27.2-linux-x86_64.zip


WORKDIR /src

COPY . .

RUN make prepare
RUN make build \
    && make plugin

FROM alpine:3.19

COPY --from=build /src/build/tookhook /app/
COPY --from=build /src/build/pachca /app/

EXPOSE 8080

CMD ["/app/tookhook"]