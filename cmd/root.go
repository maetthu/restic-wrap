package cmd

import (
	"fmt"
	"os"

	"github.com/maetthu/restic-wrap/lib/profile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.PersistentFlags().StringP("profile", "p", "", "Path to profile.yaml")
	rootCmd.PersistentFlags().StringP("backend", "b", "", "Backend to use (depending on the command, either the first one or all are used by default)")

	_ = rootCmd.MarkPersistentFlagFilename("profile")
	_ = rootCmd.MarkPersistentFlagRequired("profile")
}

func initConfig(profilePath string) (profile.Profile, error) {
	viper.SetConfigFile(profilePath)

	if err := viper.ReadInConfig(); err != nil {
		return profile.Profile{}, err
	}

	p := profile.Profile{}
	err := viper.Unmarshal(&p)

	if err != nil {
		return profile.Profile{}, err
	}

	return p, nil
}

var rootCmd = &cobra.Command{
	Use:     "restic-wrap",
	Version: version,
	Short:   "Restic wrapper tool with profile support",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
