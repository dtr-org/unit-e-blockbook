package ute

import (
	"blockbook/api"
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/golang/glog"
	"github.com/juju/errors"
)

// UniteRPC is the Unit-e coin's RPC handler
type UniteRPC struct {
	*btc.BitcoinRPC
	HTMLHandler UniteHTMLHandler
	txTypesMap  map[uint32]string
}

// UniteHTMLHandler handles Unit-e specific HTML requests
type UniteHTMLHandler struct {
	unite *UniteRPC
	t     *template.Template
}

// UniteTemplateData holds the Unit-e specific template data
// Implements interface{} passed to template as data
type UniteTemplateData struct {
	FinalizationConfig *FinalizationConfig
	FinalizationState  *FinalizationState
}

// FinalizationConfig holds the JSON response to getfinalizationconfig message
type FinalizationConfig struct {
	EpochLength               int    `json:"epochLength"`
	MinDepositSize            int64  `json:"minDepositSize"`
	DynastyLogoutDelay        int    `json:"dynastyLogoutDelay"`
	WithdrawalEpochDelay      int    `json:"withdrawalEpochDelay"`
	BountyFractionDenominator int    `json:"bountyFractionDenominator"`
	SlashFractionMultiplier   int    `json:"slashFractionMultiplier"`
	BaseInterestFactor        string `json:"baseInterestFactor"`
	BasePenaltyFactor         string `json:"basePenaltyFactor"`
}

// FinalizationState holds the JSON response to getfinalizationstate message
type FinalizationState struct {
	CurrentEpoch       int `json:"currentEpoch"`
	CurrentDynasty     int `json:"currentDynasty"`
	LastFinalizedEpoch int `json:"lastFinalizedEpoch"`
	LastJustifiedEpoch int `json:"lastJustifiedEpoch"`
	Validators         int `json:"validators"`
}

func (u *UniteRPC) txTypeToString(t uint32) string {
	v, _ := u.txTypesMap[t]
	return v
}

func (u *UniteRPC) extractVoteFromTx(tx *api.Tx) *Vote {
	return ExtractVoteFromSignature(tx.Vin[0].Hex)
}

func (u *UniteRPC) extractSlashFromTx(tx *api.Tx) *Slash {
	return ExtractSlashFromSignature(tx.Vin[0].Hex)
}

// NewUniteHTMLHandler creates the Unit-e's HTML handler and populates it's templates
func NewUniteHTMLHandler(u *UniteRPC) UniteHTMLHandler {
	h := UniteHTMLHandler{
		unite: u,
	}

	h.t = template.Must(template.New("ute").ParseFiles("./static/templates/coins/ute.html", "./static/templates/base.html"))
	return h
}

// NewUniteRPC returns UniteRPC from configuration
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
	u.HTMLHandler = NewUniteHTMLHandler(u)

	u.txTypesMap = map[uint32]string{
		0: "Standard",
		1: "Coinbase",
		2: "Deposit",
		3: "Vote",
		4: "Logout",
		5: "Slash",
		6: "Withdraw",
		7: "Admin",
	}

	return u, nil
}

type cmdGetFinalizationConfig struct {
	Method string `json:"method"`
}

type resGetFinalizationConfig struct {
	Error  *bchain.RPCError    `json:"error"`
	Result *FinalizationConfig `json:"result"`
}

// Initialize initializes UniteRPC instance.
func (u *UniteRPC) Initialize() error {
	ci, err := u.GetChainInfo()
	if err != nil {
		return err
	}
	chainName := ci.Chain

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

type cmdGetFinalizationState struct {
	Method string `json:"method"`
}

type resGetFinalizationState struct {
	Error  *bchain.RPCError   `json:"error"`
	Result *FinalizationState `json:"result"`
}

func (u *UniteRPC) getFinalizationConfig() (*FinalizationConfig, error) {
	var err error

	res := resGetFinalizationConfig{}
	req := cmdGetFinalizationConfig{Method: "getfinalizationconfig"}
	err = u.Call(&req, &res)

	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return res.Result, err
}

func (u *UniteRPC) getFinalizationState() (*FinalizationState, error) {
	var err error

	res := resGetFinalizationState{}
	req := cmdGetFinalizationState{Method: "getfinalizationstate"}
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
	return nil, errors.New("unit-e rpc: GetMempoolEntry: not implemented")
}

// GetCoinHTMLHandler returns the HTML handler
func (u *UniteRPC) GetCoinHTMLHandler() bchain.CoinHTMLHandler {
	return &u.HTMLHandler
}

// GetExtraNavItems returns the extra navigation bar items
func (h *UniteHTMLHandler) GetExtraNavItems() map[string]string {
	return map[string]string{
		"Finalization": "/coin/",
	}
}

// GetExtraFuncMap returns the extra functions to be registered in templates
func (h *UniteHTMLHandler) GetExtraFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatTxType":       h.unite.txTypeToString,
		"extractVoteFromTx":  h.unite.extractVoteFromTx,
		"extractSlashFromTx": h.unite.extractSlashFromTx,
	}
}

// HandleCoinRequest returns template for given path and it's data
func (h *UniteHTMLHandler) HandleCoinRequest(w http.ResponseWriter, r *http.Request) (*template.Template, interface{}, error) {
	state, err := h.unite.getFinalizationState()
	if err != nil {
		return nil, nil, err
	}
	config, err := h.unite.getFinalizationConfig()
	if err != nil {
		return nil, nil, err
	}

	data := UniteTemplateData{
		FinalizationConfig: config,
		FinalizationState:  state,
	}

	return h.t, data, err
}
