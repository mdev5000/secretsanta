package gen

import (
	"errors"
	"fmt"
)

type People = []string

type HistoryItem = map[string]string

type IRandomSource interface {
	Shuffle(n int, swap func(i, j int))
	Intn(n int) int
}

type PossibleSecretSantas = map[string][]string

func GenPairing(people People, history []HistoryItem, randSrc IRandomSource) (map[string]string, error) {
	mustDifferFrom := history
	exclusions := PossibleSecretSantas{}
	possibilities := PossibleSecretSantas{}
	pairings := map[string]string{}

	// figure out who people cant have

	for _, c := range people {
		exclusions[c] = make(People, 0, len(people))
		// a candidate cannot have themselves
		exclusions[c] = append(exclusions[c], c)
	}

	for _, exclusionGroup := range mustDifferFrom {
		for giving, receiving := range exclusionGroup {
			exclusions[giving] = append(exclusions[giving], receiving)
		}
	}

	// setup all possibles

	for _, c := range people {
		possibilities[c] = copyCandidates(people)
		for _, exclusion := range exclusions[c] {
			possibilities[c] = removeCandidate(possibilities[c], exclusion)
		}
	}

	pairings, err := generateSecretSantas(possibilities, randSrc)

	return pairings, err
}

func removeCandidate(candidates People, uid string) People {
	newCandidates := make(People, 0, len(candidates))
	for _, c := range candidates {
		if c != uid {
			newCandidates = append(newCandidates, c)
		}
	}
	return newCandidates
}

func ValidatePairing(pairings map[string]string, history []map[string]string) error {
	if !validateNobodyHasThemselves(pairings) {
		return errors.New("error validating pairings, somebody has themselves")
	}
	for giving, receiving := range pairings {
		for _, oldPairing := range history {
			if oldPairing[giving] == receiving {
				return fmt.Errorf("candidate %s has them same person as a previous year", giving)
			}
		}
	}
	return nil
}

func validateNobodyHasThemselves(pairings map[string]string) bool {
	for giving, receiving := range pairings {
		if giving == receiving {
			return false
		}
	}
	return true
}

func copyCandidates(s People) People {
	cp := make(People, len(s))
	copy(cp, s)
	return cp
}

func copyPossibleSecretSantas(s PossibleSecretSantas) PossibleSecretSantas {
	newS := make(PossibleSecretSantas, len(s))
	for k, v := range s {
		newS[k] = v
	}
	return newS
}
