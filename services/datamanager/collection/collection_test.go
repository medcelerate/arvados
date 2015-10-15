package collection

import (
	"git.curoverse.com/arvados.git/sdk/go/blockdigest"
	. "gopkg.in/check.v1"
	"testing"
)

// Gocheck boilerplate
func Test(t *testing.T) {
	TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})

// This captures the result we expect from
// ReadCollections.Summarize().  Because CollectionUUIDToIndex is
// indeterminate, we replace BlockToCollectionIndices with
// BlockToCollectionUuids.
type ExpectedSummary struct {
	OwnerToCollectionSize     map[string]int
	BlockToDesiredReplication map[blockdigest.DigestWithSize]int
	BlockToCollectionUuids    map[blockdigest.DigestWithSize][]string
}

func CompareSummarizedReadCollections(c *C,
	summarized ReadCollections,
	expected ExpectedSummary) {

	c.Assert(summarized.OwnerToCollectionSize, DeepEquals,
		expected.OwnerToCollectionSize)

	c.Assert(summarized.BlockToDesiredReplication, DeepEquals,
		expected.BlockToDesiredReplication)

	summarizedBlockToCollectionUuids :=
		make(map[blockdigest.DigestWithSize]map[string]struct{})
	for digest, indices := range summarized.BlockToCollectionIndices {
		uuidSet := make(map[string]struct{})
		summarizedBlockToCollectionUuids[digest] = uuidSet
		for _, index := range indices {
			uuidSet[summarized.CollectionIndexToUUID[index]] = struct{}{}
		}
	}

	expectedBlockToCollectionUuids :=
		make(map[blockdigest.DigestWithSize]map[string]struct{})
	for digest, uuidSlice := range expected.BlockToCollectionUuids {
		uuidSet := make(map[string]struct{})
		expectedBlockToCollectionUuids[digest] = uuidSet
		for _, uuid := range uuidSlice {
			uuidSet[uuid] = struct{}{}
		}
	}

	c.Assert(summarizedBlockToCollectionUuids, DeepEquals,
		expectedBlockToCollectionUuids)
}

func (s *MySuite) TestSummarizeSimple(checker *C) {
	rc := MakeTestReadCollections([]TestCollectionSpec{TestCollectionSpec{
		ReplicationLevel: 5,
		Blocks:           []int{1, 2},
	}})

	rc.Summarize(nil)

	c := rc.UUIDToCollection["col0"]

	blockDigest1 := blockdigest.MakeTestDigestWithSize(1)
	blockDigest2 := blockdigest.MakeTestDigestWithSize(2)

	expected := ExpectedSummary{
		OwnerToCollectionSize:     map[string]int{c.OwnerUUID: c.TotalSize},
		BlockToDesiredReplication: map[blockdigest.DigestWithSize]int{blockDigest1: 5, blockDigest2: 5},
		BlockToCollectionUuids:    map[blockdigest.DigestWithSize][]string{blockDigest1: []string{c.UUID}, blockDigest2: []string{c.UUID}},
	}

	CompareSummarizedReadCollections(checker, rc, expected)
}

func (s *MySuite) TestSummarizeOverlapping(checker *C) {
	rc := MakeTestReadCollections([]TestCollectionSpec{
		TestCollectionSpec{
			ReplicationLevel: 5,
			Blocks:           []int{1, 2},
		},
		TestCollectionSpec{
			ReplicationLevel: 8,
			Blocks:           []int{2, 3},
		},
	})

	rc.Summarize(nil)

	c0 := rc.UUIDToCollection["col0"]
	c1 := rc.UUIDToCollection["col1"]

	blockDigest1 := blockdigest.MakeTestDigestWithSize(1)
	blockDigest2 := blockdigest.MakeTestDigestWithSize(2)
	blockDigest3 := blockdigest.MakeTestDigestWithSize(3)

	expected := ExpectedSummary{
		OwnerToCollectionSize: map[string]int{
			c0.OwnerUUID: c0.TotalSize,
			c1.OwnerUUID: c1.TotalSize,
		},
		BlockToDesiredReplication: map[blockdigest.DigestWithSize]int{
			blockDigest1: 5,
			blockDigest2: 8,
			blockDigest3: 8,
		},
		BlockToCollectionUuids: map[blockdigest.DigestWithSize][]string{
			blockDigest1: []string{c0.UUID},
			blockDigest2: []string{c0.UUID, c1.UUID},
			blockDigest3: []string{c1.UUID},
		},
	}

	CompareSummarizedReadCollections(checker, rc, expected)
}
