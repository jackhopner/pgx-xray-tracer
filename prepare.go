package pgxxray

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jackc/pgx/v5"
)

func (t *PGXTracer) TracePrepareStart(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareStartData) context.Context {
	if t.traceEnabled[PrepareTraceType] && t.hasSegment(ctx) {
		var seg *xray.Segment
		ctx, seg = t.beginSubsegment(ctx, conn.Config(), "PREPARE")
		seg.AddMetadata("sql_name", data.Name)
		addSegmentMetadataString(seg, "sql", data.SQL)
	}

	return ctx
}

func (t *PGXTracer) TracePrepareEnd(ctx context.Context, conn *pgx.Conn, data pgx.TracePrepareEndData) {
	if t.traceEnabled[PrepareTraceType] && t.hasSegment(ctx) {
		seg := t.tryGetSegment(ctx)
		if seg != nil {
			seg.AddMetadata("sql_already_prepared", data.AlreadyPrepared)
			seg.Close(data.Err)
		}
	}
}
