FROM scratch
COPY artifacts/build/release/linux/amd64/echo-server /bin/echo-server
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/bin/echo-server"]
