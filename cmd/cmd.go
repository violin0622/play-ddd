package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	root := &cobra.Command{
		Use:     `play-ddd`,
		Short:   `Play DDD Example Program.`,
		Long:    slogan,
		Version: fmt.Sprintf(versionTpl, Version, GitSHA, BuiltAt, GoVersion),
	}

	root.AddCommand(Contents())
	return root
}

const (
	slogan = `
     ____  __               ____  ____  ____
    / __ \/ /___ ___  __   / __ \/ __ \/ __ \
   / /_/ / / __ '/ / / /  / / / / / / / / / /
  / ____/ / /_/ / /_/ /  / /_/ / /_/ / /_/ /
 /_/   /_/\__,_/\__, /  /_____/_____/_____/
               /____/

Play DDD Example Program.
`

	versionTpl = `

Version: %s
Git SHA: %s
Go Version: %s
Built At: %s
`
)

var Version, GitSHA, BuiltAt, GoVersion string
