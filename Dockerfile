FROM golang:1.9.4 AS go

RUN apt update -y

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD . /go/src/go.dutchsec.com/lootbox

ARG LDFLAGS=""

WORKDIR /go/src/go.dutchsec.com/lootbox
RUN go build -ldflags="$(go run scripts/gen-ldflags.go)" -o /go/bin/app go.dutchsec.com/lootbox

FROM debian

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=go /go/bin/app /lootbox/lootbox

ARG LDFLAGS=""
RUN mkdir /loot/

ENTRYPOINT ["/lootbox/lootbox"]
