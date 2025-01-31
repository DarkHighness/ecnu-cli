package cmd

import (
	"fmt"
	"os"

	"github.com/DarkHighness/ecnu-cli/cmd/auth"
	"github.com/DarkHighness/ecnu-cli/cmd/net"
	"github.com/DarkHighness/ecnu-cli/cmd/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ecnu-cli",
	Short:   "hdu cli",
	Version: "alpha",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("save") {
			if viper.WriteConfig() != nil {
				viper.Set("verbose", nil)
				viper.Set("save", nil)
				cobra.CheckErr(viper.SafeWriteConfig())
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ecnu-cli.yaml)")
	rootCmd.PersistentFlags().BoolP("save", "s", false, "save config")
	cobra.CheckErr(viper.BindPFlag("save", rootCmd.PersistentFlags().Lookup("save")))
	rootCmd.PersistentFlags().BoolP("verbose", "V", false, "show more info")
	cobra.CheckErr(viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")))

	rootCmd.AddCommand(net.Cmd, auth.Cmd, rpc.Cmd)
}

var cfgFile string

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ecnu-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ecnu-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
