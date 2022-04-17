package program

import "errors"

/**
 * Specifies if taskMaster should automatically restart a process if it exits when it is in the 'launched' state.
 * \always
 * The process will be unconditionally restarted when it exits, without regard to its exit code.
 * \never
 * The process will not be auto restarted.
 * \whenUnexpected
 * The process will be restarted when the program exits with an unexpected exit code.
 */
const (
	never      = "never"
	unexpected = "unexpected"
	always     = "always"
)

type Program struct {
	Name             string
	Command          string   `ini:"command"`
	ProcessNumber    int      `ini:"process_number"`
	AutoStart        bool     `ini:"auto_start"`
	StartTimeSec     int      `ini:"start_secs"`
	StartRetries     int      `ini:"start_retries"`
	AutoRestart      string   `ini:"auto_restart"`
	ExitCodes        []int    `delim:"," ini:"exit_codes"`
	StopSignal       int      `ini:"stop_signal"`
	StopWaitSecs     int      `ini:"stop_wait_secs"`
	StdErrLogfile    string   `ini:"log_file_error"`
	StdOutLogfile    string   `ini:"log_file"`
	Environments     []string `delim:"," ini:"env"`
	WorkingDirectory string   `ini:"working_directory"`
	Umask            string   `ini:"umask"`
}

func NewProgram() *Program {
	return &Program{
		ProcessNumber: 1,
		AutoStart:     true,
		StartTimeSec:  1,
		StartRetries:  3,
		AutoRestart:   unexpected,
		ExitCodes:     []int{0},
		StopSignal:    0, // TODO temp! Check real default value
		StopWaitSecs:  10,
	}
}

func (p *Program) commandIsCorrect() error {
	if len(p.Command) == 0 {
		return errors.New("program must have 'command' value")
	}
	return nil
}

func (p *Program) autoRestartIsCorrect() error {
	if p.AutoRestart != never && p.AutoRestart != unexpected && p.AutoRestart != always {
		return errors.New("auto_restart option should be '" + never + "', '" + unexpected + "' or '" + always + "' value")
	}
	return nil
}

func (p *Program) IsCorrect() error {
	if err := p.commandIsCorrect(); err != nil {
		return err
	} else if err = p.autoRestartIsCorrect(); err != nil {
		return err
	}
	return nil
}

func (p *Program) compareExitCode(lp *Program) bool {
	if len(p.ExitCodes) != len(lp.ExitCodes) {
		return false
	}
	for i := 0; i < len(p.ExitCodes); i++ {
		if p.ExitCodes[i] != lp.ExitCodes[i] {
			return false
		}
	}
	return true
}

func (p *Program) compareEnvironments(lp *Program) bool {
	if len(p.Environments) != len(lp.Environments) {
		return false
	}
	for i := 0; i < len(p.Environments); i++ {
		if p.Environments[i] != lp.Environments[i] {
			return false
		}
	}
	return true
}

func (p *Program) EqualTo(lp *Program) bool {
	return p.Name == lp.Name &&
		p.Command == lp.Command &&
		p.ProcessNumber == lp.ProcessNumber &&
		p.AutoStart == lp.AutoStart &&
		p.StartTimeSec == lp.StartTimeSec &&
		p.StartRetries == lp.StartRetries &&
		p.AutoRestart == lp.AutoRestart &&
		p.compareExitCode(lp) &&
		p.StopSignal == lp.StopSignal &&
		p.StopWaitSecs == lp.StopWaitSecs &&
		p.StdErrLogfile == lp.StdErrLogfile &&
		p.StdOutLogfile == lp.StdOutLogfile &&
		p.compareEnvironments(lp) &&
		p.WorkingDirectory == lp.WorkingDirectory &&
		p.Umask == lp.Umask
}
