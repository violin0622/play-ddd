package events_test

import (
	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	ne "play-ddd/contents/domain/novel/events"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

var _ = Describe(`NovelCreatedV2 event fromPB`, func() {
	var (
		payload, _ = anypb.New(&novelv1.NovelCreated{
			AuthorId:    ulidpb.From(ulid.Make()),
			Title:       `My Book Title`,
			Description: `My book is good!`,
			Tags:        []string{`Tokyo`, `Hot`},
			Category:    `Story`,
		})
		pbEvent = novelv1.Event{
			Id:            ulidpb.From(ulid.Make()),
			AggregateId:   ulidpb.From(ulid.Make()),
			EmitAt:        timestamppb.Now(),
			Kind:          "created",
			AggregateKind: "novel",
			Payload:       payload,
		}
	)

	It(`can unmarshal from proto`, func() {
		var e ne.NovelCreatedV2
		Ω(e.FromPB(&pbEvent)).Should(Succeed())
		Ω(e.ID()).Should(BeEquivalentTo(pbEvent.GetId().Into()))

		raw, _ := protojson.Marshal(pbEvent.GetPayload())
		Ω(e.Payload()).Should(Equal(raw))
	})
})
