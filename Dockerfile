# # Docker image versions
# ARG alpine=3.11
# ARG go=1.14.0

# FROM golang:${go}-alpine${alpine} AS buildenv

# LABEL maintainer='kathleen french <kfrench09@gmail.com>'

# ENV GO111MODULE=on
# ENV CGO_ENABLED=0
# ENV LOCAL_NETWORK=172.16.0.0/12

# ENV APP_NAME sneak
# ENV INSTALL_PATH /opt/${APP_NAME}
# ENV APP_PERSIST_DIR /opt/${APP_NAME}_data

# WORKDIR /go/src/github.com/kathleenfrench/sneak

# RUN apk update \
#     && apk add --no-cache vim curl git make openvpn easy-rsa bash netcat-openbsd zip dumb-init privoxy

# WORKDIR ${APP_INSTALL_PATH}


# RUN mkdir -p /dev/net \ 
#     && mknod /dev/net/tun c 10 200

# RUN apk add --no-cache openvpn easy-rsa bash netcat-openbsd zip dumb-init && \
#     mkdir -p ${APP_PERSIST_DIR}


# EXPOSE 8118 4444


###

ARG SNEAK_VERSION
ARG LOCAL_NETWORK=172.16.0.0/12
ARG USER=sneak
ARG UID=1000
ARG GID=1000

FROM golang:alpine as sneak-builder
ARG USER
ARG UID
ARG GID

RUN export uid=${UID} gid=${GID} username=${USER} && \
    mkdir -p /home/${USER} && \
    echo "${username}:x:${uid}:${gid}:User,,,:/home/${username}:/bin/bash" >> /etc/passwd && \
    echo "${username}:x:${uid}:" >> /etc/group && \
    chown ${uid}:${gid} -R /home/${username}

RUN mkdir -p /dev/net \ 
    && mknod /dev/net/tun c 10 200

WORKDIR /go/src/github.com/kathleenfrench/sneak
RUN apk update && apk add --no-cache openvpn easy-rsa bash netcat-openbsd zip dumb-init privoxy curl git tmux
COPY . .

FROM sneak-builder as sneak-build
ARG SNEAK_VERSION
ENV GO111MODULE=on
WORKDIR /go/src/github.com/kathleenfrench/sneak
RUN go build -mod=download -a --installsuffix cgo -ldflags "-X github.com/kathleenfrench/sneak/cmd/sneak.Version=${SNEAK_VERSION}" -o app

FROM alpine:3.11.3 as run-env
ARG USER
ARG LOCAL_NETWORK

COPY --from=sneak-build /go/src/github.com/kathleenfrench/sneak/app /usr/local/bin/sneak
RUN chmod +ugox /usr/local/bin/sneak

USER ${USER}

ENV LOCAL_NETWORK ${LOCAL_NETWORK}
ENV PATH "$PATH:/usr/local/bin/sneak"

WORKDIR /home/${USER}/

ENTRYPOINT [ "/bin/bash" ]

EXPOSE 8118 4444
# EXPOSE 8080/tcp

# EXPOSE 1194/udp
# EXPOSE 8080/tcp