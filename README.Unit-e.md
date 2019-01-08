# Unit-e Blockbook

This README is specifically for instructions on running Blockbook for the Unit-e coin.

## Running Unit-e node for Blockbook

Make sure the node is indexing transactions, has RPC server enabled, the RPC credentials are rpc:rpc and that node accepts incoming connections from Blockbook's IP.

Example flags:
```
united -txindex -server=1 -rpcport=8172 -rpcallowip=127.0.0.1/32 -rpcuser=rpc -rpcpassword=rpc
```

This uses a custom port `8172` to run the daemon which makes it easier when, for
example switching between different networks so you don't have to adapt the port
on the blockbook side.

The `-rpcallowip` option has to allow the IP address from the blockbook server.
You can use `-rpcallowip=0.0.0.0/0` to allow connections from all IP addresses
or use the specific address of the server, which you will need when it's not
running on the same host or in a container.

## Running standalone Blockbook from Docker

1. Run `make deb-blockbook-unit-e`, this will produce blockbook package in build directory
2. Run `docker build --build-arg COIN=unite -t blockbook-runtime .`
3. Run `docker run -it --rm --name blockbook -p 9172:9172 -p 9272:9272 blockbook-runtime --address 127.0.0.1:8172`

Make sure --address points to your Unit-e node's RPC port. If you run `united`
on the host you can find its IP address by running `ip a` to get the docker
network interface of the host. Or you run `docker run --entrypoint="" -it
blockbook-runtime ip route` to get the IP of the default route of the container
which is the host IP.

## Running both Blockbook and Unit-e's node in Dockers

If you want to run both Blockbook and Unit-e in separate containers, make sure both containers are in the same network.
To achieve this:
1. Create the network `docker network create blockbook-net`,
2. Get network's subnet `docker network inspect blockbook-net`
3. Configure both united and blockbook to use this subnet (you can get the IP
   addresses of the containers from `docker network inspect` when the containers
   are running):
    * Set the `united` RPC address in blockbook with the `--address` option when
      starting the container
    * Make sure the `-rpcallowip` option allows the address of the blockbook
      container to connect to `united`.
4. If containers are anonymous, get their ID's with `docker ps`
5. Connect both containers by calling `docker network connect blockbook-net
   container-id-or-name`. You can also provide the network name with the
   `--network=blockbook-net` option to the `run` command when starting the
   container.


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
3. Run `./blockbook -blockchaincfg=path_to_config.json
   -datadir=path-to-blockbook-db -sync -internal=:9172 -public=:9272
   -certfile=server/testcert -explorer= -log_dir=log-dir -dbcache=1073741824`

   This uses the test certificate from the blockbook repo. If you want to use
   your own certificate adjust the `-certfile` option accordingly.


## Creating test data

In order to see data in the Blockbook block explorer you need to connect it to
a united daemon which has data. For testing there are two ways to generate test
data:

1. Run `united` in `-regtest` mode and use the CLI to generate blocks with
   `unite-cli --rpcuser=rpc --rpcpassword=rpc generate`
2. Reuse the data created by the functional tests. For this you have to run a
   functional test with the `--nocleanup` option. This will keep the directory
   with the test data. The test output will tell you where this directory is.
   You can then run `united` with the `-datadir` option pointing to this
   directory.

   You will also need to add the esperanza options with the
   `-esperanzaconfig` option. Don't forget to quote the JSON such as
   `-esperanzaconfig='{"epochLength": 10, "minDepositSize": 1500}'`. The test
   output will tell you which exact options you need to set.

   You also will need to add the `-reindex` option so that the index required by
   the block explorer is built.
