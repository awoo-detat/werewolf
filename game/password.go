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
		"good",
		"evil",
		"tinker",
		"tailor",
		"solder",
		"spy",
		"atrocities",
	})
	passwordGenerator = spg.NewWLRecipe(3, wordList)
	passwordGenerator.SeparatorChar = " "
}
