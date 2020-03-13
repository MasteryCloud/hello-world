FROM alpine:3.11.3
COPY img /opt/hello-world/img
COPY build/* /opt/hello-world
WORKDIR /opt/hello-world
ENTRYPOINT ["./main"]
