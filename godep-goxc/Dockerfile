FROM ailispaw/ubuntu-essential

RUN apt-get -q update && \
    apt-get -q -y install --no-install-recommends curl ca-certificates git make && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

ENV GOVERSION=1.6 \
    GOROOT=/goroot \
    GOPATH=/gopath
ENV PATH ${PATH}:${GOROOT}/bin:${GOPATH}/bin

RUN mkdir -p ${GOROOT} && \
    curl https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz \
      | tar xvzf - -C ${GOROOT} --strip-components=1

RUN go get github.com/laher/goxc && \
    go get github.com/tools/godep
