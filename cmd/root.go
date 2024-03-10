package cmd

import (
	"fmt"
	"os"

	"github.com/maetthu/restic-wrap/lib/profile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var p string
var Profile profile.Profile

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&p, "profile", "p", "", "Path to profile.yaml")
	rootCmd.PersistentFlags().StringP("backend", "b", "", "Backend to use (depending on the command, either the first one or all are used by default)")

	rootCmd.MarkFlagFilename("profile")
	rootCmd.MarkFlagRequired("profile")
}

func initConfig() {
	viper.SetConfigFile(p)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	err := viper.Unmarshal(&Profile)

	if err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
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
