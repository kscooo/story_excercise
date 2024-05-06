package blockprocessor

import (
	"reflect"
	"testing"
)

func createZeroPeers(n int) []Peer {
	peers := make([]Peer, n)
	for i := range peers {
		peers[i] = &FakePeer{Blocks: []Block{}, Height: 0}
	}
	return peers
}

type FakePeer struct {
	Blocks []Block
	Height uint64
}

func (f *FakePeer) GetBlocks() ([]Block, uint64) {
	return f.Blocks, f.Height
}

func TestBlockProcessor_Initialize(t *testing.T) {
	testCases := []struct {
		peers        []Peer
		expectBlocks []Block
		expectHeight uint64
		expectError  bool
	}{
		{
			peers: []Peer{
				&FakePeer{Blocks: []Block{"A", "B", "C"}, Height: 3},
				&FakePeer{Blocks: []Block{"A", "B", "C"}, Height: 3},
				&FakePeer{Blocks: []Block{"A", "B", "C"}, Height: 3},
				&FakePeer{Blocks: []Block{"A", "B", "C"}, Height: 3},
				&FakePeer{Blocks: []Block{"A", "B", "C"}, Height: 3},
			},
			expectBlocks: []Block{"A", "B", "C"},
			expectHeight: 3,
			expectError:  false,
		},
		{
			peers: []Peer{
				&FakePeer{Blocks: []Block{"A"}, Height: 1},
				&FakePeer{Blocks: []Block{"B"}, Height: 1},
				&FakePeer{Blocks: []Block{"C"}, Height: 1},
				&FakePeer{Blocks: []Block{"D"}, Height: 1},
				&FakePeer{Blocks: []Block{"E"}, Height: 1},
			},
			expectBlocks: nil,
			expectHeight: 0,
			expectError:  false,
		},
		{
			peers: []Peer{
				&FakePeer{Blocks: []Block{"X", "Y"}, Height: 2},
				&FakePeer{Blocks: []Block{"X", "Y"}, Height: 2},
				&FakePeer{Blocks: []Block{"X", "Y"}, Height: 2},
				&FakePeer{Blocks: []Block{"A", "B"}, Height: 2},
				&FakePeer{Blocks: []Block{"A", "B"}, Height: 2},
			},
			expectBlocks: []Block{"X", "Y"},
			expectHeight: 2,
			expectError:  false,
		},
		{
			peers: []Peer{
				&FakePeer{Blocks: []Block{"A", "B", "C"}, Height: 3},
				&FakePeer{Blocks: []Block{"A", "B", "C"}, Height: 3},
				&FakePeer{Blocks: []Block{"X", "Y"}, Height: 0},
				&FakePeer{Blocks: []Block{"X", "Y"}, Height: 0},
				&FakePeer{Blocks: []Block{"X", "Y"}, Height: 0},
			},
			expectBlocks: []Block{"X", "Y"},
			expectHeight: 0,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		bp := NewBlockProcessor(tc.peers)
		err := bp.Initialize()

		if (err != nil) != tc.expectError {
			t.Errorf("Expected error: %v, got: %v", tc.expectError, err)
		}
		if !tc.expectError {
			blocks, height := bp.GetBlocks()
			if !reflect.DeepEqual(blocks, tc.expectBlocks) || height != tc.expectHeight {
				t.Errorf("Expected blocks %v and height %d, got blocks %v and height %d",
					tc.expectBlocks, tc.expectHeight, blocks, height)
			}
		}
	}
}

func TestBlockProcessor_ProcessBlocks(t *testing.T) {
	bp := NewBlockProcessor(createZeroPeers(5))

	testCases := []struct {
		startHeight       uint64
		blocks            []Block
		expectedMaxHeight uint64
	}{
		{1, []Block{"A", "B", "C", "D", "E"}, 0},
		{1, []Block{"A", "B", "C", "D", "E"}, 0},
		{1, []Block{"A", "B", "C"}, 3},
		{4, []Block{"D", "E"}, 5},
		{6, []Block{"F", "G", "H"}, 5},
		{6, []Block{"F", "G", "H"}, 5},
		{6, []Block{"F", "G", "H"}, 8},
	}

	for _, tc := range testCases {
		maxConfirmedHeight := bp.ProcessBlocks(tc.startHeight, tc.blocks)
		if maxConfirmedHeight != tc.expectedMaxHeight {
			t.Errorf("Expected confirmedHeight to be %d, got %d", tc.expectedMaxHeight, maxConfirmedHeight)
		}
	}

	blocks, _ := bp.GetBlocks()
	expectedBlocks := []Block{"A", "B", "C", "D", "E", "F", "G", "H"}
	if !reflect.DeepEqual(blocks, expectedBlocks) {
		t.Errorf("Expected confirmedBlocks to be %v, got %v", expectedBlocks, blocks)
	}
}

func TestBlockProcessor_ConfirmBlock(t *testing.T) {
	bp := NewBlockProcessor(createZeroPeers(5))

	testCases := []struct {
		block        Block
		height       uint64
		expectBlocks []Block
		expectHeight uint64
	}{
		{"A", 1, []Block{"A"}, 1},
		{"B", 2, []Block{"A", "B"}, 2},
		{"C", 3, []Block{"A", "B", "C"}, 3},
		{"D", 5, []Block{"A", "B", "C", "", "D"}, 5},
		{"E", 4, []Block{"A", "B", "C", "E", "D"}, 5},
		{"F", 7, []Block{"A", "B", "C", "E", "D", "", "F"}, 7},
		{
			"Z",
			100,
			func() []Block {
				blocks := make([]Block, 50)
				for i := 0; i < 49; i++ {
					blocks[i] = ""
				}
				blocks[49] = "Z"
				return blocks
			}(),
			100,
		},
	}

	for _, tc := range testCases {
		bp.confirmBlock(tc.height, tc.block)
		blocks, maxConfirmedHeight := bp.GetBlocks()
		if !reflect.DeepEqual(blocks, tc.expectBlocks) || maxConfirmedHeight != tc.expectHeight {
			t.Errorf("Expected confirmedBlocks %v and height %d, got confirmedBlocks %v and height %d",
				tc.expectBlocks, tc.expectHeight, blocks, maxConfirmedHeight)
		}
	}
}
