package player

import (
	"log/slog"

	"go.1password.io/spg"
)

var nameGenerator *spg.WLRecipe

func init() {
	wl, err := spg.NewWordList(spg.AgileWords)
	if err != nil {
		slog.Error("error generating played names", "error", err)
		return
	}
	nameGenerator = spg.NewWLRecipe(2, wl)
	nameGenerator.SeparatorChar = " "
}
