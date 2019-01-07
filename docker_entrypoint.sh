#!/bin/bash

USAGE="Usage: docker_entrypoint.sh [--address=HOST:PORT] [-- command]"
# Default address is localhost (127.0.0.1)
# Default command is seen below
DEFAULT_CMD="./bin/blockbook -blockchaincfg=config/blockchaincfg.json -datadir=db -sync -internal=:9172 -public=:9272 -certfile=cert/blockbook -explorer= -log_dir=logs -dbcache=1073741824"
DEFAULT_ADDRESS="127.0.0.1:8172"

CMD=${DEFAULT_CMD}
ADDRESS=${DEFAULT_ADDRESS}


while [[ $# -gt 0 ]]
do
key="$1"
case $key in
    --address)
    ADDRESS="$2"
    shift
    shift
    ;;
    --address=*)
    ADDRESS="${key#*=}"
    shift
    ;;
    -h|--help)
    echo ${USAGE}
    exit 0
    ;;
    --)
    shift
    CMD="$@"
    break
    ;;
    *)
    echo "Unknown argument, $1 exiting"
    echo ${USAGE}
    exit 1
    ;;
esac
done

echo $ADDRESS

sed -i 's/"rpc_url": "[^\"].*"/"rpc_url": "http:\/\/'${ADDRESS}'"/' config/blockchaincfg.json
eval ${CMD}
