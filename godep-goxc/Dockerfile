FROM debian:wheezy

RUN apt-get update -y && \
    apt-get install --no-install-recommends -y -q \
      curl \
      build-essential \
      ca-certificates \
      git mercurial bzr \
    && rm -rf /var/lib/apt/lists/*

ENV GOVERSION=1.5.2 \
    GOROOT=/goroot \
    GOPATH=/gopath
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin

RUN mkdir -p ${GOROOT} && \
    curl https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz \
      | tar xvzf - -C ${GOROOT} --strip-components=1

RUN go get github.com/laher/goxc && \
    go get github.com/tools/godep

CMD godep restore && goxc
