package types

import (
	"fmt"

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
