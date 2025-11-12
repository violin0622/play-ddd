package contents

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func Run() *cobra.Command {
	return &cobra.Command{
		Use:   `run`,
		Short: `start contents service`,

		RunE: func(*cobra.Command, []string) error {
			slog.Info(`Hello world!`)
			defer slog.Info(`Good bye!`)

			return nil
		},
	}
}
