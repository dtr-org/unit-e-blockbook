package vertcoin

import (
	"encoding/json"

	"github.com/golang/glog"

	"github.com/dtr-org/blockbook/bchain"
	"github.com/dtr-org/blockbook/bchain/coins/btc"
)

// VertcoinRPC is an interface to JSON-RPC bitcoind service.
type VertcoinRPC struct {
	*btc.BitcoinRPC
}

// NewVertcoinRPC returns new VertcoinRPC instance.
func NewVertcoinRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	s := &VertcoinRPC{
		b.(*btc.BitcoinRPC),
	}
	s.RPCMarshaler = btc.JSONMarshalerV2{}
	s.ChainConfig.SupportsEstimateFee = false

	return s, nil
}

// Initialize initializes VertcoinRPC instance.
func (b *VertcoinRPC) Initialize() error {
	chainName, err := b.GetChainInfoAndInitializeMempool(b)
	if err != nil {
		return err
	}

	glog.Info("Chain name ", chainName)
	params := GetChainParams(chainName)

	// always create parser
	b.Parser = NewVertcoinParser(params, b.ChainConfig)

	// parameters for getInfo request
	if params.Net == MainnetMagic {
		b.Testnet = false
		b.Network = "livenet"
	} else {
		b.Testnet = true
		b.Network = "testnet"
	}

	glog.Info("rpc: block chain ", params.Name)

	return nil
}
