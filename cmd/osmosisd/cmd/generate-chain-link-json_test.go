package cmd_test

import (
	"encoding/hex"
	"fmt"
	"github.com/osmosis-labs/osmosis/app/params"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/desmos-labs/desmos/app"
	profilescliutils "github.com/desmos-labs/desmos/x/profiles/client/utils"
	"github.com/desmos-labs/desmos/x/profiles/types"
	"github.com/osmosis-labs/osmosis/cmd/osmosisd/cmd"
)

func TestGetGenerateChainlinkJsonCmd(t *testing.T) {
	cfg := sdk.GetConfig()
	app.SetupConfig(cfg)

	keyBase := keyring.NewInMemory()
	algo := hd.Secp256k1
	hdPath := sdk.GetConfig().GetFullFundraiserPath()

	keyName := "test"
	mnemonic := "clip toilet stairs jaguar baby over mosquito capital speed mule adjust eye print voyage verify smart open crack imitate auto gauge museum planet rebel"
	_, err := keyBase.NewAccount(keyName, mnemonic, "", hdPath, algo)
	require.NoError(t, err)

	output := os.Stdout
	clientCtx := client.Context{}.
		WithKeyring(keyBase).
		WithOutput(output)

	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd.GetGenerateChainlinkJSONCmd(), []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, keyName),
		fmt.Sprintf("--%s=%s", "filename", ""),
	})
	require.NoError(t, err)

	key, err := keyBase.Key(keyName)
	addr, _ := sdk.Bech32ifyAddressBytes(params.Bech32PrefixAccAddr, key.GetAddress())
	sig, pubkey, err := clientCtx.Keyring.Sign(keyName, []byte(addr))
	require.NoError(t, err)

	cdc, _ := app.MakeCodecs()
	var data profilescliutils.ChainLinkJSON
	err = cdc.UnmarshalJSON(out.Bytes(), &data)
	require.NoError(t, err)

	expected := profilescliutils.NewChainLinkJSON(
		types.NewBech32Address(addr, params.Bech32PrefixAccAddr),
		types.NewProof(pubkey, hex.EncodeToString(sig), addr),
		types.NewChainConfig(params.Bech32PrefixAccAddr),
	)

	require.Equal(t, expected, data)

}
