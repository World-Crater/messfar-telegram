#!/bin/bash

docker build --platform linux/amd64 . -t=superj80820/messfar-telegram-service
docker push superj80820/messfar-telegram-service