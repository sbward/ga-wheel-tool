package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/thoj/go-galib"
)

var timeout bool
var scores int
var input = struct {
	ToolOrder []int
}{}

type appleToolMutator struct {
	gaussian ga.GAMutator
}

func (a *appleToolMutator) Mutate(g ga.GAGenome) ga.GAGenome {
	return a.gaussian.Mutate(g)
}

func (a *appleToolMutator) String() string { return "appleToolMutator" }

type appleToolSelector struct {
	tournament ga.GASelector
}

func (s *appleToolSelector) SelectOne(pop ga.GAGenomes) ga.GAGenome {
	// Pick the one that has the shortest Tool Selection Path
	// Assume ga.GAGenomes are type ga.GAOrderedIntGenome
	var best ga.GAGenome
	var bestLength int
	genomes := make([]*ga.GAOrderedIntGenome, len(pop))
	for i, g := range pop {
		genome := g.(*ga.GAOrderedIntGenome)
		length := toolSelectionPathLength(genome)
		if length < bestLength || bestLength == 0 {
			bestLength = length
			best = genome
		}
	}
	return best
}

func (s *appleToolSelector) String() string { return "appleToolSelector" }

var dedupInput []int

func getDedupInput() []int {
	if dedupInput != nil {
		return dedupInput
	}
	dedupInput := removeDuplicates(input.ToolOrder)
	sort.Ints(dedupInput)
	return dedupInput
}

// Returns the tool selection path length or zero if it's impossible
func toolSelectionPathLength(g *ga.GAOrderedIntGenome) int {
	dedupInput := getDedupInput()
	dedupGene := removeDuplicates(g.Gene)
	if !reflect.DeepEqual(dedupInput, dedupGene) {
		// This gene does not contain all of the tools
		return 0
	}
	// Follow the input path and measure the distance of each jump in the gene
	// fixme(sam): Assumes starting position is index 0
	var total int
	for i, t := range g.Gene {

	}
}

// Absolute value for integers
func iabs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func removeDuplicates(a []int) []int {
	result := []int{}
	seen := map[int]int{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	sort.Ints(result)
	return result
}

func main() {
	err := json.Unmarshal([]byte(strings.Join(os.Args[1:], "")), &input.ToolOrder)
	if err != nil {
		log.Fatal("Must pass [1,2,3] argument.", err)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	gao := ga.NewGA(ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector: &appleToolSelector{
			tournament: ga.NewGATournamentSelector(0.2, 5),
		},
		Breeder: new(ga.GA2PointBreeder),
		Mutator: &appleToolMutator{
			gaussian: ga.NewGAGaussianMutator(0.4, 0),
		},
		PMutate: 0.5,
		PBreed:  0.2,
	})

	fitness := func(ga *ga.GAOrderedIntGenome) float64 {
		// Check if all tools exist in answer; if not, score = 0

		// Check tool path length; score = (1 / len)

		return 1
	}

	genome := ga.NewOrderedIntGenome(input.ToolOrder, fitness)

	gao.Init(1000, genome) //Total population

	go func() {
		<-time.After(5 * time.Second)
		timeout = true
	}()

	gao.OptimizeUntil(func(best ga.GAGenome) bool {
		scores++
		if timeout {
			return true
		}
		return best.Score() <= 680
	})

	best := gao.Best().(*ga.GAOrderedIntGenome)

	fmt.Printf("%s = %f\n", best, best.Score())
	fmt.Printf("Calls to score = %d\n", scores)
}
