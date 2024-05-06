# BlockProcessor
BlockProcessor is a Go package that provides functionality for tracking and processing confirmed blocks in a blockchain network. It is designed to keep track of the most recent 50 confirmed blocks and the latest confirmed block height.
## Design Thinking and Assumptions

- Concurrency: The BlockProcessor is designed to handle concurrent requests and updates safely. It uses a read-write mutex (sync.RWMutex) to synchronize access to the shared state.
- Fault Tolerance: The BlockProcessor is designed to handle scenarios where some peer nodes may provide incorrect or malicious block data. It relies on the majority consensus among the peer nodes to determine the correct block data.
- Initialization: When the BlockProcessor is initialized, it synchronizes its state with the peer nodes by calling their GetBlocks method. It selects the majority response as the initial state.
- Block Processing: The BlockProcessor allows processing a range of blocks with the ProcessBlocks method. It keeps track of the confirmation count for each block at each height and confirms a block when it reaches the majority confirmation threshold.
- Block Confirmation: When a block is confirmed, the BlockProcessor updates its internal state, including the confirmed blocks and the latest confirmed block height. It fills in any missing blocks with empty blocks to maintain the integrity of the block chain.
- Recent Block Tracking: The BlockProcessor keeps track of the most recent 50 confirmed blocks. When new blocks are confirmed, older blocks are trimmed to maintain this limit.

## Setup and Running

1. Install the package:
```
go get github.com/kscooo/home_excercise/blockprocessor
```

2. Import the package in your Go code:
```go
import "github.com/kscooo/home_excercise"
```

3. Create an instance of the BlockProcessor:
```go
peers := []blockprocessor.Peer{peer1, peer2, peer3}
bp := blockprocessor.NewBlockProcessor(peers)
```

4. Initialize the BlockProcessor:
```go
err := bp.Initialize()
if err != nil {
    // Handle the error
}
```

5. Process blocks:
```go
maxConfirmedHeight := bp.ProcessBlocks(startHeight, blocks)
```

6. Get the confirmed blocks and height:
```go
confirmedBlocks, confirmedHeight := bp.GetBlocks()
```
## Testing
To run the tests:
```
go test ./...
```
The tests cover various scenarios, including:

- Initialization with different peer responses
- Processing blocks in order and out of order
- Handling missing blocks
- Trimming blocks to maintain the recent block limit

## Example Usage
An example usage of the BlockProcessor package can be found in the examples/ directory. It demonstrates how to create a BlockProcessor instance, initialize it, process blocks, and retrieve the confirmed blocks and height.
To run the example:
```
go run examples/main.go
```
