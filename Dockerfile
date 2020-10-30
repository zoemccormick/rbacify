FROM golang

ADD . /go/src/github.com/zoemccormick/rbacify

RUN go install github.com/zoemccormick/rbacify

ENTRYPOINT /go/bin/rbacify

EXPOSE 8000