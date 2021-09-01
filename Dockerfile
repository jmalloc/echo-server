FROM scratch
ARG TARGETPLATFORM
COPY artifacts/build/release/$TARGETPLATFORM/echo-server /bin/echo-server
ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/bin/echo-server"]
