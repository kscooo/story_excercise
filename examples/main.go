package main

import (
	"fmt"
	"sync"

	blockprocessor "github.com/kscooo/home_excercise"
)

// MockPeer is a mock implementation of the Peer interface for testing purposes
type MockPeer struct {
	blocks []blockprocessor.Block
	height uint64
}

func (p *MockPeer) GetBlocks() ([]blockprocessor.Block, uint64) {
	return p.blocks, p.height
}

func main() {
	// Create some mock peers with different block data
	peer1 := &MockPeer{
		blocks: []blockprocessor.Block{"A", "B", "C", "D", "E"},
		height: 5,
	}
	peer2 := &MockPeer{
		blocks: []blockprocessor.Block{"A", "B", "C", "D", "E"},
		height: 5,
	}
	peer3 := &MockPeer{
		blocks: []blockprocessor.Block{"A", "B", "C", "D", "E"},
		height: 5,
	}
	peer4 := &MockPeer{
		blocks: []blockprocessor.Block{"X", "Y", "Z"},
		height: 3,
	}
	peer5 := &MockPeer{
		blocks: []blockprocessor.Block{"X", "Y", "Z"},
		height: 5,
	}

	// Create a new BlockProcessor with the mock peers
	peers := []blockprocessor.Peer{peer1, peer2, peer3, peer4, peer5}
	bp := blockprocessor.NewBlockProcessor(peers)

	// Initialize the BlockProcessor to synchronize with the majority of peers
	err := bp.Initialize()
	if err != nil {
		fmt.Printf("Error initializing BlockProcessor: %v\n", err)
		return
	}

	// Get the confirmed blocks and height
	confirmedBlocks, confirmedHeight := bp.GetBlocks()
	fmt.Printf("Initialize Confirmed Blocks: %v\n", confirmedBlocks)
	fmt.Printf("Initialize Confirmed Height: %d\n", confirmedHeight)

	// Process more blocks to update the confirmed blocks and height
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		bp.ProcessBlocks(6, []blockprocessor.Block{"F", "G", "H"})
	}()
	go func() {
		defer wg.Done()
		bp.ProcessBlocks(6, []blockprocessor.Block{"X", "Y", "Z"})
	}()
	go func() {
		defer wg.Done()
		bp.ProcessBlocks(6, []blockprocessor.Block{"F", "G", "H"})
	}()
	wg.Wait()

	// Get the updated confirmed blocks and height
	confirmedBlocks, confirmedHeight = bp.GetBlocks()
	fmt.Printf("Updated Confirmed Blocks: %v\n", confirmedBlocks)
	fmt.Printf("Updated Confirmed Height: %d\n", confirmedHeight)

	wg.Add(2)
	go func() {
		defer wg.Done()
		bp.ProcessBlocks(6, []blockprocessor.Block{"X", "Y", "Z"})
	}()
	go func() {
		defer wg.Done()
		bp.ProcessBlocks(6, []blockprocessor.Block{"X", "Y", "Z"})
	}()
	wg.Wait()

	// Get the final confirmed blocks and height
	confirmedBlocks, confirmedHeight = bp.GetBlocks()
	fmt.Printf("Final Confirmed Blocks: %v\n", confirmedBlocks)
	fmt.Printf("Final Confirmed Height: %d\n", confirmedHeight)
}
