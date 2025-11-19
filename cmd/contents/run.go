package contents

import (
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"play-ddd/contents/app"
	"play-ddd/contents/iface/grpc"
	"play-ddd/contents/infra/repository/pg"
	"play-ddd/contents/infra/server"
)

func Run() *cobra.Command {
	return &cobra.Command{
		Use:   `run`,
		Short: `start contents service`,
		RunE:  runErr,
	}
}

func runErr(*cobra.Command, []string) error {
	db, err := pg.InitDB(pg.DSN)
	if err != nil {
		return err
	}

	repo := pg.New(db)
	ch := app.NewCommandHandler(repo, logr.Discard())
	cmdSvc := grpc.NewCmdService(ch)
	qh := app.NewQueryHandler()
	querySvc := grpc.NewQueryService(qh)
	server.New(
		&cmdSvc,
		&querySvc,
	)

	return nil
}
