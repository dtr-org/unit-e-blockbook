package ute

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/jakm/btcutil"

	"github.com/btcsuite/btcd/wire"
	"github.com/jakm/btcutil/chaincfg"
)

// Script opcodes
const (
	OpPushdata1 byte = 0x4c
	OpPushdata2 byte = 0x4d
	OpPushdata4 byte = 0x4e
	OpTrue      byte = 0x51
	OpIf        byte = 0x63
	OpNotif     byte = 0x64
	OpElse      byte = 0x67
	OpEndif     byte = 0x68
	OpReturn    byte = 0x6a

	OpDup          byte = 0x76
	OpEqualverify  byte = 0x88
	OpHash160      byte = 0xa9
	OpChecksig     byte = 0xac
	OpCheckvotesig byte = 0xb3
	OpSlashable    byte = 0xb4
)

// Network bytes
const (
	MainnetMagic wire.BitcoinNet = 0xeeeeaec1
	TestnetMagic wire.BitcoinNet = 0xfdfcfbfa
	RegtestMagic wire.BitcoinNet = 0xfabfb5da
)

// Blockchain parameters
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
	MainNetParams.Bech32HRPSegwit = "ue"

	TestNetParams = chaincfg.TestNet3Params
	TestNetParams.Net = TestnetMagic
	TestNetParams.Bech32HRPSegwit = "tue"

	// Address encoding magics
	TestNetParams.AddressMagicLen = 1

	RegtestParams = chaincfg.RegressionNetParams
	RegtestParams.Net = RegtestMagic
	RegtestParams.Bech32HRPSegwit = "uert"
}

// UniteParser handle
type UniteParser struct {
	*btc.BitcoinParser
	baseparser *bchain.BaseParser
	Params     *chaincfg.Params
}

// Vote data
type Vote struct {
	ValidatorAddress string
	TargetHash       string
	SourceEpoch      uint32
	TargetEpoch      uint32
}

// Slash data
type Slash struct {
	Vote1 *Vote
	Vote2 *Vote
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

// GetChainParams contains network parameters for the main Unit-e network,
// the regression test Unit-e network, the test Unite network and
// the simulation test Unit-e network, in this order
func GetChainParams(chain string) *chaincfg.Params {
	if !chaincfg.IsRegistered(&MainNetParams) {
		if err := chaincfg.Register(&MainNetParams); err != nil {
			panic(err)
		}
		if err := chaincfg.Register(&TestNetParams); err != nil {
			panic(err)
		}
		if err := chaincfg.Register(&RegtestParams); err != nil {
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

// GetVaruint returns decoded int from byte array
func GetVaruint(arr []byte) uint64 {
	l := len(arr)
	if l == 0 {
		// treating as empty
		return 0
	}

	var r uint64
	var mod uint64 = 1 << uint64(8*(l-1))
	for _, b := range arr {
		r |= uint64(b) * mod
		mod = mod >> 8
	}

	return r
}

// GetOp parses script at given offset and returns new offset and the bytes read
func GetOp(script []byte, ofs uint64) (uint64, []byte, error) {
	scriptLen := uint64(len(script))
	if scriptLen <= ofs {
		return 0, []byte{}, errors.New("invalid script, offset outside bounds")
	}

	opcode := script[ofs]
	if opcode <= OpPushdata4 {
		var nSize uint64
		var dataStart uint64 = ofs
		if opcode < OpPushdata1 {
			dataStart = ofs + 1
			nSize = uint64(opcode)
		} else {
			if opcode == OpPushdata1 {
				dataStart = ofs + 2
			} else if opcode == OpPushdata2 {
				dataStart = ofs + 3
			} else if opcode == OpPushdata4 {
				dataStart = ofs + 5
			}
			if scriptLen < dataStart {
				return 0, []byte{}, errors.New("invalid script, not enough elements after OP_PUSHDATA")
			}
			nSize = GetVaruint(script[ofs+1 : dataStart])
		}

		if scriptLen < dataStart+nSize {
			return 0, []byte{}, errors.New("invalid script, not enough elements")
		}

		return dataStart + nSize, script[dataStart : dataStart+nSize], nil
	}

	return ofs + 1, []byte{opcode}, nil
}

// DecodeVote returns decoded vote from script
func DecodeVote(voteScript []byte) *Vote {
	// read voteSig
	ofs, _, err := GetOp(voteScript, 0)
	if err != nil {
		return nil
	}

	// read validatorAddress
	ofs, validatorAddress, err := GetOp(voteScript, ofs)
	if err != nil {
		return nil
	}

	// read targetHash
	ofs, targetHash, err := GetOp(voteScript, ofs)
	if err != nil {
		return nil
	}

	// read sourceEpochVec
	ofs, sourceEpochV, err := GetOp(voteScript, ofs)
	if err != nil {
		return nil
	}
	sourceEpoch := GetVaruint(sourceEpochV)

	// read targetEpochVec
	ofs, targetEpochV, err := GetOp(voteScript, ofs)
	if err != nil {
		return nil
	}
	targetEpoch := GetVaruint(targetEpochV)

	return &Vote{ValidatorAddress: hex.EncodeToString(validatorAddress), TargetHash: hex.EncodeToString(targetHash), SourceEpoch: uint32(sourceEpoch), TargetEpoch: uint32(targetEpoch)}
}

// ExtractVoteFromSignature reads and decodes vote from signature
func ExtractVoteFromSignature(sigHex string) *Vote {
	script, err := hex.DecodeString(sigHex)
	if err != nil {
		return nil
	}

	// read txSig (ignored)
	ofs, _, err := GetOp(script, 0)
	if err != nil {
		return nil
	}

	// read vote
	ofs, vote, err := GetOp(script, ofs)
	if err != nil {
		return nil
	}

	return DecodeVote(vote)
}

// ExtractSlashFromSignature reads and decodes votes from signature of a slash
func ExtractSlashFromSignature(sigHex string) *Slash {
	script, err := hex.DecodeString(sigHex)
	if err != nil {
		return nil
	}

	// read txSig (ignored)
	ofs, _, err := GetOp(script, 0)
	if err != nil {
		return nil
	}

	// read first vote
	ofs, vote1, err := GetOp(script, ofs)
	if err != nil {
		return nil
	}

	// read second vote
	ofs, vote2, err := GetOp(script, ofs)
	if err != nil {
		return nil
	}

	return &Slash{Vote1: DecodeVote(vote1), Vote2: DecodeVote(vote2)}
}

// PackTx packs transaction to byte array using protobuf
func (p *UniteParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *UniteParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	tx, height, err := p.baseparser.UnpackTx(buf)
	return tx, height, err
}

// Extra-fast test for pay-to-script-hash CScripts:
func matchPayToPublicKeyHash(script []byte, ofs int) bool {
	matchesP2PKH := len(script)-ofs >= 25 &&
		script[ofs+0] == OpDup &&
		script[ofs+1] == OpHash160 &&
		script[ofs+2] == 0x14 &&
		script[ofs+23] == OpEqualverify &&
		script[ofs+24] == OpChecksig
	return matchesP2PKH
}

// Extra-fast test for pay-vote-slash script hash CScripts:
func matchPayVoteSlashScript(script []byte, ofs int) bool {
	return len(script)-ofs == 103 &&
		matchVoteScript(script, 0) &&
		script[ofs+35] == OpIf &&
		script[ofs+36] == OpTrue &&
		script[ofs+37] == OpElse &&
		matchSlashScript(script, 38) &&
		script[ofs+73] == OpNotif &&
		matchPayToPublicKeyHash(script, 74) &&
		script[ofs+99] == OpElse &&
		script[ofs+100] == OpTrue &&
		script[ofs+101] == OpEndif &&
		script[ofs+102] == OpEndif
}

func matchVoteScript(script []byte, ofs int) bool {
	matchesVoteScript := len(script)-ofs >= 35 && script[ofs] == 0x21 && script[ofs+34] == OpCheckvotesig
	return matchesVoteScript
}

func matchSlashScript(script []byte, ofs int) bool {
	matchesSlashScript := len(script)-ofs >= 35 && script[ofs] == 0x21 && script[ofs+34] == OpSlashable
	return matchesSlashScript
}

func extractPayVoteSlashScriptAddrs(script []byte, params *chaincfg.Params) ([]string, bool, error) {
	addr, err := btcutil.NewAddressPubKeyHash(script[77:97], params)
	if err != nil {
		return nil, false, err
	}

	return []string{addr.EncodeAddress()}, true, nil
}

func isPayVoteSlashScript(script []byte) bool {
	return len(script) == 103 && matchPayVoteSlashScript(script, 0)
}

// IsOpReturnScript returns whether script is OP_RETURN-type script
func IsOpReturnScript(script []byte) bool {
	matchesOpReturnScript := len(script) > 0 && script[0] == OpReturn
	return matchesOpReturnScript
}

// This function is given internal representation, so for now it's the address
func (p *UniteParser) outputScriptToAddresses(script []byte) ([]string, bool, error) {
	if len(script) == 20 {
		addr, err := btcutil.NewAddressPubKeyHash(script, p.Params)
		if err != nil {
			return nil, false, err
		}
		return []string{addr.EncodeAddress()}, true, nil
	} else if strings.HasPrefix(string(script), p.Params.Bech32HRPSegwit) {
		return []string{string(script)}, true, nil
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
		return nil, errors.New("unit-e parser: could not decode script hex")
	}

	if isPayVoteSlashScript(script) {
		addresses, _, err := extractPayVoteSlashScriptAddrs(script, p.Params)
		if err != nil {
			return nil, errors.New("unit-e parser: could not extract address from payvoteslash")
		}
		return p.GetAddrDescFromAddress(addresses[0])
	} else if IsOpReturnScript(script) {
		or := TryParseOPReturn(script)
		return []byte(or), nil
	}

	return nil, errors.New("unit-e parser: unknown address")
}

// GetAddrDescFromAddress Returns address encoded to bytes
func (p *UniteParser) GetAddrDescFromAddress(address string) (bchain.AddressDescriptor, error) {
	pkhAddr, err := p.convertAddrToStandard(address)
	if err != nil {
		return nil, err
	}
	da, err := btcutil.DecodeAddress(pkhAddr, p.Params)
	if err == nil {
		return da.ScriptAddress(), nil
	}

	return []byte(address), nil
}

// GetScriptFromAddrDesc returns unchanged address as it's the internal type
// Untill full segwit use addresses
func (p *UniteParser) GetScriptFromAddrDesc(addrDesc bchain.AddressDescriptor) ([]byte, error) {
	return addrDesc, nil
}

// TryParseOPReturn tries to process OpReturn script and return its string representation
func TryParseOPReturn(script []byte) string {
	// trying 3 variants of OP_RETURN data
	// 1) OP_RETURN
	// 1) OP_RETURN OP_PUSHDATA1 <datalen> <data>
	// 3) OP_RETURN <datalen> <data>
	if len(script) == 0 || script[0] != OpReturn {
		return ""
	}

	if len(script) == 1 {
		return "OP_RETURN"
	}

	var data []byte
	var l int
	if script[1] == OpPushdata1 && len(script) > 2 {
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

	return ""
}
