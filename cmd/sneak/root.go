package sneak

import (
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
)

var rootCmd = &cobra.Command{
	Use:   "sneak",
	Short: "a tool for common actions when pentesting/playing CTFs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		s, err := config.Parse(viper.GetViper())
		if err != nil {
			gui.ExitWithError(err)
		}

		sneakCfg = s
	},
}

// Execute adds all child commands to the root command set sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		gui.ExitWithError(err)
	}
}

// -------------------- init

func initGlobalFlags() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/sneak/config.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
}

func init() {
	cobra.OnInitialize(config.Initialize)
	initGlobalFlags()

	rootCmd.AddCommand(configCmd)
}
