package cmd

import (
	"github.com/spf13/cobra"

	"play-ddd/cmd/contents"
)

func Contents() *cobra.Command {
	cmd := &cobra.Command{
		Use:     `contents`,
		Aliases: []string{`cont`},
		Short:   `contents service actions`,
	}

	cmd.AddCommand(contents.Run())
	cmd.AddCommand(contents.Config())

	return cmd
}
