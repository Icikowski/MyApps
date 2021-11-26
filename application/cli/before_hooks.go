package cli

import (
	"os/user"
	"runtime"

	cliv2 "github.com/urfave/cli/v2"
)

func basicChecks(ctx *cliv2.Context) error {
	if runtime.GOOS == "windows" {
		return exitErrMsg("this application is not supported on Windows")
	}
	return nil
}

func allChecks(ctx *cliv2.Context) error {
	if err := basicChecks(ctx); err != nil {
		return err
	}

	currentUser, err := user.Current()
	if err != nil {
		return exitErrMsg("failed to get current user")
	}

	if currentUser.Username != "root" && len(ctx.Command.Name) > 0 {
		return exitErrMsg("this command must be executed as root")
	}

	return nil
}
