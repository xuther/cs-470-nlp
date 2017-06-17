package training

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Frequencies struct {
	Word         string
	Frequency    float64
	Addativefreq float64
}

func buildTransitionProbabilities(path string) map[string][]Frequencies {
	//we're gonna rad through the string, and split on spaces - for now we're just interested in word transition frequencies - so first step is to count, then we normalize.

	toReturn := make(map[string][]Frequencies)

	//initialize the map
	frequencyCount := make(map[string]map[string]int)
	totals := make(map[string]int)

	if file, err := os.Open(path); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)

		//read line by line
		for scanner.Scan() {
			curWord := "^"

			line := scanner.Text()
			if len(line) == 0 {
				continue
			}

			split := strings.Split(line, " ")
			if len(split) == 0 {
				continue
			}

			for _, word := range split {
				//make sure the word exists in the map
				if _, ok := frequencyCount[curWord]; !ok {
					frequencyCount[curWord] = make(map[string]int)
					frequencyCount[curWord][word] = 1
					totals[curWord] = 1
				} else if _, ok := frequencyCount[curWord][word]; !ok {
					frequencyCount[curWord][word] = 1
					totals[curWord] = 1 + totals[curWord]
				} else {
					frequencyCount[curWord][word] = frequencyCount[curWord][word] + 1
					totals[curWord] = 1 + totals[curWord]
				}
				curWord = word
			}
		}
	} else {
		log.Fatalf("Error opening file: %v: %v", path, err.Error())
	}

	//now we can trun our frequency counts into our transition probabilities
	for k, v := range frequencyCount {
		total := totals[k]
		freqs := []Frequencies{}

		for word, count := range v {
			freq := float64(count) / float64(total)
			newFreq := Frequencies{
				Word:      word,
				Frequency: freq,
			}
			freqs = append(freqs, newFreq)
		}

		//now we have all of our individual frequencies - go through and create 'addative' frequencies
		//we'll use these later to build sentences
		runningTotal := 0.0
		for i := range freqs {
			freqs[i].Addativefreq = runningTotal + freqs[i].Frequency
			runningTotal = freqs[i].Addativefreq
		}
		toReturn[k] = freqs
	}

	return toReturn
}

func Generate(freq map[string][]Frequencies) string {

	//for now assume that we start with the "^" character, and that we go from there
	sentence := ""
	words := []Frequencies{}

	curWord := "^"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(freq[curWord]) != 0 {
		randVal := r.Float64()
		for _, cur := range freq[curWord] {
			if randVal <= cur.Addativefreq {
				words = append(words, cur)
				curWord = cur.Word
				break
			}
		}
	}

	//build our sentence
	for i := range words {
		words := strings.Split(words[i].Word, "_")
		sentence += " " + words[0]
	}
	return sentence
}

func Main() {
	frequencies := buildTransitionProbabilities("training/trainingset.txt")

	log.Printf("Generating a new sentence")
	sentence := Generate(frequencies)
	log.Printf("%+v", sentence)

	return
}
