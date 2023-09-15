package gen

import (
	"errors"
	"math"
)

type matchingAlgState struct {
	randSrc          IRandomSource
	isRoot           bool
	remainingMatches PossibleSecretSantas
	parentState      *matchingAlgState
	giver            string
	receiver         string
}

func generateSecretSantas(possibilities PossibleSecretSantas, randSrc IRandomSource) (map[string]string, error) {
	state := &matchingAlgState{
		remainingMatches: possibilities,
		isRoot:           true,
		randSrc:          randSrc,
	}

	for !state.solved() {
		newState, err := state.chooseNextGiver()
		state = newState
		if err != nil {
			return map[string]string{}, err
		}
	}

	return state.createPairings(), nil
}

func (state *matchingAlgState) createChild() *matchingAlgState {
	return &matchingAlgState{
		randSrc:          state.randSrc,
		isRoot:           false,
		parentState:      state,
		remainingMatches: copyPossibleSecretSantas(state.remainingMatches),
	}
}

func (state *matchingAlgState) solved() bool {
	return len(state.remainingMatches) == 0
}

func (state *matchingAlgState) chooseNextGiver() (*matchingAlgState, error) {
	giver, isOk := findMinPerson(state.remainingMatches)

	// if this states choice led to a dead end, jump back and mark
	// this choice as a dead end.
	if !isOk {
		oldState, err := state.revert()
		if err != nil {
			return state, err
		}
		oldState.removeChoice(state.giver, state.receiver)
		return oldState, nil
	}

	newState := state.createChild()
	newState.giver = giver
	newState.receiver = chooseReceiver(giver, state.remainingMatches, state.randSrc)
	delete(newState.remainingMatches, giver)
	removeChoice(newState.remainingMatches, newState.receiver)
	return newState, nil
}

func (state *matchingAlgState) revert() (*matchingAlgState, error) {
	if state.isRoot {
		return state, errors.New("at the root node")
	}
	return state.parentState, nil
}

func (state *matchingAlgState) createPairings() map[string]string {
	pairings := map[string]string{}
	currentState := state
	for {
		if currentState.isRoot {
			return pairings
		}
		pairings[currentState.giver] = currentState.receiver
		currentState = currentState.parentState
	}
}

func (state *matchingAlgState) removeChoice(giver string, receiving string) {
	state.remainingMatches[giver] = removeCandidate(state.remainingMatches[giver], receiving)
}

func findMinPerson(santas PossibleSecretSantas) (string, bool) {
	// find person with least possibilities
	var minUId string = ""
	minLen := math.MaxInt32
	for uid, possibleReceivers := range santas {
		if len(possibleReceivers) < minLen {
			minUId = uid
			minLen = len(possibleReceivers)
		}
	}
	if minLen == 0 {
		// current state is impossible to solve
		return "", false
	}
	return minUId, true
}

func chooseReceiver(uid string, santas PossibleSecretSantas, randSrc IRandomSource) string {
	selectionIndex := 0
	numChoices := len(santas[uid])
	if numChoices > 1 {
		selectionIndex = randSrc.Intn(numChoices - 1)
	}
	selectedUid := santas[uid][selectionIndex]
	//for cuid, possibleReceivers := range santas {
	//	santas[cuid] = removeCandidate(possibleReceivers, selectedUid)
	//}
	return selectedUid
}

func removeChoice(possiblities PossibleSecretSantas, person string) {
	for giver, choices := range possiblities {
		possiblities[giver] = removeCandidate(choices, person)
	}
}
