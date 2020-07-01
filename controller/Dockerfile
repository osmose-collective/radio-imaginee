# build
FROM		golang:1.12-alpine as builder
RUN		apk add --no-cache git
ENV		GO111MODULE=on
WORKDIR         /go/src/github.com/osmose-collective/radio-imaginee/controller
COPY            go.mod go.sum ./
RUN		go mod download
COPY            . .
RUN		go get .
RUN             go build -ldflags '-extldflags "-static"' -tags netgo -v -o /go/bin/controller

# minimal runtime
FROM            alpine:3.12
RUN		apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR		/workdir
COPY            --from=builder /go/bin/controller /bin/controller
COPY		--from=builder /go/src/github.com/osmose-collective/radio-imaginee/controller/www www
COPY		--from=builder /etc/ssl /etc/
ENTRYPOINT      ["/bin/controller"]
#CMD             []
EXPOSE          8042
#HEALTHCHECK     CMD /bin/controller healthcheck --wait