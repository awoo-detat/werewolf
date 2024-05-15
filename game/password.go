package game

import (
	"log/slog"

	"go.1password.io/spg"
)

var passwordGenerator *spg.WLRecipe

func init() {
	wordList, err := spg.NewWordList([]string{
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

	if err != nil {
		slog.Error("error generating passwords", "error", err)
		return
	}
	passwordGenerator = spg.NewWLRecipe(3, wordList)
	passwordGenerator.SeparatorChar = " "
}
