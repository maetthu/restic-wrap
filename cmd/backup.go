package cmd

import (
	"os/exec"

	"github.com/maetthu/restic-wrap/lib/profile"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	rootCmd.AddCommand(backupCmd)
}

type LogWriter struct {
	Logger *zap.SugaredLogger
	Level  zapcore.Level
	Fields []string
}

func (w LogWriter) Write(p []byte) (int, error) {
	w.Logger.Log(w.Level, string(p))
	return len(p), nil
}

var backupCmd = &cobra.Command{
	Use:   "backup -p <path-to-profile.yaml> [flags]",
	Short: "Executes configured backup stages for all backends",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := zap.NewProductionConfig()
		config.DisableStacktrace = true
		logger, _ := config.Build()
		defer logger.Sync()

		sugar := logger.Sugar().WithOptions()

		logWriterInfo := LogWriter{
			Logger: sugar,
			Level:  zapcore.InfoLevel,
			Fields: []string{},
		}

		logWriterError := LogWriter{
			Logger: sugar,
			Level:  zapcore.ErrorLevel,
			Fields: []string{},
		}

		run := func(b *profile.Backend, args []string) error {
			r := exec.Command("restic", args...)
			r.Env = Profile.BuildEnv(b)
			r.Stdout = logWriterInfo
			r.Stderr = logWriterError

			err := r.Run()
			return err
		}

		notify := func(b *profile.Backend, stage string, level string, msg string) error {
			for _, command := range Profile.Notify {
				n := exec.Command(command, b.Name, stage, level, msg)
				n.Env = Profile.BuildEnv(b)
				n.Stdout = logWriterInfo
				n.Stderr = logWriterError

				err := n.Run()

				if err != nil {
					return err
				}
			}

			return nil
		}

		backends := Profile.Backends

		if name, err := cmd.Flags().GetString("backend"); err != nil {
			b, err := Profile.Backend(name)

			if err != nil {
				return err
			}

			backends = []*profile.Backend{b}
		}

		for _, b := range backends {
			sugar.Infow("Start backup", "backend", b.Name)

			for _, s := range Profile.Stages {
				sugar.Infow("Start backup stage", "backend", b.Name, "stage", s.Command)
				logWriterInfo.Fields = []string{"backend", b.Name, "stage", s.Command}
				logWriterError.Fields = logWriterInfo.Fields

				args := []string{s.Command}
				args = append(args, s.Args...)

				err := run(b, args)

				if err != nil {
					sugar.Errorw("Failed backup stage", "error", err.Error(), "backend", b.Name, "stage", s.Command)
					err = notify(b, s.Command, "error", err.Error())

					if err != nil {
						sugar.Warnw("Failed invoking notification command", "error", err.Error(), "backend", b.Name, "stage", s.Command)
					}

					break
				}

				sugar.Infow("Finished backup stage", "backend", b.Name, "stage", s.Command)
				err = notify(b, s.Command, "success", "")

				if err != nil {
					sugar.Warnw("Failed invoking notification command", "error", err.Error(), "backend", b.Name, "stage", s.Command)
				}
			}

			sugar.Infow("Finished backup", "backend", b.Name)
		}

		return nil
	},
}
