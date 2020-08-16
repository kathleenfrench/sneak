package sneak

import (
	"fmt"
	"time"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/kathleenfrench/sneak/internal/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

// Version is a value injected at compile time for the current version of sneak
var Version = "master"

// mountData is a bool flag used when running sneaker with mounted local .sneak configs
var mountData bool

// unMountData is a bool flag used when converting containerized configs back into compatibility with your local fs
var unMountData bool

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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := config.SafeWriteConfig(mountData, unMountData)
		if err != nil {
			gui.ExitWithError(err)
		}

		sneakCfg, err = config.Parse(viper.GetViper())
		if err != nil {
			gui.ExitWithError(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if mountData || unMountData {
			gui.Info("+1", "your configs have been updated!", nil)
			return
		}

		fmt.Println(config.Banner)
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command set sets flags appropriately
func Execute() {
	var err error

	dataDir, err = store.GetDataDirectory()
	if err != nil {
		gui.ExitWithError(err)
	}

	opts := &bolthold.Options{
		Options: &bbolt.Options{
			Timeout: 10 * time.Second,
		},
	}

	db, err = bolthold.Open(fmt.Sprintf("%s/sneak.db", dataDir), 0600, opts)
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
	rootCmd.PersistentFlags().BoolVarP(&mountData, "mount", "m", false, "used when running the full sneaker containerized environment and mounting local .sneak config files - only needs to be run once")
	// viper.BindPFlag("mounted_data", rootCmd.PersistentFlags().Lookup("mount"))
	rootCmd.PersistentFlags().BoolVarP(&unMountData, "unmount", "u", false, "used when running sneak locally after having mounted persistent data into a containerized environment to refresh config path values - only needs to be run once")
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	initGlobalFlags()
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(boxSubCmd)
	rootCmd.AddCommand(gotoCmd)
	rootCmd.AddCommand(vpnCmd)
	rootCmd.AddCommand(pipelineCmd)
}
