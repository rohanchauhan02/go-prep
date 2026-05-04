package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Vote int
const (Commit Vote = iota; Abort)
func (v Vote) String() string { if v == Commit { return "COMMIT" }; return "ABORT" }

type VoteMsg struct{ ParticipantID int; Vote Vote }
type DecisionMsg struct{ Decision Vote }

func participant(id int, voteCh chan<- VoteMsg, decisionCh <-chan DecisionMsg, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	vote := Commit
	if rand.Intn(4) == 0 { vote = Abort } // 25% chance of abort
	fmt.Printf("  Participant-%d voting: %s\n", id, vote)
	voteCh <- VoteMsg{id, vote}

	decision := <-decisionCh
	fmt.Printf("  Participant-%d executing: %s\n", id, decision.Decision)
}

func coordinator(n int) {
	voteChs := make([]chan VoteMsg, n)
	decisionChs := make([]chan DecisionMsg, n)
	for i := range voteChs {
		voteChs[i] = make(chan VoteMsg, 1)
		decisionChs[i] = make(chan DecisionMsg, 1)
	}

	var wg sync.WaitGroup
	// Phase 1: collect votes
	for i := 0; i < n; i++ {
		wg.Add(1)
		go participant(i+1, voteChs[i], decisionChs[i], &wg)
	}

	allCommit := true
	for i := 0; i < n; i++ {
		msg := <-voteChs[i]
		fmt.Printf("  Coordinator received vote from participant-%d: %s\n", msg.ParticipantID, msg.Vote)
		if msg.Vote == Abort { allCommit = false }
	}

	// Phase 2: broadcast decision
	decision := Abort
	if allCommit { decision = Commit }
	fmt.Printf("\nCoordinator decision: %s\n\n", decision)
	for i := 0; i < n; i++ {
		decisionChs[i] <- DecisionMsg{decision}
	}
	wg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=== Two-Phase Commit Simulation (3 participants) ===\n")
	coordinator(3)
	fmt.Println("Transaction complete")
}
