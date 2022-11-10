package resources

import (
	"regexp"
	"strings"
)

// MaxAuthorLen is the maximum highscore author length
const MaxAuthorLen = 6

// RegexpForbiddenChars contains the list of forbidden characters is a highscore author
var RegexpForbiddenChars = regexp.MustCompile("[[:^alnum:]]")

// Highscore is a game highscore
type Highscore struct {
	Author    string
	Movements string
}

type HighscoreTable = map[string]Highscore

// NormalizeHighScores normalizes highscores
func NormalizeHighScores(t HighscoreTable) {
	for level, highscore := range t {
		highscore.Author = strings.ToUpper(RegexpForbiddenChars.ReplaceAllLiteralString(highscore.Author, ""))
		if len(highscore.Author) > MaxAuthorLen {
			highscore.Author = highscore.Author[:MaxAuthorLen]
		}
		t[level] = highscore
	}
}
