package pgxxray

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (t *PGXTracer) TraceConnectStart(ctx context.Context, data pgx.TraceConnectStartData) context.Context {
	if t.traceEnabled[ConnectTraceType] {
		ctx, _ = t.beginSubsegment(ctx, data.ConnConfig, "CONNECT")
	}

	return ctx
}

func (t *PGXTracer) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	if t.traceEnabled[ConnectTraceType] {
		seg := t.tryGetSegment(ctx)
		if seg != nil {
			seg.Close(data.Err)
		}
	}
}
