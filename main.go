package main

import (
	"log/slog"

	"github.com/qbradq/eye-engine/internal/ui"
	"github.com/qbradq/eye-engine/internal/util"
)

func main() {
	util.InitLog("eye-engine")
	m, err := ui.NewMain("Eye Engine")
	if err != nil {
		slog.Error("error creating main UI", "error", err)
	}
	m.Main()
}
