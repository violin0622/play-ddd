package outbox_test

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

const raw = `{"@type":"type.googleapis.com/contents.novel.v1.NovelCreated","authorId":{"ulid":"01KD2Q3EVM65QRM1M28S0ZM7C6"},"title":"My Book Title","description":"My book is good!","tags":["Tokyo","Hot"],"category":"Story"}`

func TestUnmarshal(t *testing.T) {
	var aa anypb.Any
	err := protojson.Unmarshal([]byte(raw), &aa)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(&aa)
	t.Fail()
}
