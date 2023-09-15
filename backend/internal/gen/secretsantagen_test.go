package gen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func FakeRandomFn() IRandomSource {
	return rand.New(rand.NewSource(6))
}

func GenerateRandomNames(rand *rand.Rand, maxNumNames int) People {
	numNames := 0
	for numNames < 2 {
		numNames = rand.Intn(maxNumNames)
	}

	names := map[string]bool{}

	count := 0
	for count < numNames {
		names[GenerateName(count+1)] = true
		count++
	}

	people := make(People, numNames)
	i := 0
	for name, _ := range names {
		people[i] = name
		i++
	}

	return people
}

func GenerateName(count int) string {
	return string(fmt.Sprintf("a%d", count))
}

func pruneHistory(uniqueFor int, history []HistoryItem) []HistoryItem {
	start := len(history) - uniqueFor
	if start < 0 {
		start = 0
	}
	return history[start:]
}

func TestCanGenerateAValidPairingForTheMinimalAmountOfData(t *testing.T) {
	people := People{"John", "Bob"}
	pairing, err := GenPairing(people, []map[string]string{}, FakeRandomFn())
	assert.Nil(t, err)
	assert.Equal(t, pairing, map[string]string{
		"John": "Bob",
		"Bob":  "John",
	})
}

func TestErrorWhenCannotGeneratePairing(t *testing.T) {
	people := People{"John", "Bob"}
	history := []map[string]string{
		{
			"John": "Bob",
			"Bob":  "John",
		},
	}
	_, err := GenPairing(people, history, FakeRandomFn())
	assert.EqualError(t, err, "at the root node")
}

func TestTrickyGenerateAlg(t *testing.T) {
	for i := 0; i < 2000; i++ {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		possiblities := PossibleSecretSantas{
			"a3": []string{"a2"},
			"a4": []string{"a1"},
			"a5": []string{"a3"},
			"a1": []string{"a4"},
			"a2": []string{"a5"},
		}
		_, err := generateSecretSantas(possiblities, random)
		assert.Nil(t, err)
	}
}

func TestTrickyGenerateAlg2(t *testing.T) {
	for i := 0; i < 50; i++ {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		possiblities := PossibleSecretSantas{
			"a3": []string{"a4", "a5", "a2"},
			"a4": []string{"a3", "a1", "a2"},
			"a5": []string{"a3", "a1", "a2"},
			"a1": []string{"a3", "a4", "a5"},
			"a2": []string{"a4", "a5", "a1"},
		}
		_, err := generateSecretSantas(possiblities, random)
		assert.Nil(t, err, fmt.Sprintf("error running %d", i))
	}
}

func TestTrickyCasesCanWork(t *testing.T) {
	for j := 0; j < 20; j++ {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		peopleInit := []string{"a3", "a4", "a5", "a1", "a2"}
		history := []map[string]string{{
			"a2": "a3",
			"a3": "a1",
			"a1": "a2",
			"a5": "a4",
			"a4": "a5",
		}}

		numUnique := len(peopleInit) - 2
		//for i := 0; i < numUnique; i++ {
		people := copyCandidates(peopleInit)
		history = pruneHistory(numUnique, history)
		pairings, err := GenPairing(people, history, random)
		assert.Nil(t, err, fmt.Sprintf("Failed on %d", j))
		assert.Nil(t, ValidatePairing(pairings, history))
		history = append(history, pairings)
		//}
	}
}

func TestTrickyCasesCanWork2(t *testing.T) {
	for i := 0; i < 2000; i++ {
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		peopleInit := []string{"a3", "a4", "a5", "a1", "a2"}
		history := []map[string]string{
			{"a2": "a3", "a3": "a1", "a1": "a2", "a5": "a4", "a4": "a5"},
			//{"a2": "a5", "a3": "a4", "a1": "a3", "a5": "a2", "a4": "a1"},
			{"a2": "a1", "a3": "a5", "a1": "a4", "a5": "a3", "a4": "a2"},
			{"a2": "a4", "a3": "a2", "a1": "a5", "a5": "a1", "a4": "a3"},
		}

		people := copyCandidates(peopleInit)
		uniqueFor := len(peopleInit) - 1
		pairings, err := GenPairing(people, pruneHistory(uniqueFor, history), random)
		assert.Nil(t, err)
		assert.Nil(t, ValidatePairing(pairings, history))
	}
}

func TestCanGenerateNMinusOneValidatePermutations(t *testing.T) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	peopleInit := People{"John", "Bob", "Edna", "Laura"}
	history := []map[string]string{}
	numUnique := len(peopleInit) - 1
	for i := 0; i < numUnique; i++ {
		people := copyCandidates(peopleInit)
		pairings, err := GenPairing(people, pruneHistory(numUnique, history), random)
		assert.Nil(t, err)
		assert.Nil(t, ValidatePairing(pairings, history))
		history = append(history, pairings)
	}
}

func TestNeverGeneratesAnInvalidPairingWhenPossibleToGenerateAValidOne(t *testing.T) {
	numIterations := 50
	maxNumOfPeople := 40

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numIterations; i++ {
		peopleInit := GenerateRandomNames(random, maxNumOfPeople)
		numUnique := len(peopleInit) - 1
		history := make([]map[string]string, 0, numUnique)

		for j := 0; j < numUnique; j++ {
			if len(history) != j {
				assert.Failf(t, "Bad invariant", "Expected history to match j")
			}

			people := copyCandidates(peopleInit)
			pairings, err := GenPairing(people, pruneHistory(numUnique, history), random)

			if (err != nil) || (ValidatePairing(pairings, history) != nil) {
				assert.Failf(
					t,
					"Failed to generate valid pairing", "Failed to generate valid pairing\ncandidates: %s\niteration: %d\nerror: %e\npairings: %s\nhistory: %s\n",
					peopleInit,
					j,
					err,
					pairings,
					history,
				)
			}
			history = append(history, pairings)
		}
	}
}
