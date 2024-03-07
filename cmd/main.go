package main

import (
	"fmt"
	"github.com/airchains-network/decentralized-sequencer/command"
	"github.com/airchains-network/decentralized-sequencer/command/keys"
	"github.com/airchains-network/decentralized-sequencer/command/zkpCmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "station-trackd",
		Short: "Decentralized Sequencer for StaionApps",
	}

	rootCmd.AddCommand(command.StationCmd)
	rootCmd.AddCommand(command.InitCmd)
	rootCmd.AddCommand(command.KeyGenCmd)
	rootCmd.AddCommand(command.ProverGenCMD)
	command.KeyGenCmd.AddCommand(keys.JunctionKeyGenCmd)
	command.ProverGenCMD.AddCommand(zkpCmd.V1ZKP)
	keys.JunctionKeyGenCmd.Flags().StringVarP(&keys.AcountName, "accountName", "n", "", "Account Name")
	keys.JunctionKeyGenCmd.Flags().StringVarP(&keys.AccountPath, "accountPath", "p", "", "Account Path")
	keys.JunctionKeyGenCmd.MarkFlagRequired("accountName")
	keys.JunctionKeyGenCmd.MarkFlagRequired("accountPath")
	command.InitCmd.Flags().String("moniker", "", "Moniker for the sequencer")
	command.InitCmd.Flags().String("stationType", "", "Station Type for the sequencer (evm | cosmwasm | svm)")
	command.InitCmd.Flags().String("daType", "mock", "DA Type for the sequencer (avail | celestia | eigen | mock)")
	command.InitCmd.MarkFlagRequired("moniker")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
