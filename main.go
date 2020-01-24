package main

import (
	"github.com/xortim/peruse/cmd"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	cmd.Execute()
}

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
