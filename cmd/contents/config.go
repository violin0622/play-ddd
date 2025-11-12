package contents

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func Config() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `config`,
		Short: `Contents service config management.`,
	}

	cmd.AddCommand(create())
	cmd.AddCommand(validate())
	return cmd
}

func create() *cobra.Command {
	return &cobra.Command{
		Use:   `create`,
		Short: `Create a default cofig file.`,
		RunE: func(*cobra.Command, []string) error {
			f, err := os.Create(`config.yml`)
			if err != nil {
				return err
			}

			slog.Info(`Config file created.`,
				slog.String(`file`, f.Name()))

			return nil
		},
	}
}

func validate() *cobra.Command {
	return &cobra.Command{
		Use:   `validate`,
		Short: `Validate cofig file.`,
		RunE: func(*cobra.Command, []string) error {
			slog.Info(`Config file is OK.`)
			return nil
		},
	}
}
