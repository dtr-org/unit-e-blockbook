package ute

import (
	"errors"
	"encoding/hex"
	"github.com/jakm/btcutil"
	"github.com/jakm/btcutil/txscript"
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"

	"github.com/btcsuite/btcd/wire"
	"github.com/jakm/btcutil/chaincfg"
)

const (

	OP_TRUE byte = 0x51
	OP_IF byte = 0x63
	OP_NOTIF byte = 0x64
	OP_ELSE byte = 0x67
	OP_ENDIF byte = 0x68

	OP_DUP byte = 0x76
	OP_EQUALVERIFY byte = 0x88
	OP_HASH160 byte = 0xa9
	OP_CHECKSIG byte = 0xac
	OP_CHECKVOTESIG byte = 0xb3
	OP_SLASHABLE    byte = 0xb4

	MainnetMagic wire.BitcoinNet = 0xeeeeaec1
	TestnetMagic wire.BitcoinNet = 0xfdfcfbfa
	RegtestMagic wire.BitcoinNet = 0xfabfb5da
)

var (
	MainNetParams chaincfg.Params
	TestNetParams chaincfg.Params
	RegtestParams chaincfg.Params
)

func init() {
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic

	// Address encoding magics
	MainNetParams.AddressMagicLen = 1

	TestNetParams = chaincfg.TestNet3Params
	TestNetParams.Net = TestnetMagic

	// Address encoding magics
	TestNetParams.AddressMagicLen = 1

	RegtestParams = chaincfg.RegressionNetParams
	RegtestParams.Net = RegtestMagic
}

// UniteParser handle
type UniteParser struct {
	*btc.BitcoinParser
	baseparser *bchain.BaseParser
	Params *chaincfg.Params
}

// NewUniteParser returns new UniteParser instance
func NewUniteParser(params *chaincfg.Params, c *btc.Configuration) *UniteParser {
	p := &UniteParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
	p.OutputScriptToAddressesFunc = p.outputScriptToAddresses
	p.Params = params
	return p
}

// GetChainParams contains network parameters for the main Unite network,
// the regression test Unite network, the test Unite network and
// the simulation test Unite network, in this order
func GetChainParams(chain string) *chaincfg.Params {
	if !chaincfg.IsRegistered(&MainNetParams) {
		err := chaincfg.Register(&MainNetParams)
		if err == nil {
			err = chaincfg.Register(&TestNetParams)
		}
		if err == nil {
			err = chaincfg.Register(&RegtestParams)
		}
		if err != nil {
			panic(err)
		}
	}
	switch chain {
	case "test":
		return &TestNetParams
	case "regtest":
		return &RegtestParams
	default:
		return &MainNetParams
	}
}

// PackTx packs transaction to byte array using protobuf
func (p *UniteParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *UniteParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	tx, height, err := p.baseparser.UnpackTx(buf)
	return tx,height,err
}

// Extra-fast test for pay-to-script-hash CScripts:
func matchPayToPublicKeyHash(script []byte, ofs int) (bool) {
	matchesP2PKH := len(script) - ofs >= 25 &&
	       script[ofs+0] == OP_DUP &&
	       script[ofs+1] == OP_HASH160 &&
	       script[ofs+2] == 0x14 &&
	       script[ofs+23] == OP_EQUALVERIFY &&
	       script[ofs+24] == OP_CHECKSIG
	return matchesP2PKH
}

// Extra-fast test for pay-vote-slash script hash CScripts:
func matchPayVoteSlashScript(script []byte, ofs int) (bool) {
	return len(script) - ofs == 103 &&
	       matchVoteScript(script, 0) &&
	       script[ofs+35] == OP_IF &&
	       script[ofs+36] == OP_TRUE &&
	       script[ofs+37] == OP_ELSE &&
	       matchSlashScript(script, 38) &&
	       script[ofs+73] == OP_NOTIF &&
	       matchPayToPublicKeyHash(script, 74) &&
	       script[ofs+99] == OP_ELSE &&
	       script[ofs+100] == OP_TRUE &&
	       script[ofs+101] == OP_ENDIF &&
	       script[ofs+102] == OP_ENDIF
}

func matchVoteScript(script []byte, ofs int) (bool) {
	matchesVoteScript := len(script) - ofs >= 35 && script[ofs] == 0x21 && script[ofs+34] == OP_CHECKVOTESIG
	return matchesVoteScript
}

func matchSlashScript(script []byte, ofs int) (bool) {
	matchesSlashScript := len(script) - ofs >= 35 && script[ofs] == 0x21 && script[ofs+34] == OP_SLASHABLE
	return matchesSlashScript
}

func extractPayVoteSlashScriptAddrs(script []byte, params *chaincfg.Params) ([]string, bool, error) {
	addr, err := btcutil.NewAddressPubKeyHash(script[77:97], params)
	if err != nil {
		return nil, false, err
	}

	return []string{addr.EncodeAddress()}, true, nil
}

func isPayVoteSlashScript(script []byte) (bool) {
	return len(script) == 103 && matchPayVoteSlashScript(script, 0)
}

func IsOpReturnScript(script []byte) (bool) {
	matchesOpReturnScript := len(script) > 0 && script[0] == txscript.OP_RETURN
	return matchesOpReturnScript
}

// This function is given internal representation, so for now it's the address
func (p *UniteParser) outputScriptToAddresses(script []byte) ([]string, bool, error) {
	if len(script) == 20 {
		addr, err := btcutil.NewAddressPubKeyHash(script, p.Params)
		if err != nil {
			return []string{}, false, err
		}
		return []string{addr.EncodeAddress()}, true, nil
	} else if len(script) != 0 {
		return []string{string(script)}, false, nil
	}
	return nil, false, nil
}

func (p *UniteParser) convertAddrToStandard(address string) (string, error) {
	// For now different types of addresses will be indexed separately
	return address, nil
}

// GetAddrDescFromVout returns internal address representation (descriptor) of given transaction output
func (p *UniteParser) GetAddrDescFromVout(output *bchain.Vout) (bchain.AddressDescriptor, error) {
	if len(output.ScriptPubKey.Addresses) == 1 {
		return p.GetAddrDescFromAddress(output.ScriptPubKey.Addresses[0])
	}

	script, err := hex.DecodeString(output.ScriptPubKey.Hex)
	if err != nil {
		return nil, errors.New("Could not decode script hex")
	}

	if isPayVoteSlashScript(script) {
		addresses, _, err := extractPayVoteSlashScriptAddrs(script, p.Params)
		if err != nil {
			return nil, errors.New("Could not extract address from payvoteslash")
		}
		return p.GetAddrDescFromAddress(addresses[0])
	} else if IsOpReturnScript(script) {
		or := TryParseOPReturn(script)
		return []byte(or), nil
	}

	return nil, errors.New("Unknown address")
}

func (p *UniteParser) GetAddrDescFromAddress(address string) (bchain.AddressDescriptor, error) {
	pkh_addr, err := p.convertAddrToStandard(address)
	if err != nil {
		return nil, err
	}
	da, err := btcutil.DecodeAddress(pkh_addr, p.Params)
	if err == nil {
		return da.ScriptAddress(), nil
	}

	return []byte(address), nil
}

// Untill full segwit use addresses
func (p *UniteParser) GetScriptFromAddrDesc(addrDesc bchain.AddressDescriptor) ([]byte, error) {
	return addrDesc, nil
}

// TryParseOPReturn tries to process OP_RETURN script and return its string representation
func TryParseOPReturn(script []byte) string {
	// trying 3 variants of OP_RETURN data
	// 1) OP_RETURN
	// 1) OP_RETURN OP_PUSHDATA1 <datalen> <data>
	// 3) OP_RETURN <datalen> <data>
	if len(script) == 1 {
		return "OP_RETURN"
	}
	if len(script) > 1 && script[0] == txscript.OP_RETURN {
		var data []byte
		var l int
		if script[1] == txscript.OP_PUSHDATA1 && len(script) > 2 {
			l = int(script[2])
			data = script[3:]
			if l != len(data) {
				l = int(script[1])
				data = script[2:]
			}
		} else {
			l = int(script[1])
			data = script[2:]
		}
		if l == len(data) {
			isASCII := true
			for _, c := range data {
				if c < 32 || c > 127 {
					isASCII = false
					break
				}
			}
			var ed string
			if isASCII {
				ed = "(" + string(data) + ")"
			} else {
				ed = hex.EncodeToString(data)
			}
			return "OP_RETURN " + ed
		}
	}
	return ""
}
