package game

import (
	"go.1password.io/spg"
)

var passwordGenerator *spg.WLRecipe

func init() {
	wordList, _ := spg.NewWordList([]string{
		"furry",
		"awoo",
		"werewolf",
		"villager",
		"seer",
		"correct",
		"horse",
		"battery",
		"staple",
		"village",
	})
	passwordGenerator = spg.NewWLRecipe(3, wordList)
}
