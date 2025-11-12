package stats

import (
	"context"

	"google.golang.org/grpc/stats"
)

type rpcKey struct{}

var _ stats.Handler = (*payloadSize)(nil)

type PayloadSize struct {
	InBytes  int64
	OutBytes int64
}

func NewInjectPayloadSize() *payloadSize {
	return &payloadSize{}
}

type payloadSize struct{}

func FromContext(ctx context.Context) (*PayloadSize, bool) {
	ps, ok := ctx.Value(rpcKey{}).(*PayloadSize)
	return ps, ok
}

// HandleConn implements stats.Handler.
func (*payloadSize) HandleConn(context.Context, stats.ConnStats) {}

// TagConn implements stats.Handler.
func (*payloadSize) TagConn(
	ctx context.Context,
	_ *stats.ConnTagInfo,
) context.Context {
	return ctx
}

// HandleRPC implements stats.Handler.
func (*payloadSize) HandleRPC(ctx context.Context, info stats.RPCStats) {
	ps, ok := ctx.Value(rpcKey{}).(*PayloadSize)
	if !ok {
		return
	}

	switch info := info.(type) {
	case *stats.InPayload:
		ps.InBytes += int64(info.CompressedLength)
	case *stats.OutPayload:
		ps.InBytes += int64(info.CompressedLength)
	}
}

// TagRPC implements stats.Handler.
func (*payloadSize) TagRPC(
	ctx context.Context, info *stats.RPCTagInfo,
) context.Context {
	return context.WithValue(ctx, rpcKey{}, &PayloadSize{})
}
