# hello-world

A simple Docker container serving a single web page. We use this image
for testing, i.e. [our scalable AWS Fargate setup](https://github.com/MasteryCloud/aws-fargate-playground).

![Screenshot](./screenshot.png)

## About this image

The image is based on [Alpine Linux container](https://hub.docker.com/_/alpine)
and runs a simple webserver serving a single page. The webserver is written in Go.

## How to use this image

To run the docker image:

```
docker run --rm --publish 8080:80 masterycloud/hello-world
```

Then open a browser and navigate to http://localhost:8080.