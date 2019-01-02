// +build unittest

package ute

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/jakm/btcutil/chaincfg"
)

var (
	testTx1, testTx2 bchain.Tx

	testTxPacked1 = "0200030001d86cd422c107fe6c0d4c3ac462ec33407f7200dcc00ffc253ed47ffe33bef30800000000cc483045022100bacc2a56b0f64c2049867eccb84c82274531f9e83803dac9c175c15fcc4a21fb022045ff57194e4ed1e314ac78ce78ac8d6829f593d3adf129d7ee591f5b4bed987e014c8146304402200b6459709d33acf8185c5f934df3cd3cab138fde917491d0acbb5820ba92d887022025159929618327f19c966655005f22069e40d7819e6b35dd72d39661e880f71314a57e1e892f3031232356ecfddb01025579993571208ced6f6944123a69ff006ad75d72c570bb9c9b247a11dce426359c0dbadaf125010d010effffffff016844b2ec22000000672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb36351672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb46476a914a57e1e892f3031232356ecfddb0102557999357188ac6751686800000000"
	testTxPacked2 = "0a2008f3be33fe7fd43e25fc0fc0dc00727f4033ec62c43a4c0d6cfe07c122d46cd810001aad0202000200019752a960f975b3d1232f8a6f0c087bf4f3b7a7f264cc094c214f1448e433d888000000006a473044022071c7f9d3f8e260623dc0b50323779faea1bf2a660e53fe5744af32e6c486235c02206bbe9cd425d191f7ea33a1ce2df0a8799ca7794a6ee67610754a6083c9be4b24012103ce22dff385b9c2f95bd972ca1a0e3feb92a3507797f121bd35eb43fe96a816effeffffff026844b2ec22000000672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb36351672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb46476a914a57e1e892f3031232356ecfddb0102557999357188ac6751686800b4f2e7c500000017a91417fcd0eb8dea5d92da76e1494a4f3aa07f4557bc877800000020e687bae005287830f9ea113a98010a00122088d833e448144f214c09cc64f2a7b7f3f47b080c6f8a2f23d1b375f960a952971800226a473044022071c7f9d3f8e260623dc0b50323779faea1bf2a660e53fe5744af32e6c486235c02206bbe9cd425d191f7ea33a1ce2df0a8799ca7794a6ee67610754a6083c9be4b24012103ce22dff385b9c2f95bd972ca1a0e3feb92a3507797f121bd35eb43fe96a816ef28feffffff0f42a4010a05037e11d3a410001a672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb36351672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb46476a914a57e1e892f3031232356ecfddb0102557999357188ac67516868220c706179766f7465736c6173682a226d76627a73644c54756b364370444e554c56366547705152675a705475424131536a42530a0513ca65120010011a17a91417fcd0eb8dea5d92da76e1494a4f3aa07f4557bc87220a736372697074686173682a23324d75533450666254786b4a566763483174784c583448357466694179684e63487766"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:       "0200030001d86cd422c107fe6c0d4c3ac462ec33407f7200dcc00ffc253ed47ffe33bef30800000000cc483045022100bacc2a56b0f64c2049867eccb84c82274531f9e83803dac9c175c15fcc4a21fb022045ff57194e4ed1e314ac78ce78ac8d6829f593d3adf129d7ee591f5b4bed987e014c8146304402200b6459709d33acf8185c5f934df3cd3cab138fde917491d0acbb5820ba92d887022025159929618327f19c966655005f22069e40d7819e6b35dd72d39661e880f71314a57e1e892f3031232356ecfddb01025579993571208ced6f6944123a69ff006ad75d72c570bb9c9b247a11dce426359c0dbadaf125010d010effffffff016844b2ec22000000672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb36351672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb46476a914a57e1e892f3031232356ecfddb0102557999357188ac6751686800000000",
		Blocktime: 1544455142,
		Time:      1544455142,
		Txid:      "1a2b48881af23daca0759fa9c9a07b101f4ca457266f2aa8fdffd9e14ed926e4",
		LockTime:  0,
		Vin: []bchain.Vin{
			{
				Coinbase: "0178202a0dbc1bfc55107b64c8a1be4e66e8df23aad8c3a1bdd55aacf466458b29cae00101",
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(500010376),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "21024ad5f27f1d5360b2bf2b30a3ffa7ca4db389912d4b8487bc2e4dce2657d63105ac",
					Addresses: []string{
						"mo1XAXjBVpofyXns4XXpnh3jxpm3cz47Z4",
					},
					Type: "pubkey",
				},
			},
			{
				ValueSat: *big.NewInt(0),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex:  "6a24aa21a9ed32d4e8028e9a5a650f55498f8fd3b11689adfcb359be0b1281d2c097cfbe114d",
					Type: "nulldata",
				},
			},
		},
	}

	testTx2 = bchain.Tx{
		Hex:       "02000200019752a960f975b3d1232f8a6f0c087bf4f3b7a7f264cc094c214f1448e433d888000000006a473044022071c7f9d3f8e260623dc0b50323779faea1bf2a660e53fe5744af32e6c486235c02206bbe9cd425d191f7ea33a1ce2df0a8799ca7794a6ee67610754a6083c9be4b24012103ce22dff385b9c2f95bd972ca1a0e3feb92a3507797f121bd35eb43fe96a816effeffffff026844b2ec22000000672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb36351672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb46476a914a57e1e892f3031232356ecfddb0102557999357188ac6751686800b4f2e7c500000017a91417fcd0eb8dea5d92da76e1494a4f3aa07f4557bc8778000000",
		Blocktime: 1544455142,
		Time:      1544455142,
		Txid:      "08f3be33fe7fd43e25fc0fc0dc00727f4033ec62c43a4c0d6cfe07c122d46cd8",
		LockTime:  120,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "473044022071c7f9d3f8e260623dc0b50323779faea1bf2a660e53fe5744af32e6c486235c02206bbe9cd425d191f7ea33a1ce2df0a8799ca7794a6ee67610754a6083c9be4b24012103ce22dff385b9c2f95bd972ca1a0e3feb92a3507797f121bd35eb43fe96a816ef",
				},
				Txid:     "88d833e448144f214c09cc64f2a7b7f3f47b080c6f8a2f23d1b375f960a95297",
				Vout:     0,
				Sequence: 4294967294,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(14999999396),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "2103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb36351672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb46476a914a57e1e892f3031232356ecfddb0102557999357188ac67516868",
					Addresses: []string{
						"mvbzsdLTuk6CpDNULV6eGpQRgZpTuBA1Sj",
					},
					Type: "payvoteslash",
				},
			},
			{
				ValueSat: *big.NewInt(85000000000),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "a91417fcd0eb8dea5d92da76e1494a4f3aa07f4557bc87",
					Addresses: []string{
						"2MuS4PfbTxkJVgcH1txLX4H5tfiAyhNcHwf",
					},
					Type: "scripthash",
				},
			},
		},
	}
}

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func TestParseOpReturn(t *testing.T) {
	parser := NewUniteParser(btc.GetChainParams("regtest"), &btc.Configuration{})

	// OP_RETURN
	opReturnHexString := "6a"
	opReturnHexBytes, err := hex.DecodeString(opReturnHexString)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
		return
	}
	if !IsOpReturnScript(opReturnHexBytes) {
		t.Errorf("Script should be OpReturn script")
		return
	}

	opRetAddr := TryParseOPReturn(opReturnHexBytes)
	if opRetAddr != "OP_RETURN" {
		t.Errorf("OpRetAddr is incorrect %s", opRetAddr)
		return
	}

	// OP_RETURN datalen data
	opReturnHexString = "6a24aa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2"
	opReturnHexBytes, err = hex.DecodeString(opReturnHexString)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
		return
	}
	if !IsOpReturnScript(opReturnHexBytes) {
		t.Errorf("Script should be OpReturn script")
		return
	}

	opRetAddr = TryParseOPReturn(opReturnHexBytes)
	if opRetAddr != "OP_RETURN aa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2" {
		t.Errorf("OpRetAddr is incorrect %s", opRetAddr)
		return
	}

	vout := bchain.Vout{
		ValueSat: *big.NewInt(0),
		N:        1,
		ScriptPubKey: bchain.ScriptPubKey{
			Hex:  "6a24aa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2",
			Type: "nulldata",
		},
	}

	desc, err := parser.GetAddrDescFromVout(&vout)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
		return
	}
	if string(desc) != "OP_RETURN aa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2" {
		t.Errorf("Desc from Vout is incorrect: %s", desc)
	}

	desc, err = parser.GetAddrDescFromAddress(string(desc))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
		return
	}
	if string(desc) != "OP_RETURN aa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2" {
		t.Errorf("Desc from Address is incorrect: %s", desc)
	}

	addrs, searchable, err := parser.GetAddressesFromAddrDesc(desc)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
		return
	}
	if searchable {
		t.Errorf("OP_RETURN should not be searchable")
		return
	}
	if len(addrs) != 1 || addrs[0] != "OP_RETURN aa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2" {
		t.Errorf("Address from Desc is incorrect: %s", addrs[0])
	}

	script, err := parser.GetScriptFromAddrDesc(desc)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
		return
	}
	if string(script) != string(desc) {
		t.Errorf("Script does not match desc")
	}
}

func TestGetAddrDesc(t *testing.T) {
	type args struct {
		tx           bchain.Tx
		parser       *UniteParser
		wantsReverse bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ute-1",
			args: args{
				tx:           testTx1,
				parser:       NewUniteParser(btc.GetChainParams("regtest"), &btc.Configuration{}),
				wantsReverse: true,
			},
		},
		{
			name: "ute-2",
			args: args{
				tx:           testTx2,
				parser:       NewUniteParser(GetChainParams("regtest"), &btc.Configuration{}),
				wantsReverse: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for n, vout := range tt.args.tx.Vout {
				if len(vout.ScriptPubKey.Addresses) == 1 {
					got1, err := tt.args.parser.GetAddrDescFromVout(&vout)
					if err != nil {
						t.Errorf("GetAddrDescFromAddress() error = %v, vout = %d", err, n)
						return
					}
					got2, err := tt.args.parser.GetAddrDescFromAddress(vout.ScriptPubKey.Addresses[0])
					if err != nil {
						t.Errorf("GetAddrDescFromAddress() error = %v, vout = %d", err, n)
						return
					}
					got3, _, err := tt.args.parser.GetAddressesFromAddrDesc(got1)
					if err != nil {
						t.Errorf("GetAddressesFromAddrDesc() error = %v, vout = %d", err, n)
						return
					}
					if !bytes.Equal(got1, got2) {
						t.Errorf("Address descriptors mismatch: got1 = %v, got2 = %v", got1, got2)
					}
					if got3[0] != vout.ScriptPubKey.Addresses[0] && tt.args.wantsReverse {
						t.Errorf("Address reverse lookup mismatch: got3 = %v, address = %v", got3, vout.ScriptPubKey.Addresses[0])
					}
				}
			}
		})
	}
}

func TestPackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *UniteParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ute-1",
			args: args{
				tx:        testTx1,
				height:    292272,
				blockTime: 1544455142,
				parser:    NewUniteParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "ute-2",
			args: args{
				tx:        testTx2,
				height:    292217,
				blockTime: 1544455142,
				parser:    NewUniteParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestUnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *UniteParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "ute-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewUniteParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    nil,
			want1:   0,
			wantErr: true,
		},
		{
			name: "ute-2",
			args: args{
				packedTx: testTxPacked2,
				parser:   NewUniteParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx2,
			want1:   292217,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetOp1(t *testing.T) {
	tests := []struct {
		hex         string
		off         uint64
		noff        uint64
		expected    string
		expectedErr error
	}{
		{
			hex:         "01",
			noff:        0,
			off:         2,
			expected:    "00",
			expectedErr: errors.New("invalid script, offset outside bounds"),
		},
		{
			hex:         "01",
			noff:        0,
			expected:    "00",
			expectedErr: errors.New("invalid script, not enough elements"),
		},
		{
			hex:         "ff",
			noff:        1,
			expected:    "ff",
			expectedErr: nil,
		},
		{
			hex:         "4dff",
			noff:        0,
			expected:    "00",
			expectedErr: errors.New("invalid script, not enough elements after OP_PUSHDATA"),
		},
		{
			hex:         "4dffff",
			noff:        0,
			expected:    "00",
			expectedErr: errors.New("invalid script, not enough elements"),
		},
		{
			hex:         "4effffffff00",
			noff:        0,
			expected:    "00",
			expectedErr: errors.New("invalid script, not enough elements"),
		},
		{
			hex:         "01ff",
			noff:        2,
			expected:    "ff",
			expectedErr: nil,
		},
		{
			hex:         "4e0000000105",
			noff:        6,
			expected:    "05",
			expectedErr: nil,
		},
	}

	for i, test := range tests {
		script, err := hex.DecodeString(test.hex)
		if err != nil {
			t.Errorf("test case %v: unexpected parsing error %v", i, err)
			return
		}

		noff, op, err := GetOp(script, test.off)
		if test.expectedErr != nil {
			if err == nil {
				t.Errorf("test case %v: did not get error, want %v", i, test.expectedErr)
			}
			if err.Error() != test.expectedErr.Error() {
				t.Errorf("test case %v: unexpected error %v, want %v", i, err, test.expectedErr)
				return
			}
			continue
		}

		if err != nil {
			t.Errorf("test case %v: unexpected error %v", i, err)
			return
		}

		if noff != test.noff {
			t.Errorf("test case %v: unexpected new offset %v, want %v", i, noff, test.noff)
			return
		}

		got := hex.EncodeToString(op)
		if got != test.expected {
			t.Errorf("test case %v: unexpected value %v, want %v", i, got, test.expected)
			return
		}
	}
}

func TestGetOp2(t *testing.T) {
	hexString := "473044022041007ad95eaf56b4d5d2629cf96d2f968ad4488102c93ee7418bfa4c3c7a3c3d02207f0c36db3a9bb30995e88ed321cd144f612bba9ab84d0fdd28f0eb7a17055ffb014c82473045022100e77bd5fe006cd973f8d231ffeb26b80eddbeab93ed9f382f55276ddf98c882840220027a8ad1efcc41606dd77b48f5658bf8270fb51ada97b2b36d6a60649280c6c114a57e1e892f3031232356ecfddb0102557999357120c057ed3a9a722f0ebab32aa4231df45deeec409f45206f0942d9103d9504d556010e010f"
	script, err := hex.DecodeString(hexString)
	if err != nil {
		t.Errorf("unexpected error = %v", err)
		return
	}

	noff, op, err := GetOp(script, 0)
	if err != nil {
		t.Errorf("unexpected error = %v", err)
		return
	}

	if noff != 72 {
		t.Errorf("unexpected new offset = %v", noff)
		return
	}

	got := hex.EncodeToString(op)
	if got != "3044022041007ad95eaf56b4d5d2629cf96d2f968ad4488102c93ee7418bfa4c3c7a3c3d02207f0c36db3a9bb30995e88ed321cd144f612bba9ab84d0fdd28f0eb7a17055ffb01" {
		t.Errorf("unexpected value = %v", got)
		return
	}

	noff, op, err = GetOp(script, noff)
	if err != nil {
		t.Errorf("unexpected error = %v", err)
		return
	}

	if noff != 204 {
		t.Errorf("unexpected new offset = %v", noff)
		return
	}

	got = hex.EncodeToString(op)
	if got != "473045022100e77bd5fe006cd973f8d231ffeb26b80eddbeab93ed9f382f55276ddf98c882840220027a8ad1efcc41606dd77b48f5658bf8270fb51ada97b2b36d6a60649280c6c114a57e1e892f3031232356ecfddb0102557999357120c057ed3a9a722f0ebab32aa4231df45deeec409f45206f0942d9103d9504d556010e010f" {
		t.Errorf("unexpected value = %v", got)
		return
	}
}

func TestGetVaruint(t *testing.T) {
	tests := []struct {
		arr      []byte
		expected uint64
	}{
		{
			arr:      []byte{0x0e},
			expected: 14,
		},
		{
			arr:      []byte{0x00, 0x0e},
			expected: 14,
		},
		{
			arr:      []byte{0x01, 0x0e},
			expected: 270,
		},
	}

	for _, test := range tests {
		got := GetVaruint(test.arr)
		if got != test.expected {
			t.Errorf("incorrect value %v, expected %v", got, test.expected)
			return
		}
	}
}

func TestExtractVote(t *testing.T) {
	hexString := "473044022041007ad95eaf56b4d5d2629cf96d2f968ad4488102c93ee7418bfa4c3c7a3c3d02207f0c36db3a9bb30995e88ed321cd144f612bba9ab84d0fdd28f0eb7a17055ffb014c82473045022100e77bd5fe006cd973f8d231ffeb26b80eddbeab93ed9f382f55276ddf98c882840220027a8ad1efcc41606dd77b48f5658bf8270fb51ada97b2b36d6a60649280c6c114a57e1e892f3031232356ecfddb0102557999357120c057ed3a9a722f0ebab32aa4231df45deeec409f45206f0942d9103d9504d556010e010f"
	vote := ExtractVoteFromSignature(hexString)
	if vote.ValidatorAddress != "a57e1e892f3031232356ecfddb01025579993571" {
		t.Errorf("unexpected ValidatorAddress = %v", vote.ValidatorAddress)
		return
	}

	if vote.TargetHash != "c057ed3a9a722f0ebab32aa4231df45deeec409f45206f0942d9103d9504d556" {
		t.Errorf("unexpected targetHash = %v", vote.TargetHash)
		return
	}

	if vote.SourceEpoch != 0x0e {
		t.Errorf("unexpected SourceEpoch = %v", vote.SourceEpoch)
		return
	}

	if vote.TargetEpoch != 0x0f {
		t.Errorf("unexpected TargetEpoch = %v", vote.TargetEpoch)
		return
	}
}
