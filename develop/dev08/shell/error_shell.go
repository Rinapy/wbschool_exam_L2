package shell

type UnknownCommand struct{}

func (c *UnknownCommand) Error() string {
	return "unknown command."
}

type NoArgumentsPassed struct{}

func (c *NoArgumentsPassed) Error() string {
	return "no arguments passed."
}

type PorcKillError struct{}

func (c *PorcKillError) Error() string {
	return "filed to kill process"
}
