package pgxxray

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jackc/pgx/v5"
)

func (t *PGXTracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	if t.traceEnabled[BatchTraceType] {
		var seg *xray.Segment
		ctx, seg = t.beginSubsegment(ctx, conn.Config(), "BATCH")
		seg.AddMetadata("sql_batch_length", data.Batch.Len())
	}

	return ctx
}

func (t *PGXTracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
	if t.traceEnabled[BatchTraceType] {
		seg := t.tryGetSegment(ctx)
		if seg != nil {
			addSegmentMetadataString(seg, "sql", data.SQL)
			addSegmentMetadataArray(seg, "sql_args", data.Args)
			seg.AddMetadata("sql_rows_affected", data.CommandTag.RowsAffected())
			seg.Close(data.Err)
		}
	}
}

func (t *PGXTracer) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
	if t.traceEnabled[BatchTraceType] {
		seg := t.tryGetSegment(ctx)
		if seg != nil {
			seg.Close(data.Err)
		}
	}
}
