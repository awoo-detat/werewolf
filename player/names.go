package player

import (
	"go.1password.io/spg"
)

var nameGenerator *spg.WLRecipe

func init() {
	wl, _ := spg.NewWordList(spg.AgileWords)
	nameGenerator = spg.NewWLRecipe(2, wl)
	nameGenerator.SeparatorChar = " "
}
