package pgxxray

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jackc/pgx/v5"
)

func (t *PGXTracer) TraceCopyFromStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	if t.traceEnabled[CopyFromTraceType] && t.hasSegment(ctx) {
		var seg *xray.Segment
		ctx, seg = t.beginSubsegment(ctx, conn.Config(), "COPY")
		seg.AddMetadata("sql_column_names", data.ColumnNames)
		seg.AddMetadata("sql_table_names", data.TableName)
	}
	return ctx
}

func (t *PGXTracer) TraceCopyFromEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceCopyFromEndData) {
	if t.traceEnabled[CopyFromTraceType] && t.hasSegment(ctx) {
		seg := t.tryGetSegment(ctx)
		if seg != nil {
			seg.AddMetadata("sql_rows_affected", data.CommandTag.RowsAffected())
			seg.Close(data.Err)
		}
	}
}
