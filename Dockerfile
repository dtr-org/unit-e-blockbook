FROM debian:9

ARG COIN

ENV RUNTIME_DEPS nano libbz2-1.0 libc6 libgcc1 libgflags2v5 \
    liblz4-1 libsnappy1v5 libstdc++6 libzmq5 zlib1g \
    init-system-helpers coreutils passwd findutils psmisc

ADD build/blockbook*.deb /root/

RUN apt-get update && \
    apt-get install -y $RUNTIME_DEPS && \
    dpkg --ignore-depends backend-${COIN} -i /root/*.deb

EXPOSE 9172 9272

WORKDIR /opt/coins/blockbook/${COIN}

ADD docker_entrypoint.sh /.
RUN chmod 777 /docker_entrypoint.sh

ENTRYPOINT ["/docker_entrypoint.sh"]
