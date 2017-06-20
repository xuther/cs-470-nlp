package training

type DynamicEntry struct {
	POS         string
	Probability float64
	PrevIndicie int
}

func Label(TransitionFreq map[string]Frequencies, EmissionFrequencies map[string]Frequencies, PossibleLabels []string, toLabel []string) []string {
	//viterbi algorithms

	//initialize our dynamic 'table' (really we're just starting at the beginning to find the most likely value)
	//transitionMatrix := [][]

	//initialize our dynamic tables
	dynamicTable := make([][]DynamicEntry, len(toLabel)+1)
	for i := range dynamicTable {
		if i == 0 {
			dynamicTable[i] = make([]DynamicEntry, 1)
			dynamicTable[i][0] = DynamicEntry{"^", 1, -1}
			continue
		}
		dynamicTable[i] = make([]DynamicEntry, len(PossibleLabels))
	}

	//run for the entire length building up our most likely label path
	for i := 1; i < len(toLabel)+1; i++ {
		for j := 0; j < len(PossibleLabels); j++ {
			//find the best transititon to the current label from the last state
			bestProb := 0.0
			bestVal := 0
			for k := 0; k < len(dynamicTable[i-1]); k++ {
				challenger := (dynamicTable[i-1][k].Probability * GetFreqForPOS(TransitionFreq[PossibleLabels[k]], PossibleLabels[j]).Frequency) * GetFreqForEmis(EmissionFrequencies[PossibleLabels[j]], toLabel[i-1]).Frequency
				if challenger > bestProb {
					bestProb = challenger
					bestVal = k
				}
			}

			dynamicTable[i][j] = DynamicEntry{PossibleLabels[j], bestProb, bestVal}
		}
	}

	//now that we've been through it all, we need to check the probability to the end state

	curBest := 0.0
	bestVal := 0

	//go through and find our best value
	for i := 0; i < len(PossibleLabels); i++ {
		if dynamicTable[len(dynamicTable)-1][i].Probability > curBest {
			curBest = dynamicTable[len(dynamicTable)-1][i].Probability
			bestVal = i
		}
	}

	labels := make([]string, len(toLabel))

	cur := bestVal
	//build our label set using the back pointer
	for i := len(toLabel); i > 0; i-- {
		curEntry := dynamicTable[i][cur]
		labels[i-1] = curEntry.POS
		cur = curEntry.PrevIndicie
	}

	return labels
}

func GetFreqForEmis(freqs []Frequency, word string) Frequency {
	for _, v := range freqs {
		if v.Word == word {
			return v
		}
	}
	return Frequency{word, "", 0.00000001, 0.0}
}

func GetFreqForPOS(freqs []Frequency, POS string) Frequency {
	for _, v := range freqs {
		if v.PartOfSpeech == POS {
			return v
		}
	}
	return Frequency{"", POS, 0.00000001, 0.0}
}
