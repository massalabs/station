package runner

type Runner interface {
	// Run runs the given command and returns the combined output of stdout and stderr.
	Run(args ...string) error
}
