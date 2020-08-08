package sneak

import (
	"fmt"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version is a value injected at compile time for the current version of sneak
var Version = "master"

// local
var (
	// sneakCfg are the struct representation of sneak settings
	sneakCfg config.Settings
	// cfgfile is the path to the sneak config file different than the default
	cfgFile string
	// dataDir is the path where sneak stores its key/value database
	dataDir string
)

var rootCmd = &cobra.Command{
	Use:   "sneak",
	Short: "a tool for common actions when pentesting/playing CTFs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Banner)

		err := config.SafeWriteConfig()
		if err != nil {
			gui.ExitWithError(err)
		}

		sneakCfg, err = config.Parse(viper.GetViper())
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

// Execute adds all child commands to the root command set sets flags appropriately
func Execute() {
	cobra.OnInitialize(config.InitConfig)
	initGlobalFlags()

	if err := rootCmd.Execute(); err != nil {
		gui.ExitWithError(err)
	}
}

// -------------------- init

func initGlobalFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sneak/.sneak.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	rootCmd.PersistentFlags().StringVar(&dataDir, "data", "", "database dir default is $HOME/.sneak")
	viper.BindPFlag("data", rootCmd.PersistentFlags().Lookup("data"))
}

func init() {
	rootCmd.AddCommand(configCmd)
}
