package quack

import (
	"testing"
	"time"
)

type testFmtStruct struct {
	Aaa      int
	TtString string
}

func (*testFmtStruct) Run([]string) {

}

type testFmtSub struct {
}

func (*testFmtSub) SubCommands() map[string]Unit {
	return map[string]Unit{
		"testin": &testFmtStruct{},
		"test2":  &testFmtStruct{},
	}
}

func TestFmtHelpGroup(t *testing.T) {
	fh := fmtHelp("test", new(testFmtSub))
	want := `Usage:    test <cmd> [args]
	test2
	testin
`

	if fh != want {
		t.Errorf("\n%x\n!=\n%x", fh, want)
	}
}

type allTypes struct {
	Int      int
	Int8     int8
	Int16    int16
	Int32    int32
	Int64    int64
	Uint     uint
	Uint8    uint8
	Uint16   uint16
	Uint32   uint32
	Uint64   uint64
	Float32  float32
	Float64  float64
	Bool     bool
	String   string
	Duration time.Duration
}

func (allTypes) Run([]string) {

}

func TestAllTypes(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want allTypes
	}{
		{
			"None set, no args",
			[]string{},
			allTypes{},
		},
		{
			"Bool set",
			[]string{"--bool"},
			allTypes{
				Bool: true,
			},
		},
		{
			"int set",
			[]string{"--int", "100"},
			allTypes{
				Int: 100,
			},
		},
		{
			"floats set",
			[]string{"--float-32", "3.2", "--float-64", "6.4"},
			allTypes{
				Float32: 3.2,
				Float64: 6.4,
			},
		},
		{
			"duration",
			[]string{"--duration", "4h"},
			allTypes{
				Duration: time.Hour * 4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(allTypes)
			fs := getFlags(tt.name, got)
			err := fs.Parse(tt.args)
			if err != nil {
				t.Errorf("getFlags() error %s", err)
			}

			if *got != tt.want {
				t.Errorf("hasHelpArg() = \n%+v want \n%+v", *got, tt.want)
			}
		})
	}
}
