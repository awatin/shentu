package simulation

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cert/types"
)

func TestDecodeStore(t *testing.T) {
	app := simapp.Setup(false)
	cdc := app.AppCodec()

	rand.Seed(time.Now().UnixNano())

	certifier := types.Certifier{
		Address:     RandomAccount().Address.String(),
		Proposer:    RandomAccount().Address.String(),
		Description: "this is a test case.",
	}

	validatorPubKey := RandomAccount().PubKey
	pkAny, err := codectypes.PackAny(validatorPubKey)
	if err != nil {
		panic(err)
	}
	validator := types.Validator{
		Pubkey:    pkAny,
		Certifier: RandomAccount().Address.String(),
	}

	platformPubKey := RandomAccount().PubKey
	pkAny, err = codectypes.PackAny(platformPubKey)
	if err != nil {
		panic(err)
	}
	platform := types.Platform{
		ValidatorPubkey:   pkAny,
		Description: "This is a test case.",
	}

	library := types.Library{
		Address:   sdk.AccAddress("f23908hf932").String(),
		Publisher: sdk.AccAddress("0092uf32").String(),
	}

	aliasCertifier := types.Certifier{
		Address:     RandomAccount().Address.String(),
		Alias:       "Alice",
		Proposer:    RandomAccount().Address.String(),
		Description: "this is a test case.",
	}

	certifierAddr, err := sdk.AccAddressFromBech32(certifier.Address)
	require.NoError(t, err)
	libraryAddr, err := sdk.AccAddressFromBech32(library.Address)
	require.NoError(t, err)
	
	kvPairs := []kv.Pair{
		kv.Pair{Key: types.CertifierStoreKey(certifierAddr), Value: cdc.MustMarshalBinaryLengthPrefixed(&certifier)},
		kv.Pair{Key: types.ValidatorStoreKey(validatorPubKey), Value: cdc.MustMarshalBinaryLengthPrefixed(&validator)},
		kv.Pair{Key: types.PlatformStoreKey(platformPubKey), Value: cdc.MustMarshalBinaryLengthPrefixed(&platform)},
		kv.Pair{Key: types.LibraryStoreKey(libraryAddr), Value: cdc.MustMarshalBinaryLengthPrefixed(&library)},
		kv.Pair{Key: types.CertifierAliasStoreKey(aliasCertifier.Alias), Value: cdc.MustMarshalBinaryLengthPrefixed(&aliasCertifier)},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Certifier", fmt.Sprintf("%v\n%v", certifier, certifier)},
		{"Validator", fmt.Sprintf("%v\n%v", validator, validator)},
		{"Platform", fmt.Sprintf("%v\n%v", platform, platform)},
		{"Library", fmt.Sprintf("%v\n%v", library, library)},
		{"Alias certifier", fmt.Sprintf("%v\n%v", aliasCertifier, aliasCertifier)},
		{"other", ""},
	}

	decoder := NewDecodeStore(cdc)
	
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 {
				require.Panics(t, func() { decoder(kvPairs[i], kvPairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, decoder(kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}

// RandomAccount generates a random Account object.
func RandomAccount() simtypes.Account {
	privkeySeed := make([]byte, 15)
	rand.Read(privkeySeed)

	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	pubKey := privKey.PubKey()
	address := sdk.AccAddress(pubKey.Address())

	return simtypes.Account{
		PrivKey: privKey,
		PubKey:  pubKey,
		Address: address,
	}
}
