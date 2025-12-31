package contents

import (
	"github.com/spf13/cobra"
)

func Contents() *cobra.Command {
	cmd := &cobra.Command{
		Use:     `contents`,
		Aliases: []string{`cont`},
		Short:   `contents service actions`,
	}

	cmd.AddCommand(Run())
	cmd.AddCommand(Config())

	return cmd
}
