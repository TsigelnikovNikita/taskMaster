package task

import "errors"

/**
 * Specifies if taskMaster should automatically restart a process if it exits when it is in the 'launched' state.
 * \always
 * The process will be unconditionally restarted when it exits, without regard to its exit code.
 * \never
 * The process will not be auto restarted.
 * \whenUnexpected
 * The process will be restarted when the task exits with an unexpected exit code.
 */
const (
	never      = "never"
	unexpected = "unexpected"
	always     = "always"
)

type Task struct {
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

func NewTask() *Task {
	return &Task{
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

func (t *Task) commandIsCorrect() error {
	if len(t.Command) == 0 {
		return errors.New("task must have 'command' value")
	}
	return nil
}

func (t *Task) autoRestartIsCorrect() error {
	if t.AutoRestart != never && t.AutoRestart != unexpected && t.AutoRestart != always {
		return errors.New("auto_restart option should be '" + never + "', '" + unexpected + "' or '" + always + "' value")
	}
	return nil
}

func (t *Task) IsCorrect() error {
	if err := t.commandIsCorrect(); err != nil {
		return err
	} else if err = t.autoRestartIsCorrect(); err != nil {
		return err
	}
	return nil
}

func (t *Task) compareExitCode(lp *Task) bool {
	if len(t.ExitCodes) != len(lp.ExitCodes) {
		return false
	}
	for i := 0; i < len(t.ExitCodes); i++ {
		if t.ExitCodes[i] != lp.ExitCodes[i] {
			return false
		}
	}
	return true
}

func (t *Task) compareEnvironments(lp *Task) bool {
	if len(t.Environments) != len(lp.Environments) {
		return false
	}
	for i := 0; i < len(t.Environments); i++ {
		if t.Environments[i] != lp.Environments[i] {
			return false
		}
	}
	return true
}

func (t *Task) EqualTo(lp *Task) bool {
	return t.Name == lp.Name &&
		t.Command == lp.Command &&
		t.ProcessNumber == lp.ProcessNumber &&
		t.AutoStart == lp.AutoStart &&
		t.StartTimeSec == lp.StartTimeSec &&
		t.StartRetries == lp.StartRetries &&
		t.AutoRestart == lp.AutoRestart &&
		t.compareExitCode(lp) &&
		t.StopSignal == lp.StopSignal &&
		t.StopWaitSecs == lp.StopWaitSecs &&
		t.StdErrLogfile == lp.StdErrLogfile &&
		t.StdOutLogfile == lp.StdOutLogfile &&
		t.compareEnvironments(lp) &&
		t.WorkingDirectory == lp.WorkingDirectory &&
		t.Umask == lp.Umask
}
