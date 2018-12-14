# Unit-e Blockbook

This README is specifically for instructions on running Blockbook for the Unit-e coin.

## Running Unit-e node for Blockbook

Make sure the node is indexing transactions, has RPC server enabled, the RPC credentials are rpc:rpc and that node accepts incoming connections from Blockbook's IP.

Example flags:
```
united -txindex -server=1 -rpcbind=127.0.0.1  -rpcallowip=127.0.0.1/32 \
       -rpcauth=rpc:68c9195bf83979cc3bc3717cf17e61ce$de438b1ec0f1d4046ceb02f7603529610115e24c545fbead89f82396931cb04d
```

## Running standalone Blockbook from Docker

1. Run `make deb-blockbook-unit-e`, this will produce blockbook package in build directory
2. Run `docker build --build-arg COIN=unite --build-arg RPC_ADDRESS=127.0.0.1:8172 -t blockbook-runtime .`
3. Run `docker run -it --rm --name blockbook -p 9172:9172 -p 9272:9272 blockbook-runtime`

Make sure `RPC_ADDRESS` points to your Unit-e node's RPC port.

## Running both Blockbook and Unit-e's node in Dockers

If you want to run both Blockbook and Unit-e in separate containers, make sure both containers are in the same network.
To achieve this:
1. Create the network `docker network create blockbook-net`,
2. Get network's subnet `docker network inspect blockbook-net`
3. Configure both united and blockbook to use this subnet (replace `RPC_ADDRESS` or `rpc_url` in Blockbook and `rpc_bind` in Unit-e)
4. If cointainers are annonymous, get their ID's with `docker ps`
5. Connect both containers by calling `docker network connect blockbook-net container-id-or-name`


## Running Blockbook from repo

1. Build blockbook as in [original README](README.md)
2. Make sure you have Blockbook's config file. Example config:
```
{
    "coin_name": "Unit-e",
    "coin_shortcut": "UTE",
    "coin_label": "Unit-e",
    "rpc_url": "http://127.0.0.1:8172",
    "rpc_user": "rpc",
    "rpc_pass": "rpc",
    "rpc_timeout": 25,
    "parse": true,
    "message_queue_binding": "tcp://127.0.0.1:38330",
    "subversion": "",
    "address_format": "",
    "mempool_workers": 8,
    "mempool_sub_workers": 2,
    "block_addresses_to_keep": 300
}
```
3. Rename cert and key files from `server` directory to `blockbook.crt` and `blockbook.key` or create new and sign them
4. Run `blockbook -blockchaincfg=path_to_config.json -datadir=path-to-blockbook-db -sync -internal=:9172 -public=:9272 -certfile=path_to_cert_dir -explorer= -log_dir=log-dir -dbcache=1073741824`
