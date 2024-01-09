FROM node:21 as nodebuild

WORKDIR /ui

COPY ui/ .

RUN npm install --network-timeout=100000 && \
    npm run build


FROM golang:1-alpine AS build

RUN apk --no-cache add git go-bindata

COPY . /go/src/git.nemunai.re/nemunaire/hathoris
COPY --from=nodebuild /ui/build /go/src/git.nemunai.re/nemunaire/hathoris/ui/build
WORKDIR /go/src/git.nemunai.re/nemunaire/hathoris
RUN go get && go generate && go build -ldflags="-s -w"


FROM alpine:3.18

EXPOSE 8080
CMD ["/srv/hathoris"]

COPY --from=build /go/src/git.nemunai.re/nemunaire/hathoris/hathoris /srv/hathoris
