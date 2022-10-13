package pgxxray

import (
	"context"
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jackc/pgx/v5"
)

type segmentContextKey string

const xraySegmentKey segmentContextKey = "xray_segment"

type TraceType string

const (
	BatchTraceType    TraceType = "batch"
	ConnectTraceType  TraceType = "connect"
	CopyFromTraceType TraceType = "copy_from"
	PrepareTraceType  TraceType = "prepare"
	QueryTraceType    TraceType = "query"
)

type PGXTracer struct {
	traceEnabled map[TraceType]bool
}

func (t *PGXTracer) beginSubsegment(ctx context.Context, cfg *pgx.ConnConfig, prefix string) (context.Context, *xray.Segment) {
	ctx, seg := xray.BeginSubsegment(ctx, t.segmentName(cfg, prefix))

	return context.WithValue(ctx, xraySegmentKey, seg), seg
}

func (t *PGXTracer) segmentName(cfg *pgx.ConnConfig, prefix string) string {
	return fmt.Sprintf("%s-%s/%s", prefix, cfg.Host, cfg.Database)
}

func (t *PGXTracer) tryGetSegment(ctx context.Context) *xray.Segment {
	segRaw := ctx.Value(xraySegmentKey)
	if segRaw != nil {
		seg, ok := segRaw.(*xray.Segment)
		if ok {
			return seg
		}
	}

	return nil
}

func NewPGXTracer(traceTypes ...TraceType) *PGXTracer {
	traceEnabled := map[TraceType]bool{}
	if len(traceTypes) == 0 {
		traceEnabled[BatchTraceType] = true
		traceEnabled[ConnectTraceType] = true
		traceEnabled[CopyFromTraceType] = true
		traceEnabled[PrepareTraceType] = true
		traceEnabled[QueryTraceType] = true
	} else {
		for _, typ := range traceTypes {
			traceEnabled[typ] = true
		}
	}

	return &PGXTracer{
		traceEnabled: traceEnabled,
	}
}
