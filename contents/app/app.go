package app

import (
	"play-ddd/contents/domain/novel"
)

type service struct {
	novelRepo novel.Repo
	novelFact novel.Factory
}

func New(
	novelRepo novel.Repo,
	novelFact novel.Factory,
) service {
	return service{
		novelRepo: novelRepo,
		novelFact: novelFact,
	}
}
