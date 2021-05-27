FROM golang:1.16.4-buster as builder

RUN apt-get update && apt-get install -y uuid-runtime unzip zip \
  && rm -rf /var/lib/apt/lists/*
COPY . /src
WORKDIR /src
RUN make makefiles; make release

FROM scratch
COPY --from=builder /src/artifacts/build/release/linux/amd64/echo-server /bin/echo-server
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/bin/echo-server"]
