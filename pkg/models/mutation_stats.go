package models

type MutationStats struct {
	MutantsEscaped int
	MutantsKilled  int
	MutantsSkipped int
	Duplicated     int
	UnknownResults int
}

func (ms MutationStats) Score() float64 {
	if ms.Total() == 0 {
		return 0.0
	}

	return float64(ms.MutantsKilled) / float64(ms.Total())
}

func (ms MutationStats) Total() int {
	return ms.MutantsEscaped + ms.MutantsKilled
}
