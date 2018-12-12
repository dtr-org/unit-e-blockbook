package ute

import (
	"net/http"
	"html/template"
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/juju/errors"
)

type UniteRPC struct {
	*btc.BitcoinRPC
	HtmlHandler UniteHtmlHandler
}

type UniteHtmlHandler struct {
	unite *UniteRPC
	t *template.Template
}

type UniteTemplateData struct {
	EsperanzaState *EsperanzaState
}

type EsperanzaState struct {
	CurrentEpoch        int `json:"currentEpoch"`
	CurrentDynasty      int `json:"currentDynasty"`
	LastFinalizedEpoch  int `json:"lastFinalizedEpoch"`
	LastJustifiedEpoch  int `json:"lastJustifiedEpoch"`
	Validators          int `json:"validators"`
}

func NewUniteHtmlHandler(u *UniteRPC) (UniteHtmlHandler) {
	h := UniteHtmlHandler{
		unite: u,
	}
	h.t = template.Must(template.New("ute").ParseFiles("./static/templates/coins/ute.html", "./static/templates/base.html"))
	return h
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
	u.HtmlHandler = NewUniteHtmlHandler(u)

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
	Error  *bchain.RPCError `json:"error"`
	Result *EsperanzaState  `json:"result"`
}

func (u *UniteRPC) GetFinalizationState() (*EsperanzaState, error) {
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

	return res.Result, err
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

func (u *UniteRPC) GetCoinHtmlHandler() bchain.CoinHtmlHandler {
	return &u.HtmlHandler
}

func (h *UniteHtmlHandler) GetExtraNavItems() map[string]string {
	return map[string]string {
		"Esperanza": "/coin/",
	}
}

func (h *UniteHtmlHandler) HandleCoinRequest(w http.ResponseWriter, r *http.Request) (*template.Template, interface{}, error) {
	state, err := h.unite.GetFinalizationState()
	if err != nil {
		return nil, nil, err
	}

	data := UniteTemplateData{
		EsperanzaState: state,
	}

	return h.t, data, err
}

