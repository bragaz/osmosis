package cmd

import (
	"encoding/hex"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	profilescliutils "github.com/desmos-labs/desmos/x/profiles/client/utils"
	"github.com/desmos-labs/desmos/x/profiles/types"
	"github.com/osmosis-labs/osmosis/app/params"
)

// GetGenerateChainlinkJSONCmd returns the command allowing to generate the chain link json file for creating chain link
func GetGenerateChainlinkJSONCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-chain-link-json",
		Short: "generate the chain link json for creating chain link with the key specified using the --from flag",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cdc := params.MakeEncodingConfig()

			chainLinkJSON, err := GenerateChainLinkJSON(
				clientCtx,
				params.Bech32PrefixAccAddr,
				cdc,
			)
			if err != nil {
				return err
			}

			bz, err := cdc.Marshaler.MarshalJSON(&chainLinkJSON)
			if err != nil {
				return err
			}

			filename, _ := cmd.Flags().GetString("filename")
			if strings.TrimSpace(filename) != "" {
				if err := ioutil.WriteFile("data.json", bz, 0600); err != nil {
					return err
				}
			}
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("filename", "data.json", "The name of output chain link json file. It does not generate the file if it is empty.")
	return cmd
}

// GenerateChainLinkJSON returns ChainLinkJSON instance for creating chain link
func GenerateChainLinkJSON(clientCtx client.Context, prefix string, cdc params.EncodingConfig) (profilescliutils.ChainLinkJSON, error) {

	// generate signature
	addr, _ := sdk.Bech32ifyAddressBytes(prefix, clientCtx.GetFromAddress())
	sig, pubkey, err := clientCtx.Keyring.Sign(clientCtx.GetFromName(), []byte(addr))
	if err != nil {
		return profilescliutils.ChainLinkJSON{}, err
	}

	chainLinkJSON := profilescliutils.NewChainLinkJSON(
		types.NewBech32Address(addr, prefix),
		types.NewProof(pubkey, hex.EncodeToString(sig), addr),
		types.NewChainConfig(prefix),
	)
	if err := chainLinkJSON.UnpackInterfaces(cdc.Marshaler); err != nil {
		return profilescliutils.ChainLinkJSON{}, err
	}
	return chainLinkJSON, nil
}
