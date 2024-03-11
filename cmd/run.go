package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run -p <path-to-profile.yaml> [flags] <restic command and flags>",
	Short: "Run adhoc restic command with all the necessary environment variables set for a specific backend",
	RunE: func(cmd *cobra.Command, args []string) error {
		profilePath, err := cmd.Flags().GetString("profile")

		if err != nil {
			return err
		}

		prof, err := initConfig(profilePath)

		if err != nil {
			return err
		}

		restic := exec.Command("restic", args...)
		backend := prof.Backends[0]

		if b, err := cmd.Flags().GetString("backend"); err == nil && b != "" {
			backend, err = prof.Backend(b)

			if err != nil {
				return err
			}
		}

		restic.Env = prof.BuildEnv(backend)
		restic.Stdout = os.Stdout
		restic.Stderr = os.Stderr

		return restic.Run()
	},
}
