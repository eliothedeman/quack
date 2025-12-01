package quack

import (
	"fmt"
	"os"

	"github.com/eliothedeman/check"
)

type ExistingFilePath string

func (e ExistingFilePath) Validate() error {
	st, err := os.Stat(string(e))
	if err != nil {
		return fmt.Errorf("%w not able to find file %s", err, e)
	}
	if st.IsDir() {
		return fmt.Errorf("%s is a directory. expected to be file", e)
	}
	return nil
}

func (e *ExistingFilePath) Open() *os.File {
	return check.Must(os.Open(string(*e)))
}

func (e *ExistingFilePath) OpenWith(flags int, mode os.FileMode) *os.File {
	return check.Must(os.OpenFile(string(*e), flags, mode))
}
