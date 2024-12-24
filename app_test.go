package quack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type argCount struct {
	args     int
	opts     int
	flags    int
	reqArgs  int
	reqOpts  int
	reqFlags int
}

func toCount(f *funcDesc) argCount {
	return argCount{
		args:     len(f.args),
		opts:     len(f.opts),
		flags:    len(f.flags),
		reqArgs:  len(f.reqArgs),
		reqOpts:  len(f.reqOpts),
		reqFlags: len(f.reqFlags),
	}
}

func TestExtractFuncDefs(t *testing.T) {
	tf, err := extract(func(args struct {
		Time int
	}) {
	})
	if err != nil {
		t.Fatalf("Err during extraction %s", err)
	}
	assert.Equal(t, argCount{}, toCount(tf))
}

func TestExtractOpt(t *testing.T) {
	tf, err := extract(func(args struct {
		Time Opt[int]
	}) {
	})
	if err != nil {
		t.Fatalf("Err during extraction %s", err)
	}
	assert.Equal(t, argCount{
		opts: 1,
	}, toCount(tf))
}

func TestExtractArg(t *testing.T) {
	tf, err := extract(func(args struct {
		Time Arg[int]
	}) {
	})
	if err != nil {
		t.Fatalf("Err during extraction %s", err)
	}
	assert.Equal(t, argCount{
		args: 1,
	}, toCount(tf))
}

func TestExtractFlag(t *testing.T) {
	tf, err := extract(func(args struct {
		IsSet Flag
	}) {
	})
	if err != nil {
		t.Fatalf("Err during extraction %s", err)
	}
	assert.Equal(t, argCount{
		flags: 1,
	}, toCount(tf))
}

/*
func TestExtractRequired(t *testing.T) {
	tf, err := extract(func(args struct {
		IsSet Req[Opt[int]]
	}) {
	})
	if err != nil {
		t.Fatalf("Err during extraction %s", err)
	}
	assert.Equal(t, argCount{
		reqOpts: 1,
	}, toCount(tf))
}

func TestCall(t *testing.T) {
}
*/
