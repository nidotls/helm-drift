package command

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/nikhilsbhat/helm-drift/pkg/deviation"
)

func (cmd *command) RunKubeCmd(deviation deviation.Deviation) (deviation.Deviation, error) {
	cmd.log.Debugf("envionment variables that would be used: %v", cmd.baseCmd.Environ())

	out, err := cmd.baseCmd.CombinedOutput()
	if err != nil {
		var exerr *exec.ExitError
		if errors.As(err, &exerr) {
			switch exerr.ExitCode() {
			case 1:
				deviation.HasDrift = true
				deviation.Deviations = string(out)
				cmd.log.Debugf("found diffs for '%s' with name '%s'", deviation.Kind, deviation.Kind)
			default:
				return deviation, fmt.Errorf("running kubectl diff errored with exit code: %w ,with message: %s", err, string(out))
			}
		}
	} else {
		cmd.log.Debugf("no diffs found for '%s' with name '%s'", deviation.Kind, deviation.Kind)
	}

	return deviation, nil
}