package models

type Options struct {
	General struct {
		Debug                bool
		DoNotRemoveTmpFolder bool
		Verbose              bool
		FailOnly             bool
	} `group:"General Options"`
	Files struct {
		Blacklist []string
		ListFiles bool
		PrintAST  bool
	} `group:"File Options"`
	Mutator struct {
		AvailableMutators []MutatorName
		DisableMutators   []MutatorName
		ListMutators      bool
	} `group:"Mutator Options"`
	Filter struct {
		Match string
	} `group:"Filter Options"`
	Exec struct {
		Exec    string
		NoExec  bool
		Timeout int64
		Jobs    int
	} `group:"Exec Options"`
	Test struct {
		Recursive bool
		Score     float64
	} `group:"Test Options"`
	Remaining struct {
		Targets []string
	}
}
