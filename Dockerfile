# build stage
FROM golang:1.13.8-alpine AS build-env
ADD . /src
RUN cd /src && go build -o app

# final stage
FROM alpine
WORKDIR /app
COPY config.yaml.example /app/config.yaml
RUN mkdir images
COPY --from=build-env /src/app /app/
ENTRYPOINT ./app