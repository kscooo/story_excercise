package blockprocessor

import (
	"fmt"
	"sync"
)

type Block string

type Peer interface {
	GetBlocks() ([]Block, uint64)
}

// BlockProcessor is responsible for tracking the most recent 50 confirmed blocks and the latest confirmed block height
type BlockProcessor struct {
	confirmedBlocks    []Block
	confirmedHeight    uint64
	blockConfirmations map[uint64]map[Block]int
	peerNodes          []Peer

	mutex sync.RWMutex
}

// NewBlockProcessor creates a new BlockProcessor instance
func NewBlockProcessor(peers []Peer) *BlockProcessor {
	p := &BlockProcessor{
		confirmedBlocks:    make([]Block, 0, 50),
		blockConfirmations: make(map[uint64]map[Block]int),
		peerNodes:          peers,
	}
	return p
}

// GetBlocks returns the most recent 50 confirmed blocks and the latest confirmed block height
func (p *BlockProcessor) GetBlocks() ([]Block, uint64) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	blocks := make([]Block, len(p.confirmedBlocks))
	copy(blocks, p.confirmedBlocks)
	return blocks, p.confirmedHeight
}

// Initialize synchronizes the BlockProcessor with peer nodes
func (p *BlockProcessor) Initialize() error {
	type item struct {
		blocks []Block
		height uint64
	}

	responses := make([]item, len(p.peerNodes))
	var wg sync.WaitGroup
	for i, peer := range p.peerNodes {
		wg.Add(1)
		go func(i int, peer Peer) {
			defer wg.Done()

			blocks, height := peer.GetBlocks()
			responses[i] = item{
				blocks: blocks,
				height: height,
			}
		}(i, peer)
	}
	wg.Wait()

	peerCounts := make(map[string]int)
	var majorityBlocks []Block
	var majorityHeight uint64
	for _, response := range responses {
		key := fmt.Sprintf("%v:%v", response.blocks, response.height)
		peerCounts[key]++
		if peerCounts[key] >= (len(p.peerNodes)+1)/2 {
			majorityBlocks = response.blocks
			majorityHeight = response.height
			break
		}
	}

	if majorityHeight != uint64(len(majorityBlocks)) {
		return fmt.Errorf("majority height does not match the number of blocks, height: %d, len(blocks): %d",
			majorityHeight, len(majorityBlocks))
	}

	if len(majorityBlocks) > 50 {
		majorityBlocks = majorityBlocks[len(majorityBlocks)-50:]
	}

	p.confirmedBlocks = majorityBlocks
	p.confirmedHeight = majorityHeight
	return nil
}

// ProcessBlocks processes a range of blocks and returns the max confirmed block height
func (p *BlockProcessor) ProcessBlocks(startHeight uint64, blocks []Block) uint64 {
	if len(blocks) == 0 || startHeight == 0 {
		return p.confirmedHeight
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i, block := range blocks {
		height := startHeight + uint64(i)

		if _, ok := p.blockConfirmations[height]; !ok {
			p.blockConfirmations[height] = make(map[Block]int)
		}
		p.blockConfirmations[height][block]++

		if p.blockConfirmations[height][block] >= (len(p.peerNodes)+1)/2 {
			p.confirmBlock(height, block)
		}
	}

	return p.confirmedHeight
}

func (p *BlockProcessor) confirmBlock(height uint64, block Block) {
	if height > p.confirmedHeight {
		// Fill in missing blocks with empty blocks if any
		for i := p.confirmedHeight + 1; i < height; i++ {
			p.confirmedBlocks = append(p.confirmedBlocks, "")
		}
		p.confirmedBlocks = append(p.confirmedBlocks, block)
		p.confirmedHeight = height
	} else if height <= p.confirmedHeight {
		// Potentially update blocks in between
		index := height - (p.confirmedHeight - uint64(len(p.confirmedBlocks)) + 1)
		if index >= 0 && index < uint64(len(p.confirmedBlocks)) {
			p.confirmedBlocks[index] = block
		}
	}

	// Trim blocks to keep only the most recent 50
	if len(p.confirmedBlocks) > 50 {
		p.confirmedBlocks = p.confirmedBlocks[len(p.confirmedBlocks)-50:]
	}
}

// Reset clears the confirmed blocks and height
func (p *BlockProcessor) Reset() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.confirmedBlocks = make([]Block, 0, 50)
	p.confirmedHeight = 0
	p.blockConfirmations = make(map[uint64]map[Block]int)
}
