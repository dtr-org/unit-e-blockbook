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
	testTx1, testCommitTx, testAdminTx bchain.Tx

	testTxPacked1    = "0200030001d86cd422c107fe6c0d4c3ac462ec33407f7200dcc00ffc253ed47ffe33bef30800000000cc483045022100bacc2a56b0f64c2049867eccb84c82274531f9e83803dac9c175c15fcc4a21fb022045ff57194e4ed1e314ac78ce78ac8d6829f593d3adf129d7ee591f5b4bed987e014c8146304402200b6459709d33acf8185c5f934df3cd3cab138fde917491d0acbb5820ba92d887022025159929618327f19c966655005f22069e40d7819e6b35dd72d39661e880f71314a57e1e892f3031232356ecfddb01025579993571208ced6f6944123a69ff006ad75d72c570bb9c9b247a11dce426359c0dbadaf125010d010effffffff016844b2ec22000000672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb36351672103d3a0ea21c92eb4be4b3f6699cd5a8939aab1fe71a18378dbb4e19fef67e3b1ffb46476a914a57e1e892f3031232356ecfddb0102557999357188ac6751686800000000"
	testCommitPacked = "0a20917dc015d56c9411abfb8e0cdcacf4cc299b3afc35aa7dd2277893124540aa2b10001ac8020200030001861807bdc996daea3eefd695b556dde8e9675a020ad2d07b0afa2783d3b6470b00000000cc473044022033e7eec7f60fc5ae1ec9d35a430043b6d4b352f628760b45c0a612e38945363a022022c8bf85df13c9555d71ab3897d434ec14b1dd12537c44a2122146c70e9d51ff014c8247304502210089d29935624a10455a6f1bc2c39acf6dd3557abb236dece4f70a7021ee6933b302202e5ee0bd83b9f45bfb4fff0b05024e7c53578244f0fa6b41af05f8c2cc652c34147bc15980b3c7bfd125f3cc02cc624c81bf352f4e20d310485d7f9dd9a65cb480e7c55434dfca99c587304c51509b0f8d2597e72f63011f0120ffffffff0108366352bfc60100402103923e38f1f92061da3f7c8387f72cbbafb0f390d8c7a2ddc18aab94273a1eaa88b363516776a9147bc15980b3c7bfd125f3cc02cc624c81bf352f4e88ac6800000000209cf7d7e505280030f9ea113afb010a0012200b47b6d38327fa0a7bd0d20a025a67e9e8dd56b595d6ef3eeada96c9bd071886180022cc01473044022033e7eec7f60fc5ae1ec9d35a430043b6d4b352f628760b45c0a612e38945363a022022c8bf85df13c9555d71ab3897d434ec14b1dd12537c44a2122146c70e9d51ff014c8247304502210089d29935624a10455a6f1bc2c39acf6dd3557abb236dece4f70a7021ee6933b302202e5ee0bd83b9f45bfb4fff0b05024e7c53578244f0fa6b41af05f8c2cc652c34147bc15980b3c7bfd125f3cc02cc624c81bf352f4e20d310485d7f9dd9a65cb480e7c55434dfca99c587304c51509b0f8d2597e72f63011f012028ffffffff0f42790a0701c6bf5263360810001a402103923e38f1f92061da3f7c8387f72cbbafb0f390d8c7a2ddc18aab94273a1eaa88b363516776a9147bc15980b3c7bfd125f3cc02cc624c81bf352f4e88ac682206636f6d6d69742a226d726f4b32704a454c4253314c4a72424136696246416b6b466f736974684a756b61"
	testAdminPacked  = "0a20ad8c1088b1a5bf59da97f12cb777826f7637963560fb055aaa3bc0145a3f24e510001ab805020007000001015416d589fd46d73c47980491fe8805a4e8b824db3475eaa4d6a61b417d272ae20100000023220020e9875ec09fdcb7e709dbba6f6f188bcf974b2ca37123dcf05b55ad06aa4c87f1ffffffff020000000000000000fd38016a4d3401000921033eeb725e33b37981e08f88afdb267e11ae96f30f69bba47be214b9a90c7b85202103c0a816c8a56eaf0c6c0776a90237ee24f466bc7b65726f2e7e3bf101413f98212103abb99839a7235381f7bf9be880498706d4a6cc9427ce4bd5126f7ec3e9bc6ccd2102f5f4f46b504604625d8326e7e42956cb5daa63c8405cff478ef91797df7a44a321022cfc9cbf9f8fe3425519931d54a56208baa88fdb1295d77418c4dd9ca59516ca2102f8e734e3cef8f986d908c46414dad5257372604268b21c1f4911dc685a9c934d21023357cf75a6e4abe9b420d487880b1db18e8e9295e40ab50532338b3cf12e893221020eba4f980f90d862520f0334c203007db5fc43b2122686042c30032baf30c2552103923e38f1f92061da3f7c8387f72cbbafb0f390d8c7a2ddc18aab94273a1eaa88c09ee6050000000017a914f558a88dc06faccf65dbc69dca31b11326238ef8870400473044022071107a54a3943d702732c4cf47a8f32fc1d2540065aacabb913bf947ad975b1c022064b4fa901be0b845d9acae5ea31b948e22cefd585d77b15de32adb25f959fb3e01483045022100cf0297a83e37fd72b21c0191335383943e862c8a78490e21d18dacd58961a62f02201be56517f54294521d778850ac5de484ad66073920d7e4ba5b8d71b3b85e711d0169522102630a75cd35adc6c44ca677e83feb8e4a7e539baaa49887c455e8242e3e3b1c052103946025d10e3cdb30a9cd73525bc9acc4bd92e184cdd9c9ea7d0ebc6b654afcc5210290f45494a197cbd389181b3d7596a90499a93368159e8a6e9f9d0d460799d33d53ae00000000209cf7d7e505280030f9ea113a510a001220e22a277d411ba6d6a4ea7534db24b8e8a40588fe910498473cd746fd89d5165418002223220020e9875ec09fdcb7e709dbba6f6f188bcf974b2ca37123dcf05b55ad06aa4c87f128ffffffff0f42c70210001ab8026a4d3401000921033eeb725e33b37981e08f88afdb267e11ae96f30f69bba47be214b9a90c7b85202103c0a816c8a56eaf0c6c0776a90237ee24f466bc7b65726f2e7e3bf101413f98212103abb99839a7235381f7bf9be880498706d4a6cc9427ce4bd5126f7ec3e9bc6ccd2102f5f4f46b504604625d8326e7e42956cb5daa63c8405cff478ef91797df7a44a321022cfc9cbf9f8fe3425519931d54a56208baa88fdb1295d77418c4dd9ca59516ca2102f8e734e3cef8f986d908c46414dad5257372604268b21c1f4911dc685a9c934d21023357cf75a6e4abe9b420d487880b1db18e8e9295e40ab50532338b3cf12e893221020eba4f980f90d862520f0334c203007db5fc43b2122686042c30032baf30c2552103923e38f1f92061da3f7c8387f72cbbafb0f390d8c7a2ddc18aab94273a1eaa8822086e756c6c6461746142520a0405e69ec010011a17a914f558a88dc06faccf65dbc69dca31b11326238ef887220a736372697074686173682a23324e4663566b4e704655387877654346415161716576637532474c545a6f4356567956"
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

	testCommitTx = bchain.Tx{
		Hex:       "0200030001861807bdc996daea3eefd695b556dde8e9675a020ad2d07b0afa2783d3b6470b00000000cc473044022033e7eec7f60fc5ae1ec9d35a430043b6d4b352f628760b45c0a612e38945363a022022c8bf85df13c9555d71ab3897d434ec14b1dd12537c44a2122146c70e9d51ff014c8247304502210089d29935624a10455a6f1bc2c39acf6dd3557abb236dece4f70a7021ee6933b302202e5ee0bd83b9f45bfb4fff0b05024e7c53578244f0fa6b41af05f8c2cc652c34147bc15980b3c7bfd125f3cc02cc624c81bf352f4e20d310485d7f9dd9a65cb480e7c55434dfca99c587304c51509b0f8d2597e72f63011f0120ffffffff0108366352bfc60100402103923e38f1f92061da3f7c8387f72cbbafb0f390d8c7a2ddc18aab94273a1eaa88b363516776a9147bc15980b3c7bfd125f3cc02cc624c81bf352f4e88ac6800000000",
		Blocktime: 1555430300,
		Time:      1555430300,
		Txid:      "917dc015d56c9411abfb8e0cdcacf4cc299b3afc35aa7dd2277893124540aa2b",
		LockTime:  0,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "473044022033e7eec7f60fc5ae1ec9d35a430043b6d4b352f628760b45c0a612e38945363a022022c8bf85df13c9555d71ab3897d434ec14b1dd12537c44a2122146c70e9d51ff014c8247304502210089d29935624a10455a6f1bc2c39acf6dd3557abb236dece4f70a7021ee6933b302202e5ee0bd83b9f45bfb4fff0b05024e7c53578244f0fa6b41af05f8c2cc652c34147bc15980b3c7bfd125f3cc02cc624c81bf352f4e20d310485d7f9dd9a65cb480e7c55434dfca99c587304c51509b0f8d2597e72f63011f0120",
				},
				Txid:     "0b47b6d38327fa0a7bd0d20a025a67e9e8dd56b595d6ef3eeada96c9bd071886",
				Vout:     0,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(499999999997448),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "2103923e38f1f92061da3f7c8387f72cbbafb0f390d8c7a2ddc18aab94273a1eaa88b363516776a9147bc15980b3c7bfd125f3cc02cc624c81bf352f4e88ac68",
					Addresses: []string{
						"mroK2pJELBS1LJrBA6ibFAkkFosithJuka",
					},
					Type: "commit",
				},
			},
		},
	}

	testAdminTx = bchain.Tx{
		Hex:       "020007000001015416d589fd46d73c47980491fe8805a4e8b824db3475eaa4d6a61b417d272ae20100000023220020e9875ec09fdcb7e709dbba6f6f188bcf974b2ca37123dcf05b55ad06aa4c87f1ffffffff020000000000000000fd38016a4d3401000921033eeb725e33b37981e08f88afdb267e11ae96f30f69bba47be214b9a90c7b85202103c0a816c8a56eaf0c6c0776a90237ee24f466bc7b65726f2e7e3bf101413f98212103abb99839a7235381f7bf9be880498706d4a6cc9427ce4bd5126f7ec3e9bc6ccd2102f5f4f46b504604625d8326e7e42956cb5daa63c8405cff478ef91797df7a44a321022cfc9cbf9f8fe3425519931d54a56208baa88fdb1295d77418c4dd9ca59516ca2102f8e734e3cef8f986d908c46414dad5257372604268b21c1f4911dc685a9c934d21023357cf75a6e4abe9b420d487880b1db18e8e9295e40ab50532338b3cf12e893221020eba4f980f90d862520f0334c203007db5fc43b2122686042c30032baf30c2552103923e38f1f92061da3f7c8387f72cbbafb0f390d8c7a2ddc18aab94273a1eaa88c09ee6050000000017a914f558a88dc06faccf65dbc69dca31b11326238ef8870400473044022071107a54a3943d702732c4cf47a8f32fc1d2540065aacabb913bf947ad975b1c022064b4fa901be0b845d9acae5ea31b948e22cefd585d77b15de32adb25f959fb3e01483045022100cf0297a83e37fd72b21c0191335383943e862c8a78490e21d18dacd58961a62f02201be56517f54294521d778850ac5de484ad66073920d7e4ba5b8d71b3b85e711d0169522102630a75cd35adc6c44ca677e83feb8e4a7e539baaa49887c455e8242e3e3b1c052103946025d10e3cdb30a9cd73525bc9acc4bd92e184cdd9c9ea7d0ebc6b654afcc5210290f45494a197cbd389181b3d7596a90499a93368159e8a6e9f9d0d460799d33d53ae00000000",
		Blocktime: 1555430300,
		Time:      1555430300,
		Txid:      "ad8c1088b1a5bf59da97f12cb777826f7637963560fb055aaa3bc0145a3f24e5",
		LockTime:  0,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "220020e9875ec09fdcb7e709dbba6f6f188bcf974b2ca37123dcf05b55ad06aa4c87f1",
				},
				Txid:     "e22a277d411ba6d6a4ea7534db24b8e8a40588fe910498473cd746fd89d51654",
				Vout:     0,
				Sequence: 4294967295,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(0),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex:  "6a4d3401000921033eeb725e33b37981e08f88afdb267e11ae96f30f69bba47be214b9a90c7b85202103c0a816c8a56eaf0c6c0776a90237ee24f466bc7b65726f2e7e3bf101413f98212103abb99839a7235381f7bf9be880498706d4a6cc9427ce4bd5126f7ec3e9bc6ccd2102f5f4f46b504604625d8326e7e42956cb5daa63c8405cff478ef91797df7a44a321022cfc9cbf9f8fe3425519931d54a56208baa88fdb1295d77418c4dd9ca59516ca2102f8e734e3cef8f986d908c46414dad5257372604268b21c1f4911dc685a9c934d21023357cf75a6e4abe9b420d487880b1db18e8e9295e40ab50532338b3cf12e893221020eba4f980f90d862520f0334c203007db5fc43b2122686042c30032baf30c2552103923e38f1f92061da3f7c8387f72cbbafb0f390d8c7a2ddc18aab94273a1eaa88",
					Type: "nulldata",
				},
			},
			{
				ValueSat: *big.NewInt(99000000),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "a914f558a88dc06faccf65dbc69dca31b11326238ef887",
					Addresses: []string{
						"2NFcVkNpFU8xweCFAQaqevcu2GLTZoCVVyV",
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
	if opRetAddr[0] != OpReturn {
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
	if hex.EncodeToString(opRetAddr) != "6aaa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2" {
		t.Errorf("OpRetAddr is incorrect %s", hex.EncodeToString(opRetAddr))
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

	if hex.EncodeToString(desc) != "6aaa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2" {
		t.Errorf("Desc from Vout is incorrect: %s", hex.EncodeToString(desc))
	}

	desc, err = parser.GetAddrDescFromAddress(string(desc))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
		return
	}

	if hex.EncodeToString(desc) != "6aaa21a9ed941b41ad0b6a2c6b93d98dd32ff22ae87dcf3b95b721102f5956b9d7a3d299f2" {
		t.Errorf("Desc from Address is incorrect: %s", hex.EncodeToString(desc))
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
		t.Errorf("Address from Desc is incorrect: %s", addrs)
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
		vout   []bchain.Vout
		parser *UniteParser
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ute-1",
			args: args{
				vout:   testTx1.Vout,
				parser: NewUniteParser(btc.GetChainParams("regtest"), &btc.Configuration{}),
			},
		},
		{
			name: "ute-commit",
			args: args{
				vout:   testCommitTx.Vout,
				parser: NewUniteParser(GetChainParams("test"), &btc.Configuration{}),
			},
		},
		{
			name: "ute-admin",
			args: args{
				vout:   testAdminTx.Vout,
				parser: NewUniteParser(GetChainParams("test"), &btc.Configuration{}),
			},
		},
		{
			name: "ute-rs-keyhash",
			args: args{
				vout: []bchain.Vout{
					bchain.Vout{
						ValueSat: *big.NewInt(977517107),
						N:        1,
						ScriptPubKey: bchain.ScriptPubKey{
							Hex:  "5114ca305850196fe3286cc480fa6e9dd445512755d920bb258d8efb687e03537da2a4fc9333ba76346ef4de46e91e83db728ca69c0a0c",
							Type: "witness_v1_remotestake_keyhash",
							Addresses: []string{
								"n3Qu8Gd7VjUeg1V1wLEeSHjDKqbaJjafc5",
							},
						},
					},
				},
				parser: NewUniteParser(GetChainParams("regtest"), &btc.Configuration{}),
			},
		},
		{
			name: "ute-rs-scripthash",
			args: args{
				vout: []bchain.Vout{
					bchain.Vout{
						ValueSat: *big.NewInt(1),
						N:        0,
						ScriptPubKey: bchain.ScriptPubKey{
							Hex:  "521409908c6b1618b28f9ba847f1d986a2c9980b734d2047fc43af4ebfee85d46e5af86e16573dcd6b1cd4713e30493d8b16b3d5884c77",
							Type: "witness_v2_remotestake_scripthash",
							Addresses: []string{
								"uert1qgl7y8t6whlhgt4rwttuxu9jh8hxkk8x5wylrqjfa3vtt84vgf3ms2wvjrk",
							},
						},
					},
				},
				parser: NewUniteParser(GetChainParams("regtest"), &btc.Configuration{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for n, vout := range tt.args.vout {
				if len(vout.ScriptPubKey.Addresses) != 0 {
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
					if got3[0] != vout.ScriptPubKey.Addresses[0] {
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
			name: "ute-commit",
			args: args{
				tx:        testCommitTx,
				height:    292217,
				blockTime: 1555430300,
				parser:    NewUniteParser(GetChainParams("test"), &btc.Configuration{}),
			},
			want:    testCommitPacked,
			wantErr: false,
		},
		{
			name: "ute-admin",
			args: args{
				tx:        testAdminTx,
				height:    292217,
				blockTime: 1555430300,
				parser:    NewUniteParser(GetChainParams("test"), &btc.Configuration{}),
			},
			want:    testAdminPacked,
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
			name: "ute-commit",
			args: args{
				packedTx: testCommitPacked,
				parser:   NewUniteParser(GetChainParams("test"), &btc.Configuration{}),
			},
			want:    &testCommitTx,
			want1:   292217,
			wantErr: false,
		},
		{
			name: "ute-admin",
			args: args{
				packedTx: testAdminPacked,
				parser:   NewUniteParser(GetChainParams("test"), &btc.Configuration{}),
			},
			want:    &testAdminTx,
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
			t.Errorf("test case %d: unexpected parsing error %v", i, err)
			return
		}

		noff, op, err := GetOp(script, test.off)
		if test.expectedErr != nil {
			if err == nil {
				t.Errorf("test case %d: did not get error, want %v", i, test.expectedErr)
			}
			if err.Error() != test.expectedErr.Error() {
				t.Errorf("test case %d: unexpected error %v, want %v", i, err, test.expectedErr)
				return
			}
			continue
		}

		if err != nil {
			t.Errorf("test case %d: unexpected error %v", i, err)
			return
		}

		if noff != test.noff {
			t.Errorf("test case %d: unexpected new offset %d, want %d", i, noff, test.noff)
			return
		}

		got := hex.EncodeToString(op)
		if got != test.expected {
			t.Errorf("test case %d: unexpected value %s, want %s", i, got, test.expected)
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
		{
			arr:      []byte{0x07, 0x21, 0xda, 0x0e},
			expected: 0x0721da0e,
		},
		{
			arr:      []byte{0xdf, 0x33, 0xda, 0x0e, 0x51, 0xca, 0xcd, 0x05},
			expected: 0xdf33da0e51cacd05,
		},
	}

	for _, test := range tests {
		got := GetVaruint(test.arr)
		if got != test.expected {
			t.Errorf("incorrect value %d, expected %d", got, test.expected)
			return
		}
	}
}

func TestExtractVote(t *testing.T) {
	tests := []struct {
		hex      string
		expected Vote
	}{
		{
			hex: "473044022041007ad95eaf56b4d5d2629cf96d2f968ad4488102c93ee7418bfa4c3c7a3c3d02207f0c36db3a9bb30995e88ed321cd144f612bba9ab84d0fdd28f0eb7a17055ffb014c82473045022100e77bd5fe006cd973f8d231ffeb26b80eddbeab93ed9f382f55276ddf98c882840220027a8ad1efcc41606dd77b48f5658bf8270fb51ada97b2b36d6a60649280c6c114a57e1e892f3031232356ecfddb0102557999357120c057ed3a9a722f0ebab32aa4231df45deeec409f45206f0942d9103d9504d556010e010f",
			expected: Vote{
				ValidatorAddress: "71359979550201dbfdec56232331302f891e7ea5",
				TargetHash:       "56d504953d10d942096f20459f40ecee5df41d23a42ab3ba0e2f729a3aed57c0",
				SourceEpoch:      0x0e,
				TargetEpoch:      0x0f,
			},
		},
		{
			hex: "47304402204b9bb63f9b055a7d82841f064167df5d9b774f91a5d76eb807559a03f51dc39f02203af15ccb70a77801afdac05ef1723b07a59da9d3b19a4ced37e53cdc9a0db1bc014c5b200303b3fd308c053db885ce2ea666be319aa68af96bbdf383573d150053895c1114e1d0077c611cdb8e4bc5314e272f3b74edef8dc3205abcb1b1868582266bdd2d683ece9396cc44673085e0738bc5f173ef8e248912010a0164",
			expected: Vote{
				ValidatorAddress: "c38defed743b2f274e31c54b8edb1c617c07d0e1",
				TargetHash:       "1289248eef73f1c58b73e085306744cc9693ce3e682ddd6b26828586b1b1bc5a",
				SourceEpoch:      10,
				TargetEpoch:      100,
			},
		},
		{
			hex: "483045022100853efdc85fda74644e34f02ceda60aa2134c75c511b8eb62ad75e07e500ae77f02200fc66eb6745be8fac148c5d8819c76f86452a8920595b114cc71ddbdde0796f4014c83463044022001e9c6109b20e13d9b6bc802d425c995bbff082dd2c742e40a3355586ccc6a690220219ef555017536161085fe336a268779a871a49f61a4d2de6839cceb152af972146d68510f7d530a363dff1f727ac3649a7d73f4a7202ca64a150c31b460b93b83ea70d4e757f0d660f70804de17f1af9e369a7c799802ce0002cf00",
			expected: Vote{
				ValidatorAddress: "a7f4737d9a64c37a721fff3d360a537d0f51686d",
				TargetHash:       "98797c9a369eaff117de0408f760d6f057e7d470ea833bb960b4310c154aa62c",
				SourceEpoch:      206,
				TargetEpoch:      207,
			},
		},
	}
	for _, test := range tests {
		got := ExtractVoteFromSignature(test.hex)
		if *got != test.expected {
			t.Errorf("incorrect vote %v, expected %v", got, test.expected)
			return
		}
	}
}

func TestExtractSlash(t *testing.T) {
	hexString := "473044022051a23d4aef6086a272a8a81931d624d079d315babb626c52e3da0eee438c04ca02205420e72ab7a758abb36a963fc0edba88e285b77d80f4c1e205c0b19356727b6a014c81463044022023696de39b4aedb9e8fff73300d0378e707520084d276fa81cfa962f9bc986b1022007c3c196dab315517842adea0646018d2004909ee59a04e93c5b1cdffd9cb22a14f53ab1d04f4f13126fc30b46e2b0cdb7266dda20207414ff6c2cd52e9cca363cf2455e5a51d694ad2e165a6671caad40b900f68201010e01184c82473045022100bf2d6061ee1c279dc6f18c344484b1d07d297b8b53b60eefe7259935154bfa8202205455061688af7f05b602150e66ad45e7e9e16b75e24ba0bbde4df6635fcbbacb14f53ab1d04f4f13126fc30b46e2b0cdb7266dda20204b1708ad7c03cd582d1aadf92be80338365b2822be481bf3bd3a07abd47035e1010e0118"
	slash := ExtractSlashFromSignature(hexString)

	expectedVote1 := Vote{ValidatorAddress: "20da6d26b7cdb0e2460bc36f12134f4fd0b13af5", TargetHash: "0182f600b940adca71665a162ead94d6515a5e45f23c36ca9c2ed52c6cff1474", SourceEpoch: 14, TargetEpoch: 24}

	expectedVote2 := Vote{ValidatorAddress: "20da6d26b7cdb0e2460bc36f12134f4fd0b13af5", TargetHash: "e13570d4ab073abdf31b48be22285b363803e82bf9ad1a2d58cd037cad08174b", SourceEpoch: 14, TargetEpoch: 24}

	if *slash.Vote1 != expectedVote1 {
		t.Errorf("first vote %v is not as expected %v", slash.Vote1, expectedVote1)
		return
	}

	if *slash.Vote2 != expectedVote2 {
		t.Errorf("second vote %v is not as expected %v", slash.Vote2, expectedVote2)
		return
	}

}
