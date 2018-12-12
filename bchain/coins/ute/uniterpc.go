package ute

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/juju/errors"
)

type UniteRPC struct {
	*btc.BitcoinRPC
}

func NewUniteRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}
	u := &UniteRPC{
		BitcoinRPC: b.(*btc.BitcoinRPC),
	}
	u.RPCMarshaler = btc.JSONMarshalerV1{}
	u.ChainConfig.SupportsEstimateSmartFee = false
	return u, nil
}

// Initialize initializes UniteRPC instance.
func (u *UniteRPC) Initialize() error {
	chainName, err := u.GetChainInfoAndInitializeMempool(u)
	if err != nil {
		return err
	}

	params := GetChainParams(chainName)

	u.Parser = NewUniteParser(params, u.ChainConfig)

	// parameters for getInfo request
	if params.Net == MainnetMagic {
		u.Testnet = false
		u.Network = "mainnet"
	} else {
		u.Testnet = true
		u.Network = "testnet"
	}

	glog.Info("rpc: block chain ", params.Name)

	return nil
}

type CmdGetFinalizationState struct {
	Method string `json:"method"`
}

type ResGetEsperanzaFinalizationState struct {
	Error  *bchain.RPCError          `json:"error"`
	Result struct {
        CurrentEpoch uint64 `json:"currentEpoch"`
        CurrentDynasty uint64 `json:"currentDynasty"`
        LastJustifiedEpoch uint64 `json:"lastJustifiedEpoch"`
        LastFinalizedEpoch uint64 `json:"lastFinalizedEpoch"`
        Validators uint64 `json:"validators"`
	} `json:"result"`
}

func (u *UniteRPC) GetFinalizationState() (*ResGetEsperanzaFinalizationState, error) {
	var err error

	res := ResGetEsperanzaFinalizationState{}
	req := CmdGetFinalizationState{Method: "getfinalizationstate"}
	err = u.Call(&req, &res)

	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return &res, err
}

// GetBlock returns block with given hash.
func (u *UniteRPC) GetBlock(hash string, height uint32) (*bchain.Block, error) {
	var err error
	if hash == "" && height > 0 {
		hash, err = u.GetBlockHash(height)
		if err != nil {
			return nil, err
		}
	}

	glog.V(1).Info("rpc: getblock (verbosity=1) ", hash)

	res := btc.ResGetBlockThin{}
	req := btc.CmdGetBlock{Method: "getblock"}
	req.Params.BlockHash = hash
	req.Params.Verbosity = 1
	err = u.Call(&req, &res)

	if err != nil {
		return nil, errors.Annotatef(err, "hash %v", hash)
	}
	if res.Error != nil {
		return nil, errors.Annotatef(res.Error, "hash %v", hash)
	}

	txs := make([]bchain.Tx, 0, len(res.Result.Txids))
	for _, txid := range res.Result.Txids {
		tx, err := u.GetTransaction(txid)
		if err != nil {
			if isInvalidTx(err) {
				glog.Errorf("rpc: getblock: skipping transanction in block %s due error: %s", hash, err)
				continue
			}
			return nil, err
		}
		txs = append(txs, *tx)
	}
	block := &bchain.Block{
		BlockHeader: res.Result.BlockHeader,
		Txs:         txs,
	}
	return block, nil
}

func isInvalidTx(err error) bool {
	switch e1 := err.(type) {
	case *errors.Err:
		switch e2 := e1.Cause().(type) {
		case *bchain.RPCError:
			if e2.Code == -5 { // "No information available about transaction"
				return true
			}
		}
	}

	return false
}

// GetMempoolEntry returns mempool data for given transaction
func (u *UniteRPC) GetMempoolEntry(txid string) (*bchain.MempoolEntry, error) {
	return nil, errors.New("GetMempoolEntry: not implemented")
}
