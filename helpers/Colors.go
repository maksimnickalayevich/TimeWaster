package helpers

// Colorogo Simple struct that defines system's terminal colors
type Colorogo struct {
	Red    string
	Green  string
	Blue   string
	Cyan   string
	Yellow string
	Purple string
	Gray   string
	Reset  string
}

// InitColorogo Initializes struct with default values
// Red, Green, Yellow, Blue and etc.
func InitColorogo() Colorogo {
	return Colorogo{
		Red:    "\033[31m",
		Green:  "\033[32m",
		Yellow: "\033[33m",
		Blue:   "\033[34m",
		Purple: "\033[35m",
		Cyan:   "\033[36m",
		Gray:   "\033[37m",
		Reset:  "\033[0m",
	}
}

//TODO: Add string parsing and tag-like e.g. [Red] [/]
