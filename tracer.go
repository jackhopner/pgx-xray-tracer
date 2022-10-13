package pgxxray

import (
	"context"
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jackc/pgx/v5"
)

type segmentContextKey string

const xraySegmentKey segmentContextKey = "xray_segment"

type PGXTracer struct {
}

// TraceQueryStart is called at the beginning of Query, QueryRow, and Exec calls. The returned context is used for the
// rest of the call and will be passed to TraceQueryEnd.
func (t *PGXTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx, seg := xray.BeginSubsegment(ctx, segmentName(conn))
	seg.AddMetadata("sql", data.SQL)
	seg.AddMetadata("sql_args", data.Args)

	return context.WithValue(ctx, xraySegmentKey, seg)
}

func (t *PGXTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	segRaw := ctx.Value(xraySegmentKey)
	if segRaw != nil {
		seg, ok := segRaw.(*xray.Segment)
		if ok {
			seg.AddMetadata("sql_rows_affected", data.CommandTag.RowsAffected())
			seg.Close(data.Err)
		}
	}
}

func segmentName(conn *pgx.Conn) string {
	return fmt.Sprintf("%s/%s", conn.Config().Host, conn.Config().Database)
}
