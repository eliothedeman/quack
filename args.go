package quack

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
)

// filedNameToArg formats a filed name to it's name as cli arg.
func fieldNameToArg(f string) string {
	return strcase.ToKebab(f)
}

// This is a very long function used to add a flag to a flagset based on a reflection type.
// I don't like it, but I don't see a better way.
func setFlag(v reflect.Value, f reflect.StructField, fs *pflag.FlagSet) {
	t := f.Tag
	help := t.Get("help")
	strVal := t.Get("default")
	intVal, _ := strconv.Atoi(strVal)
	floatVal, _ := strconv.ParseFloat(strVal, 64)
	durationVal, _ := time.ParseDuration(strVal)
	short := t.Get("short")
	hasShort := short != ""
	boolVal := strVal == "true"
	argName := fieldNameToArg(f.Name)
	addr := v.Addr().Interface()
	rawAddr := unsafe.Pointer(v.UnsafeAddr())
	switch v.Kind() {
	case reflect.Bool:
		if hasShort {
			fs.BoolVarP((*bool)(rawAddr), argName, short, boolVal, help)
		} else {
			fs.BoolVar((*bool)(rawAddr), argName, boolVal, help)
		}
	case reflect.Int:
		if hasShort {
			fs.IntVarP((*int)(rawAddr), argName, short, intVal, help)
		} else {
			fs.IntVar((*int)(rawAddr), argName, intVal, help)
		}
		// handle a few types that are also int64
	case reflect.Int64:
		if addr, ok := addr.(*time.Duration); ok {
			if hasShort {
				fs.DurationVarP(addr, argName, short, durationVal, help)
			} else {
				fs.DurationVar(addr, argName, durationVal, help)
			}
		} else {
			if hasShort {
				fs.Int64VarP((*int64)(rawAddr), argName, short, int64(intVal), help)
			} else {
				fs.Int64Var((*int64)(rawAddr), argName, int64(intVal), help)
			}
		}
	case reflect.Int32:
		if hasShort {
			fs.Int32VarP((*int32)(rawAddr), argName, short, int32(intVal), help)
		} else {
			fs.Int32Var((*int32)(rawAddr), argName, int32(intVal), help)
		}
	case reflect.Int16:
		if hasShort {
			fs.Int16VarP((*int16)(rawAddr), argName, short, int16(intVal), help)
		} else {
			fs.Int16Var((*int16)(rawAddr), argName, int16(intVal), help)
		}
	case reflect.Int8:
		if hasShort {
			fs.Int8VarP((*int8)(rawAddr), argName, short, int8(intVal), help)
		} else {
			fs.Int8Var((*int8)(rawAddr), argName, int8(intVal), help)
		}
	case reflect.Uint:
		if hasShort {
			fs.UintVarP((*uint)(rawAddr), argName, short, uint(intVal), help)
		} else {
			fs.UintVar((*uint)(rawAddr), argName, uint(intVal), help)
		}
	case reflect.Uint64:
		if hasShort {
			fs.Uint64VarP((*uint64)(rawAddr), argName, short, uint64(intVal), help)
		} else {
			fs.Uint64Var((*uint64)(rawAddr), argName, uint64(intVal), help)
		}
	case reflect.Uint32:
		if hasShort {
			fs.Uint32VarP((*uint32)(rawAddr), argName, short, uint32(intVal), help)
		} else {
			fs.Uint32Var((*uint32)(rawAddr), argName, uint32(intVal), help)
		}
	case reflect.Uint16:
		if hasShort {
			fs.Uint16VarP((*uint16)(rawAddr), argName, short, uint16(intVal), help)
		} else {
			fs.Uint16Var((*uint16)(rawAddr), argName, uint16(intVal), help)
		}
	case reflect.Uint8:
		if hasShort {
			fs.Uint8VarP((*uint8)(rawAddr), argName, short, uint8(intVal), help)
		} else {
			fs.Uint8Var((*uint8)(rawAddr), argName, uint8(intVal), help)
		}
	case reflect.Float32:
		if hasShort {
			fs.Float32VarP((*float32)(rawAddr), argName, short, float32(floatVal), help)
		} else {
			fs.Float32Var((*float32)(rawAddr), argName, float32(floatVal), help)
		}
	case reflect.Float64:
		if hasShort {
			fs.Float64VarP((*float64)(rawAddr), argName, short, float64(floatVal), help)
		} else {
			fs.Float64Var((*float64)(rawAddr), argName, float64(floatVal), help)
		}
	case reflect.String:
		if hasShort {
			fs.StringVarP((*string)(rawAddr), argName, short, strVal, help)
		} else {
			fs.StringVar((*string)(rawAddr), argName, strVal, help)
		}
	default:
		log.Panicf("Unable to handle type set flags for %v", f)
	}
}

func getFlags(name string, c Command) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	v := reflect.Indirect(reflect.ValueOf(c))
	t := v.Type()
	v.FieldByNameFunc(func(s string) bool {
		f := v.FieldByName(s)
		sf, ok := t.FieldByName(s)
		if !ok {
			panic("wtf")
		}

		// respect the ignore flag
		if val, ok := sf.Tag.Lookup("ignore"); ok {
			if val != "false" {
				return false
			}
		}

		setFlag(f, sf, fs)
		return false
	})
	return fs
}

func fmtHelp(name string, u Unit) string {
	var b strings.Builder

	switch u := u.(type) {
	case Command:
		fmt.Fprintf(&b, "Usage:    %s [args]\n", name)
		f := getFlags("", u)
		b.WriteString(f.FlagUsages())

	case Group:
		fmt.Fprintf(&b, "Usage:    %s <cmd> [args]\n", name)
		if h, ok := u.(Helper); ok {
			b.WriteString(h.Help())
			b.WriteByte('\n')
		}

		cmds := u.SubCommands()
		k := keys(cmds)
		sort.Slice(k, func(i, j int) bool {
			return k[i] < k[j]
		})
		for _, name := range k {
			fmt.Fprintf(&b, "\t%s", name)
			c := cmds[name]
			if h, ok := c.(Helper); ok {
				fmt.Fprintf(&b, " \"%s\"", h.Help())
			}
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func printHelp(name string, u Unit) {
	fmt.Println(fmtHelp(name, u))
}

func hasHelpArg(args []string, shortHelp bool) bool {
	for i, a := range args {
		// if our first arg is not an option, don't look for help in the subsequent flags.
		if i == 0 {
			if !strings.HasPrefix(a, "-") {
				return false
			}
		}
		if a == "--help" {
			return true
		}
		if shortHelp && a == "-h" {
			return true
		}
	}
	return false
}

func helpError(name string, u Unit) error {
	return errors.New(fmtHelp(name, u))
}
