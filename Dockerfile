FROM golang:1.12.5

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.3/dep-linux-amd64 && chmod +x /usr/local/bin/dep

WORKDIR /go/src/taco
COPY . /go/src/taco/.
COPY taco/Gopkg.toml taco/Gopkg.lock ./

RUN dep ensure -vendor-only

ENTRYPOINT [ "make", "dev" ]

EXPOSE 8000