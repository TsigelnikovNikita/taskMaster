package program

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
	Command          string   `ini:"command"`
	ProcessNumber    int      `ini:"process_number"`
	AutoStart        bool     `ini:"auto_start"`
	StartTimeSec     int      `ini:"start_secs"`
	StartRetries     int      `ini:"start_retries"`
	AutoRestart      string   `ini:"autoRestart"`
	ExitCodes        []int    `ini:"exit_codes"`
	StopSignal       int      `ini:"stop_signal"`
	StopWaitSecs     int      `ini:"stop_wait_secs"`
	StdErrLogfile    string   `ini:"log_file_error"`
	StdOutLogfile    string   `ini:"log_file"`
	Environments     []string `ini:"env"`
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
