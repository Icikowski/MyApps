package common

import (
	"fmt"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	cliv2 "github.com/urfave/cli/v2"
)

// Formatters
var (
	// FmtHeader is a header formatting function
	FmtHeader = color.New(color.FgHiWhite, color.Bold).SprintfFunc()

	// FmtFirstCol is a first column formatting function
	FmtFirstCol = color.New(color.FgGreen, color.Bold).SprintfFunc()
)

// Box is a wrapper for box element
type Box string

func (b Box) Print(message string) {
	fmt.Printf("  %s %s\n", b, message)
}

// Boxes
var (
	// BoxInfo is a string representing blue-colored information sign box
	BoxInfo = Box(color.New(color.FgHiBlue, color.Bold).Sprint("[i]"))

	// BoxError is a string representing red-colored error sign box
	BoxErr = Box(color.New(color.FgHiRed, color.Bold).Sprint("[✗]"))

	// BoxWarn is a string representing yellow-colored warning sign box
	BoxWarn = Box(color.New(color.FgHiYellow, color.Bold).Sprint("[!]"))

	// BoxSuccess is a string representing green-colored success sign box
	BoxSuccess = Box(color.New(color.FgHiGreen, color.Bold).Sprint("[✓]"))
)

// ProgressBar is a thread-safe wrapper for the progress bar
type ProgressBar struct {
	pb    *pb.ProgressBar
	mutex sync.Mutex
}

// NewProgressBar creates a new progress bar with given capacity
func NewProgressBar(length int) *ProgressBar {
	progressBar := pb.ProgressBarTemplate(fmt.Sprintf(
		`  {{bar . "%s" "█" (cycle . "░" "▒" "▓" "▒") " " "%s"}} {{percent . }}`,
		FmtFirstCol("│"), FmtFirstCol("│"),
	)).New(length)

	defer progressBar.Start()

	return &ProgressBar{
		pb: progressBar,
	}
}

// Increment increments the progress bar by one
func (pb *ProgressBar) Increment() {
	defer pb.mutex.Unlock()
	pb.mutex.Lock()
	pb.pb.Increment()
}

// Finish finishes the progress bar and cleans the line
func (pb *ProgressBar) Finish() {
	defer pb.mutex.Unlock()
	pb.mutex.Lock()
	pb.pb.Finish()
	fmt.Print("\033[1A\033[K")
}

// ExitWithErrMsg exits application with error message
func ExitWithErrMsg(msg string) error {
	return cliv2.Exit(color.New(color.FgRed, color.Bold).Sprint("Error:")+" "+msg, 1)
}

// ExitWithWarnMsg exits application with warning message
func ExitWithWarnMsg(msg string) error {
	return cliv2.Exit(color.New(color.FgYellow, color.Bold).Sprint("Warning:")+" "+msg, 2)
}

// PrintErrorWhileMsg prints formatted "Error while <while> <subject>: <err>" message
func PrintErrorWhileMsg(while, subject string, err error) {
	fmt.Printf("%s while %s %s: %s\n",
		color.RedString("Error"), while,
		color.BlueString(subject), err.Error(),
	)
}

// PrintSuccessfullyMsg prints formatted "Successfully <what> <subject>" message
func PrintSuccessfullyMsg(what, subject string) {
	fmt.Printf("%s %s %s\n",
		color.GreenString("Successfully"), what,
		color.BlueString(subject),
	)
}
