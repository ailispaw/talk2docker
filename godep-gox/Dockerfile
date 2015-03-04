FROM tcnksm/gox:1.4.2

RUN go get github.com/tools/godep

CMD godep restore && gox
