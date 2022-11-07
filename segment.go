package pgxxray

import (
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xray"
)

const arrayLengthLimit = 30
const stringLengthLimit = 1000

func addSegmentMetadataArray(seg *xray.Segment, key string, value []any) {
	seg.AddMetadata(fmt.Sprintf("%s_length", key), len(value))
	if len(value) >= arrayLengthLimit {
		seg.AddMetadata(key, append(value[0:arrayLengthLimit-1], "..."))
		return
	}

	seg.AddMetadata(key, value)
}

func addSegmentMetadataString(seg *xray.Segment, key, value string) {
	if len(value) >= stringLengthLimit {
		seg.AddMetadata(key, fmt.Sprintf("%s...", value[0:stringLengthLimit-3]))
		return
	}

	seg.AddMetadata(key, value)
}
