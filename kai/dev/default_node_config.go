// Defines default configs used for initializing nodes in dev settings.

package dev

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/kardiachain/go-kardia/lib/crypto"
	"github.com/kardiachain/go-kardia/lib/log"
	"github.com/kardiachain/go-kardia/types"
)

type DevNodeConfig struct {
	PrivKey     *ecdsa.PrivateKey
	VotingPower int64
	NodeID      string
}

type DevEnvironmentConfig struct {
	DevNodeSet []DevNodeConfig

	proposalIndex  int
	VotingStrategy map[VoteTurn]int
}

type node struct {
	key         string
	votingPower int64
	nodeID      string
}

type VoteTurn struct {
	Height   int
	Round    int
	VoteType int
}

type account struct {
	address string
	balance int64
}

const (
	// password is used to get keystore
	password  = "KardiaChain"
	ChainData = "chaindata"
	DbCache   = 16
	DbHandles = 16
)

// GenesisAccounts are used to initialized accounts in genesis block
var GenesisAccounts = map[string]int64{
	"0xc1fe56E3F58D3244F606306611a5d10c8333f1f6": 100000000,
	"0x7cefC13B6E2aedEeDFB7Cb6c32457240746BAEe5": 100000000,
	"0xfF3dac4f04dDbD24dE5D6039F90596F0a8bb08fd": 100000000,
	"0x071E8F5ddddd9f2D4B4Bdf8Fc970DFe8d9871c28": 100000000,
	"0x94FD535AAB6C01302147Be7819D07817647f7B63": 100000000,
	"0xa8073C95521a6Db54f4b5ca31a04773B093e9274": 100000000,
	"0xe94517a4f6f45e80CbAaFfBb0b845F4c0FDD7547": 100000000,
	"0xBA30505351c17F4c818d94a990eDeD95e166474b": 100000000,
	"0x212a83C0D7Db5C526303f873D9CeaA32382b55D0": 100000000,
	// TODO(namdoh): Re-enable after parsing node index fixed in main.go
	//"0x36BE7365e6037bD0FDa455DC4d197B07A2002547": 100000000,
}

//  GenesisAddrKeys maps genesis account addresses to private keys.
var GenesisAddrKeys = map[string]string{
	"0xc1fe56E3F58D3244F606306611a5d10c8333f1f6": "8843ebcb1021b00ae9a644db6617f9c6d870e5fd53624cefe374c1d2d710fd06",
	"0x7cefC13B6E2aedEeDFB7Cb6c32457240746BAEe5": "77cfc693f7861a6e1ea817c593c04fbc9b63d4d3146c5753c008cfc67cffca79",
	"0xfF3dac4f04dDbD24dE5D6039F90596F0a8bb08fd": "98de1df1e242afb02bd5dc01fbcacddcc9a4d41df95a66f629139560ca6e4dbb",
	"0x071E8F5ddddd9f2D4B4Bdf8Fc970DFe8d9871c28": "32f5c0aef7f9172044a472478421c63fd8492640ff2d0eaab9562389db3a8efe",
	"0x94FD535AAB6C01302147Be7819D07817647f7B63": "68b53a92d846baafdc782cb9cad65d77020c8d747eca7b621370b52b18c91f9a",
	"0xa8073C95521a6Db54f4b5ca31a04773B093e9274": "049de018e08c3bcd59c1a21f0cf7de8f17fe51f8ce7d9c2120d17b1f0251b265",
	"0xe94517a4f6f45e80CbAaFfBb0b845F4c0FDD7547": "9fdd56a3c2a536dc8f981d935f0f3f2ea04e125547fdfffa37e157ce86ff1007",
	"0xBA30505351c17F4c818d94a990eDeD95e166474b": "ae1a52546294bed6e734185775dbc84009de00bdf51b709471e2415c31ceeed7",
	"0x212a83C0D7Db5C526303f873D9CeaA32382b55D0": "b34bd81838a4a335fb3403d0bf616eca1eb9a4b4716c7dda7c617503cfeaab67",

	// TODO(namdoh): Re-enable after parsing node index fixed in main.go
	//"e049a09c992c882bc2deb780323a247c6ee0951f8b4c5c1dd0fc2fc22ce6493d": "0x36BE7365e6037bD0FDa455DC4d197B07A2002547",
}

var nodes = []node{
	{"8843ebcb1021b00ae9a644db6617f9c6d870e5fd53624cefe374c1d2d710fd06", 100, "enode://724fbdc7067814bdd60315d836f28175ff9c72e4e1d86513a2b578f9cd769e688d6337550778b89e4861a42580613f1f1dec23f17f7a1627aa99104cc4204eb1@[::]:3000"},
	{"77cfc693f7861a6e1ea817c593c04fbc9b63d4d3146c5753c008cfc67cffca79", 100, "enode://b5997edd19d36bd5a1afa701563ca5b505bacb55c840ba35aa961c46af6484c4a5069c29062fa1c9fe6754ebdb42f6b16a549a2832d080f5faf13b42993d31b8@[::]:3001"},
	{"98de1df1e242afb02bd5dc01fbcacddcc9a4d41df95a66f629139560ca6e4dbb", 100, "enode://a840f4fa933e741f0c2f4dcafc6e0bbd05a364fa7ee61a312dcefbfb9a7081a5f3028419e6d94acfd95af2e3a8c884bf48d286f3015d1f7db7a4ba6030b3e66a@[::]:3002"},
	{"32f5c0aef7f9172044a472478421c63fd8492640ff2d0eaab9562389db3a8efe", 100, "enode://e00744773f3f3d641df1df591eb1212a80c6ecc49a1f1d43ee8e30e602259726274fb1817b346b7a8fc61db5cfb531c6827044cc4c7f6d2b67dda03c5cee8b8e@[::]:3003"},
	{"68b53a92d846baafdc782cb9cad65d77020c8d747eca7b621370b52b18c91f9a", 100, "enode://2773b8f006193853112e529dbab168e4c62da80beebfa151b82547ca8ba54a5cef94b30ecd3c4dd5bf54100fc642062dcec87ae98c08e504e5c985038f994325@[::]:3004"},
	{"049de018e08c3bcd59c1a21f0cf7de8f17fe51f8ce7d9c2120d17b1f0251b265", 100, "enode://754d2ab0c56c963a64f08a11d4c5d92436e81b81a1422959bc0a5ae802b099499188646a6e6ecc3878aa0a7e4edf599fc7101d4556ae4f818c7de05bd8e810b7@[::]:3005"},
	{"9fdd56a3c2a536dc8f981d935f0f3f2ea04e125547fdfffa37e157ce86ff1007", 100, "enode://ebdfa0502b2e095d493780f50024fe6746f52893d986c74dfedb44c01f834fc68f1373df7c2035feb8542dfa86d27bbdc71c0f065545c135baf55a6a84b24870@[::]:3006"},
	{"ae1a52546294bed6e734185775dbc84009de00bdf51b709471e2415c31ceeed7", 100, "enode://ca299c6ba014473c3ac53d3e9fa92f60dc34d838095c07fff2aab350439374bc3dbfe4411757e984c6b191868898755a328c3e97fc8abb37f6eed4b73ca7f67b@[::]:3007"},
	{"b34bd81838a4a335fb3403d0bf616eca1eb9a4b4716c7dda7c617503cfeaab67", 100, "enode://e9f1abb546000dbfba59894610053e4bf337ef1db64253e5d30162c8e17ff49b9980150363e09ddeac1db6e23b4ed683ad3c1bda1ad5886112b18c47f7ad9eae@[::]:3008"},
	// TODO(namdoh): Re-enable after parsing node index fixed in main.go
	//{"e049a09c992c882bc2deb780323a247c6ee0951f8b4c5c1dd0fc2fc22ce6493d", 100, ""},
}

func CreateDevEnvironmentConfig() *DevEnvironmentConfig {
	var devEnv DevEnvironmentConfig
	devEnv.proposalIndex = 0 // Default to 0-th node as the proposer.
	devEnv.DevNodeSet = make([]DevNodeConfig, len(nodes))
	for i, n := range nodes {
		privKey, _ := crypto.ToECDSA([]byte(n.key[:32]))
		devEnv.DevNodeSet[i].PrivKey = privKey
		devEnv.DevNodeSet[i].VotingPower = n.votingPower
		devEnv.DevNodeSet[i].NodeID = n.nodeID
	}

	return &devEnv
}

func (devEnv *DevEnvironmentConfig) SetVotingStrategy(votingStrategy string) {
	if strings.HasSuffix(votingStrategy, "csv") {
		devEnv.VotingStrategy = map[VoteTurn]int{}
		csvFile, _ := os.Open(votingStrategy)
		reader := csv.NewReader(bufio.NewReader(csvFile))

		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Error("error", error)
			}
			var height, _ = strconv.Atoi(line[0])
			var round, _ = strconv.Atoi(line[1])
			var voteType, _ = strconv.Atoi(line[2])
			var result, _ = strconv.Atoi(line[3])

			var _, ok = devEnv.GetScriptedVote(height, round, voteType)
			if ok {
				log.Error(fmt.Sprintf("VoteTurn already exists with height = %v, round = %v, voteType = %v", height, round, voteType))
			} else {
				devEnv.VotingStrategy[VoteTurn{height, round, voteType}] = result
			}
		}
	}
}

func (devEnv *DevEnvironmentConfig) GetScriptedVote(height int, round int, voteType int) (int, bool) {
	if val, ok := devEnv.VotingStrategy[VoteTurn{height, round, voteType}]; ok {
		return val, ok
	}
	return 0, false
}

func (devEnv *DevEnvironmentConfig) SetProposerIndex(index int) {
	if index < 0 || index >= devEnv.GetNodeSize() {
		log.Error(fmt.Sprintf("Proposer index must be within %v and %v", 0, devEnv.GetNodeSize()))
	}
	devEnv.proposalIndex = index
}

func (devEnv *DevEnvironmentConfig) GetDevNodeConfig(index int) *DevNodeConfig {
	return &devEnv.DevNodeSet[index]
}

func (devEnv *DevEnvironmentConfig) GetNodeSize() int {
	return len(devEnv.DevNodeSet)
}

func (devEnv *DevEnvironmentConfig) GetValidatorSet(numVal int) *types.ValidatorSet {
	if numVal < 0 || numVal >= devEnv.GetNodeSize() {
		log.Error(fmt.Sprintf("Number of validator must be within %v and %v", 0, devEnv.GetNodeSize()))
	}
	validators := make([]*types.Validator, numVal)
	for i := 0; i < numVal; i++ {
		node := devEnv.DevNodeSet[i]
		validators[i] = types.NewValidator(node.PrivKey.PublicKey, node.VotingPower)
	}

	validatorSet := types.NewValidatorSet(validators)
	validatorSet.TurnOnKeepSameProposer()
	validatorSet.SetProposer(validators[devEnv.proposalIndex])
	return validatorSet
}