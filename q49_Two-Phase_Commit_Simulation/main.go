//go:build ignore
// Remove the above line when implementing

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Vote int

const (Commit Vote = iota; Abort)

func (v Vote) String() string {
	if v == Commit { return "COMMIT" }
	return "ABORT"
}

// TODO: participant votes randomly (25% abort chance)
// sends vote to voteCh, then waits for decision on decisionCh
func participant(id int, voteCh chan<- Vote, decisionCh <-chan Vote, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

	// TODO: decide vote (Commit or Abort randomly)
	vote := Commit
	fmt.Printf("  P-%d votes: %s
", id, vote)
	voteCh <- vote

	// TODO: receive and print decision
	decision := <-decisionCh
	fmt.Printf("  P-%d executes: %s
", id, decision)
}

// TODO: coordinator runs 2-phase commit:
// Phase 1: collect all votes
// Phase 2: if all COMMIT → broadcast COMMIT, else broadcast ABORT
func coordinator(n int) {
	voteChans := make([]chan Vote, n)
	decisionChans := make([]chan Vote, n)
	for i := range voteChans {
		voteChans[i] = make(chan Vote, 1)
		decisionChans[i] = make(chan Vote, 1)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go participant(i+1, voteChans[i], decisionChans[i], &wg)
	}

	// TODO: collect votes, determine decision, broadcast
	allCommit := true
	for i := 0; i < n; i++ {
		v := <-voteChans[i]
		fmt.Printf("Coordinator got vote: %s
", v)
		if v == Abort { allCommit = false }
	}

	decision := Abort
	if allCommit { decision = Commit }
	fmt.Println("Decision:", decision)

	for i := 0; i < n; i++ { decisionChans[i] <- decision }
	wg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=== Two-Phase Commit (3 participants) ===")
	coordinator(3)
	fmt.Println("Done")
}
