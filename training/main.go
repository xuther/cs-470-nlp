package training

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

type Frequencies []Frequency

func (f Frequencies) Len() int {
	return len(f)
}

func (f Frequencies) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f Frequencies) Less(i, j int) bool {
	return f[i].Frequency > f[j].Frequency
}

type Frequency struct {
	Word         string
	PartOfSpeech string
	Frequency    float64
	Addativefreq float64
}

func buildTransitionProbabilities(path string) (map[string]Frequencies, map[string]Frequencies) {
	//we're gonna rad through the string, and split on spaces - for now we're just interested in word transition frequencies - so first step is to count, then we normalize.

	toReturnPos := make(map[string]Frequencies)
	toReturnWord := make(map[string]Frequencies)

	//initialize the map
	frequencyCount := make(map[string]map[string]int)
	totals := make(map[string]int)

	wordToPartFrequencies := make(map[string]map[string]int)
	wordTotals := make(map[string]int)

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
				splits := strings.Split(word, "_")
				word := splits[0]
				pos := splits[1]
				//make sure the word exists in the map
				if _, ok := frequencyCount[curWord]; !ok {
					frequencyCount[curWord] = make(map[string]int)
					frequencyCount[curWord][pos] = 1
					totals[curWord] = 1
				} else if _, ok := frequencyCount[curWord][pos]; !ok {
					frequencyCount[curWord][pos] = 1
					totals[curWord] = 1 + totals[curWord]
				} else {
					frequencyCount[curWord][pos] = frequencyCount[curWord][pos] + 1
					totals[curWord] = 1 + totals[curWord]
				}
				curWord = pos

				if _, ok := wordTotals[word]; !ok {
					wordTotals[word] = 1
					wordToPartFrequencies[word] = make(map[string]int)
					wordToPartFrequencies[word][pos] = 1
					wordTotals[word] = 1
				} else if _, ok := wordToPartFrequencies[word][pos]; !ok {
					wordToPartFrequencies[word][pos] = 1
					wordTotals[word] = 1 + wordTotals[word]
				} else {
					wordToPartFrequencies[word][pos]++
					wordTotals[word] = 1 + wordTotals[word]
				}
			}
		}
	} else {
		log.Fatalf("Error opening file: %v: %v", path, err.Error())
	}

	//now we can trun our frequency counts into our transition probabilities
	for k, v := range frequencyCount {
		total := totals[k]
		freqs := Frequencies{}

		for word, count := range v {
			freq := float64(count) / float64(total)
			newFreq := Frequency{
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
		sort.Sort(freqs)
		toReturnPos[k] = freqs
	}

	for k, v := range wordToPartFrequencies {
		total := wordTotals[k]
		freqs := Frequencies{}

		for pos, count := range v {
			freq := float64(count) / float64(total)
			newFreq := Frequency{
				PartOfSpeech: pos,
				Frequency:    freq,
			}
			freqs = append(freqs, newFreq)
		}
		runningTotal := 0.0
		for i := range freqs {
			freqs[i].Addativefreq = runningTotal + freqs[i].Frequency
			runningTotal = freqs[i].Addativefreq
		}
		sort.Sort(freqs)
		toReturnWord[k] = freqs
	}

	return toReturnPos, toReturnWord
}

func Generate(freq map[string][]Frequency) string {

	//for now assume that we start with the "^" character, and that we go from there
	sentence := ""
	words := []Frequency{}

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
	frequencies, word := buildTransitionProbabilities("training/trainingset.txt")

	log.Printf("Frequencies: ")
	log.Printf("Base: NN")
	log.Printf("Vals: ")
	for _, val := range frequencies["JJ"] {
		log.Printf("%v: %v", val.Word, val.Frequency)
	}

	log.Printf("Word Frequencies")
	log.Printf("Base: the")
	log.Printf("Vals: ")
	for _, val := range word["the"] {
		log.Printf("%v: %v", val.PartOfSpeech, val.Frequency)
	}

	//sentence := Generate(frequencies)
	//log.Printf("%+v", sentence)

	return
}
