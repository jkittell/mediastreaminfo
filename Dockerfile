FROM --platform=$BUILDPLATFORM golang:alpine AS builder
WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY . ./
ARG TARGETOS
ARG TARGETARCH
ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH
RUN go build -o mediastreaminfo.exe ./cmd

FROM alpine:latest
RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg
COPY --from=builder ["/build/mediastreaminfo.exe", "/"]
EXPOSE 3000
RUN chmod +x /mediastreaminfo.exe
ENTRYPOINT [ "/mediastreaminfo.exe" ]

