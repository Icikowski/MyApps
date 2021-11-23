package cli

import (
	"fmt"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	cliv2 "github.com/urfave/cli/v2"
)

func printErrorWhileMsg(while, subject string, err error) {
	fmt.Printf("%s while %s %s: %s\n",
		color.RedString("Error"), while,
		color.BlueString(subject), err.Error(),
	)
}

func printSuccessfully(what, subject string) {
	fmt.Printf("%s %s %s\n",
		color.GreenString("Successfully"), what,
		color.BlueString(subject),
	)
}

func exitErrMsg(msg string) error {
	return cliv2.Exit(color.New(color.FgRed, color.Bold).Sprint("Error:")+" "+msg, 1)
}

func exitWarnMsg(msg string) error {
	return cliv2.Exit(color.New(color.FgYellow, color.Bold).Sprint("Warning:")+" "+msg, 2)
}

var (
	headerFormatter      = color.New(color.FgHiWhite, color.Bold).SprintfFunc()
	firstColumnFormatter = color.New(color.FgGreen, color.Bold).SprintfFunc()
)

var progressBar = pb.ProgressBarTemplate(`{{bar . "│" "█" (cycle . "░" "▒" "▓" "▒") " " "│"}} {{percent . }}`)

func finishProgressBar(bar *pb.ProgressBar) {
	bar.Finish()
	fmt.Print("\033[1A\033[K")
}

var (
	infoBox    = color.New(color.FgHiBlue, color.Bold).Sprint("[i]")
	errorBox   = color.New(color.FgHiRed, color.Bold).Sprint("[✗]")
	warningBox = color.New(color.FgHiYellow, color.Bold).Sprint("[!]")
	successBox = color.New(color.FgHiGreen, color.Bold).Sprint("[✓]")
)

func printBox(box, message string) {
	fmt.Printf("  %s %s\n", box, message)
}
