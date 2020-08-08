package sneak

import (
	"fmt"

	"github.com/kathleenfrench/common/gui"
	"github.com/kathleenfrench/sneak/internal/helpers"
	"github.com/kathleenfrench/sneak/internal/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resetBucket string

var dbCmd = &cobra.Command{
	Use:    "db",
	Short:  "interact with the sneak database",
	Hidden: true,
}

var dbViewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"connect", "see", "show"},
	PreRun: func(cmd *cobra.Command, args []string) {
		err := db.Close()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := fmt.Sprintf("%s/sneak.db", viper.GetString("data"))
		err := store.Audit(dbPath)
		if err != nil {
			gui.ExitWithError(err)
		}
	},
}

var dbResetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"nuke", "clear"},
	Run: func(cmd *cobra.Command, args []string) {
		helpers.Spacer()
		if resetBucket != "" {
			store.EmptyBuckets(db, resetBucket)
			gui.Info("fire", "bucket emptied", resetBucket)
		} else {
			store.EmptyBuckets(db, "/all")
			gui.Info("fire", "your local DB has been wiped", "bolt")
		}
	},
}

var dbBackupCmd = &cobra.Command{
	Use:     "backup",
	Aliases: []string{"bu"},
	PreRun: func(cmd *cobra.Command, args []string) {
		err := db.Close()
		if err != nil {
			gui.ExitWithError(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		helpers.Spacer()
		dbdir := viper.GetString("data")

		gui.Info("popcorn", "backing up your local DB", fmt.Sprintf("%s/sneak_backup.db", dbdir))

		err := store.Backup(dbdir)
		if err != nil {
			gui.ExitWithError(err)
		}

		gui.Info("+1", "local db backed up!", fmt.Sprintf("%s/sneak_backup.db", dbdir))

	},
}

func init() {
	dbResetCmd.Flags().StringVarP(&resetBucket, "bucket", "b", "", "reset a specific bucket")
	dbCmd.AddCommand(dbViewCmd)
	dbCmd.AddCommand(dbResetCmd)
	dbCmd.AddCommand(dbBackupCmd)
}
