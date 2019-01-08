# Unit-e Blockbook

This README is specifically for instructions on running Blockbook for the Unit-e coin.

## Running Unit-e node for Blockbook

Make sure the node is indexing transactions, has RPC server enabled, the RPC credentials are rpc:rpc and that node accepts incoming connections from Blockbook's IP.

Example flags:
```
united -txindex -server=1 -rpcbind=127.0.0.1:8172 -rpcallowip=127.0.0.1/32 -rpcuser=rpc -rpcpassword=rpc
```

## Running standalone Blockbook from Docker

1. Run `make deb-blockbook-unit-e`, this will produce blockbook package in build directory
2. Run `docker build --build-arg COIN=unite -t blockbook-runtime .`
3. Run `docker run -it --rm --name blockbook -p 9172:9172 -p 9272:9272 blockbook-runtime --address 127.0.0.1:8172`

Make sure --address points to your Unit-e node's RPC port.

## Running both Blockbook and Unit-e's node in Dockers

If you want to run both Blockbook and Unit-e in separate containers, make sure both containers are in the same network.
To achieve this:
1. Create the network `docker network create blockbook-net`,
2. Get network's subnet `docker network inspect blockbook-net`
3. Configure both united and blockbook to use this subnet (make sure to provide correct rpc address in Blockbook and `rpc_bind` in Unit-e)
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
