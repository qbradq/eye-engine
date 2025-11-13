package main

import (
	"log/slog"
	"os"

	"github.com/qbradq/eye-engine/internal/ui"
	"github.com/qbradq/eye-engine/internal/util"
)

func main() {
	util.InitLog("eye-engine")
	m := ui.NewMain("Eye Engine")
	if m == nil {
		slog.Error("error creating main UI")
		os.Exit(1)
	}
	m.Main()
}
