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
		restic := exec.Command("restic", args...)
		backend := Profile.Backends[0]

		if b, err := cmd.Flags().GetString("backend"); err != nil {
			backend, err = Profile.Backend(b)

			if err != nil {
				return err
			}
		}

		restic.Env = Profile.BuildEnv(backend)
		restic.Stdout = os.Stdout
		restic.Stderr = os.Stderr

		return restic.Run()
	},
}
