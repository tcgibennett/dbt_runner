#
# NOTE: THIS DOCKERFILE IS GENERATED VIA "apply-templates.sh"
#
# PLEASE DO NOT EDIT IT DIRECTLY.
#

FROM ubuntu:latest

ENV VIRTUAL_ENV=/opt/venv
RUN echo "alias env_dbt='source ~/dbt-env/bin/activate'" >> ~/.bashrc
RUN apt-get update && apt-get install -y curl wget gpg git libpq-dev python-dev python3-pip python3.8-venv
RUN apt-get remove python-cffi
RUN pip install --upgrade cffi
RUN pip install cryptography~=3.4
RUN python3 -m venv ~/dbt-env
ENV PATH="~/dbt-env/bin:$PATH"
RUN pip install dbt

RUN dbt --version
# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
#RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

ENV PATH /usr/local/go/bin:$PATH

ENV GOLANG_VERSION 1.16.6

RUN export \
    # set GOROOT_BOOTSTRAP such that we can actually build Go
    GOROOT_BOOTSTRAP="$(go env GOROOT)" \
    # ... and set "cross-building" related vars to the installed system's values so that we create a build targeting the proper arch
    # (for example, if our build host is GOARCH=amd64, but our build env/image is GOARCH=386, our build needs GOARCH=386)
    GOOS="$(go env GOOS)" \
    GOARCH="$(go env GOARCH)" \
    GOHOSTOS="$(go env GOHOSTOS)" \
    GOHOSTARCH="$(go env GOHOSTARCH)" \
    ; \
    \
    # https://github.com/golang/go/issues/38536#issuecomment-616897960
    url='https://golang.org/dl/go1.16.6.linux-amd64.tar.gz'; \
    sha256='be333ef18b3016e9d7cb7b1ff1fdb0cac800ca0be4cf2290fe613b3d069dfe0d'; \
    \
    wget -O go.tgz.asc "$url.asc"; \
    wget -O go.tgz "$url"; \
    echo "$sha256 *go.tgz" | sha256sum -c -; \
    \
    # https://github.com/golang/go/issues/14739#issuecomment-324767697
    export GNUPGHOME="$(mktemp -d)"; \
    # https://www.google.com/linuxrepositories/
    gpg --batch --keyserver keyserver.ubuntu.com --recv-keys 'EB4C 1BFD 4F04 2F6D DDCC EC91 7721 F63B D38B 4796'; \
    gpg --batch --verify go.tgz.asc go.tgz; \
    gpgconf --kill all; \
    rm -rf "$GNUPGHOME" go.tgz.asc; \
    \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz; \
    \
    goEnv="$(go env | sed -rn -e '/^GO(OS|ARCH|ARM|386)=/s//export \0/p')"; \
    eval "$goEnv"; \
    [ -n "$GOOS" ]; \
    [ -n "$GOARCH" ]; \
    ( \
    cd /usr/local/go/src; \
    ./make.bash; \
    ); \
    \
    #apk del --no-network .build-deps; \
    \
    # pre-compile the standard library, just like the official binary release tarballs do
    go install std; \
    # go install: -race is only supported on linux/amd64, linux/ppc64le, linux/arm64, freebsd/amd64, netbsd/amd64, darwin/amd64 and windows/amd64
    #	go install -race std; \
    \
    # remove a few intermediate / bootstrapping files the official binary release tarballs do not contain
    rm -rf \
    /usr/local/go/pkg/*/cmd \
    /usr/local/go/pkg/bootstrap \
    /usr/local/go/pkg/obj \
    /usr/local/go/pkg/tool/*/api \
    /usr/local/go/pkg/tool/*/go_bootstrap \
    /usr/local/go/src/cmd/dist/dist \
    ; \
    \
    go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
#RUN git config --global user.email GIT_AUTHOR_EMAIL
#RUN git config --global user.name GIT_AUTHOR_NAME
COPY ./runner /go/bin
WORKDIR $GOPATH
