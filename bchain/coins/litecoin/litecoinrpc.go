package litecoin

import (
	"encoding/json"

	"github.com/golang/glog"

	"github.com/dtr-org/blockbook/bchain"
	"github.com/dtr-org/blockbook/bchain/coins/btc"
)

// LitecoinRPC is an interface to JSON-RPC bitcoind service.
type LitecoinRPC struct {
	*btc.BitcoinRPC
}

// NewLitecoinRPC returns new LitecoinRPC instance.
func NewLitecoinRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	s := &LitecoinRPC{
		b.(*btc.BitcoinRPC),
	}
	s.RPCMarshaler = btc.JSONMarshalerV2{}
	s.ChainConfig.SupportsEstimateFee = false

	return s, nil
}

// Initialize initializes LitecoinRPC instance.
func (b *LitecoinRPC) Initialize() error {
	chainName, err := b.GetChainInfoAndInitializeMempool(b)
	if err != nil {
		return err
	}

	glog.Info("Chain name ", chainName)
	params := GetChainParams(chainName)

	// always create parser
	b.Parser = NewLitecoinParser(params, b.ChainConfig)

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
