package exec

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"
)

type Config struct {
	args    []string
	writers []io.WriteCloser
}

func newDefaultConfig() *Config {
	return &Config{
		writers: []io.WriteCloser{os.Stderr},
	}
}

func Cmd(execName string, opts ...OptionFunc) (execCmdErr error) {
	return execCmd(execName, opts...)
}

func execCmd(execName string, opts ...OptionFunc) (execCmdErr error) {
	config := newDefaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	cmd := exec.Command(execName, config.args...)

	terminalReaders, terminalReadersErr := func() ([]*bufio.Reader, error) {
		stdout, stdoutErr := cmd.StdoutPipe()
		if stdoutErr != nil {
			return nil, stdoutErr
		}
		stderr, stderrErr := cmd.StderrPipe()
		if stderrErr != nil {
			return nil, stderrErr
		}
		return []*bufio.Reader{bufio.NewReader(stdout), bufio.NewReader(stderr)}, nil
	}()
	if terminalReadersErr != nil {
		return terminalReadersErr
	}

	var wg sync.WaitGroup

	for _, terminalReader := range terminalReaders {
		terminalReader := terminalReader
		wg.Add(1)
		go func() {
			defer wg.Done()
			execCmdErr = persistOutput(terminalReader, config.writers...)
		}()
	}

	runErr := cmd.Run()
	if runErr != nil {
		return runErr
	}
	wg.Wait()

	return nil
}

func persistOutput(reader *bufio.Reader, writers ...io.WriteCloser) (persistOutputErr error) {
	dump2Writer := func(data []byte) error {
		for _, w := range writers {
			if _, err := w.Write(data); err != nil {
				return err
			}
		}
		return nil
	}

	closeWriter := func() error {
		for _, w := range writers {
			if err := w.Close(); err != nil {
				return err
			}
		}
		return nil
	}

	defer func() {
		if err := closeWriter(); err != nil {
			persistOutputErr = err
		}
	}()

	outputBytes := make([]byte, 256)

	for {
		n, err := reader.Read(outputBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err = dump2Writer(outputBytes[:n]); err != nil {
			return err
		}

	}
	return nil
}
