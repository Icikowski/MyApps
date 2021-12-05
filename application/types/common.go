package types

import (
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
)

var (
	headerFormatter      = color.New(color.FgHiWhite, color.Bold).SprintfFunc()
	firstColumnFormatter = color.New(color.FgGreen, color.Bold).SprintfFunc()
)

var progressBar = pb.ProgressBarTemplate(fmt.Sprintf(
	`  {{bar . "%s" "█" (cycle . "░" "▒" "▓" "▒") " " "%s"}} {{percent . }}`,
	firstColumnFormatter("│"), firstColumnFormatter("│"),
))

func finishProgressBar(bar *pb.ProgressBar) {
	bar.Finish()
	fmt.Print("\033[1A\033[K")
}

// Temp returns temporary directory that can be used by install / update scenarios
func Temp() (dir string, cleanup func(), err error) {
	dir, err = os.MkdirTemp("", "myapps")
	if err != nil {
		err = fmt.Errorf("Failed to create temp directory for this scenario: %w", err)
	}
	cleanup = func() {
		os.RemoveAll(dir)
	}
	return
}
