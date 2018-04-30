package simple

import (
	. "../ctx"
	. "../mid"
	"errors"
	"strings"
)

// sequence an array of atomids
func SeqAtomMids(mids ...interface{}) func(*ApiContext) {
	var midList []AtomMid
	for _, mid := range mids {
		switch item := mid.(type) {
		case AtomMid:
			midList = append(midList, item)
		case string:
			//
			midList = append(midList, parseAtomMidText(item)...)
		default:
			panic(errors.New("Expect type AtomMid or string"))
		}
	}
	return seqAtomMids(midList, 0, len(midList)-1)
}

func parseAtomMidText(text string) []AtomMid {
	var mids []AtomMid
	parts := strings.Split(text, "\\")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			mids = append(mids, SimpleAtomMid(part))
		}
	}

	return mids
}

func seqAtomMids(mids []AtomMid, start int, end int) func(*ApiContext) {
	if start > end {
		return NoopMid
	}

	mid := mids[start]
	// recursive rest
	rest := seqAtomMids(mids, start+1, end)

	return func(ctx *ApiContext) {
		mid(ctx, func() {
			rest(ctx)
		})
	}
}

func NoopMid(ctx *ApiContext) {}

// TODO branch atom mids
