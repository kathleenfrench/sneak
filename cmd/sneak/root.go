package sneak

import (
	"fmt"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/kathleenfrench/sneak/internal/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/timshannon/bolthold"
)

// Version is a value injected at compile time for the current version of sneak
var Version = "master"

var testing string

// local
var (
	// sneakCfg are the struct representation of sneak settings
	sneakCfg *config.Settings
	// cfgfile is the path to the sneak config file different than the default
	cfgFile string
	// dataDir is the path where sneak stores its key/value database
	dataDir string

	db *bolthold.Store
)

var rootCmd = &cobra.Command{
	Use:   "sneak",
	Short: "a tool for common actions when pentesting/playing CTFs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Banner)
		cmd.Usage()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
	var err error

	dataDir, err = store.GetDataDirectory()
	if err != nil {
		gui.ExitWithError(err)
	}

	db, err = bolthold.Open(fmt.Sprintf("%s/sneak.db", dataDir), 0600, nil)
	if err != nil {
		gui.ExitWithError(fmt.Sprintf("could not initialize database - %s", err))
	}

	defer db.Close()

	if err = rootCmd.Execute(); err != nil {
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
	initGlobalFlags()
	cobra.OnInitialize(config.InitConfig)

	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUpdateCmd)
	configCmd.AddCommand(configDelCmd)

	rootCmd.AddCommand(configCmd)
}
