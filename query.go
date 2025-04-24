package pgxxray

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jackc/pgx/v5"
)

func (t *PGXTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if t.traceEnabled[QueryTraceType] && t.hasSegment(ctx) {
		var seg *xray.Segment
		ctx, seg = t.beginSubsegment(ctx, conn.Config(), "QUERY")
		addSegmentMetadataString(seg, "sql", data.SQL)
		addSegmentMetadataArray(seg, "sql_args", data.Args)
	}
	return ctx
}

func (t *PGXTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if t.traceEnabled[QueryTraceType] && t.hasSegment(ctx) {
		seg := t.tryGetSegment(ctx)
		if seg != nil {
			seg.AddMetadata("sql_rows_affected", data.CommandTag.RowsAffected())
			seg.Close(data.Err)
		}
	}
}
