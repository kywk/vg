// Copyright © 2017 Stream
//

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/GetStream/vg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// initSettingsCmd represents the initSettings command
var initSettingsCmd = &cobra.Command{
	Use:   "initSettings [workspaceName]",
	Short: "This command initializes the settings file for a certain workspace",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too much arguments specified")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		workspace := ""
		if len(args) == 1 {
			workspace = args[0]
		} else {
			workspace, err = os.Getwd()
			if err != nil {
				return errors.WithStack(err)
			}
			workspace = filepath.Base(workspace)

		}
		fmt.Println(workspace)

		path := utils.SettingsPath(workspace)
		if err != nil {
			return errors.WithStack(err)
		}

		dir := filepath.Dir(path)

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return errors.WithStack(err)
		}

		// Check if it's a new workspace. Only continue if this is the case or
		// if force is set.
		_, err = os.Stat(dir)
		if err != nil {
			if !os.IsNotExist(err) {
				return errors.WithStack(err)
			}
		} else if !force {
			return nil
		}

		settings := utils.WorkspaceSettings{}
		settings.GlobalFallback, err = cmd.Flags().GetBool("global-fallback")
		if err != nil {
			return errors.WithStack(err)
		}

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return errors.WithStack(err)
		}

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return errors.WithStack(err)
		}

		err = toml.NewEncoder(file).Encode(settings)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(initSettingsCmd)
	initSettingsCmd.PersistentFlags().Bool("global-fallback", false, "Fallback to global packages when they are not present in workspace")
	initSettingsCmd.PersistentFlags().BoolP("force", "f", false, "")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initSettingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initSettingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
