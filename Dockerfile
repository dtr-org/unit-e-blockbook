FROM debian:9

ARG COIN
ARG RPC_ADDRESS

ENV RUNTIME_DEPS nano libbz2-1.0 libc6 libgcc1 libgflags2v5 liblz4-1 libsnappy1v5 libstdc++6 libzmq5 zlib1g init-system-helpers coreutils passwd findutils psmisc

ADD build/blockbook*.deb /root/
RUN apt update && \
	apt install -y $RUNTIME_DEPS && \
	dpkg --ignore-depends backend-${COIN} -i /root/*.deb

RUN sed -i 's/"rpc_url": "[^\"].*"/"rpc_url": "http:\/\/'${RPC_ADDRESS}'"/' /opt/coins/blockbook/${COIN}/config/blockchaincfg.json

EXPOSE 9172 9272

WORKDIR /opt/coins/blockbook/${COIN}
CMD ["./bin/blockbook", "-blockchaincfg=config/blockchaincfg.json", "-datadir=db", "-sync", "-internal=:9172", "-public=:9272", "-certfile=cert/blockbook", "-explorer=", "-log_dir=logs", "-dbcache=1073741824"]
