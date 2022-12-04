package quack

import (
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

const (
	helpTag    = "help"
	defaultTag = "default"
	shortTag   = "short"
	longTag    = "long"
	ignoreTag  = "ignore"
)

// filedNameToArg formats a filed name to it's name as cli arg.
func fieldNameToArg(f string) string {
	return strcase.ToKebab(f)
}

func rawAddr[T any](v reflect.Value) *T {
	return (*T)(unsafe.Pointer(v.UnsafeAddr()))
}

// This is a very long function used to add a flag to a flagset based on a reflection type.
// I don't like it, but I don't see a better way.
func setFlag(v reflect.Value, f reflect.StructField, fs *pflag.FlagSet) {
	t := f.Tag
	help := t.Get(helpTag)
	strVal := t.Get(defaultTag)
	intVal, _ := strconv.Atoi(strVal)
	floatVal, _ := strconv.ParseFloat(strVal, 64)
	durationVal, _ := time.ParseDuration(strVal)
	short := t.Get(shortTag)
	hasShort := short != ""
	boolVal := strVal == "true"
	argName := t.Get(longTag)
	if argName == "" {
		argName = fieldNameToArg(f.Name)
	}
	addr := v.Addr().Interface()
	switch v.Kind() {
	case reflect.Bool:
		if hasShort {
			fs.BoolVarP(rawAddr[bool](v), argName, short, boolVal, help)
		} else {
			fs.BoolVar(rawAddr[bool](v), argName, boolVal, help)
		}
	case reflect.Int:
		if hasShort {
			fs.IntVarP(rawAddr[int](v), argName, short, intVal, help)
		} else {
			fs.IntVar(rawAddr[int](v), argName, intVal, help)
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
				fs.Int64VarP(rawAddr[int64](v), argName, short, int64(intVal), help)
			} else {
				fs.Int64Var(rawAddr[int64](v), argName, int64(intVal), help)
			}
		}
	case reflect.Int32:
		if hasShort {
			fs.Int32VarP(rawAddr[int32](v), argName, short, int32(intVal), help)
		} else {
			fs.Int32Var(rawAddr[int32](v), argName, int32(intVal), help)
		}
	case reflect.Int16:
		if hasShort {
			fs.Int16VarP(rawAddr[int16](v), argName, short, int16(intVal), help)
		} else {
			fs.Int16Var(rawAddr[int16](v), argName, int16(intVal), help)
		}
	case reflect.Int8:
		if hasShort {
			fs.Int8VarP(rawAddr[int8](v), argName, short, int8(intVal), help)
		} else {
			fs.Int8Var(rawAddr[int8](v), argName, int8(intVal), help)
		}
	case reflect.Uint:
		if hasShort {
			fs.UintVarP(rawAddr[uint](v), argName, short, uint(intVal), help)
		} else {
			fs.UintVar(rawAddr[uint](v), argName, uint(intVal), help)
		}
	case reflect.Uint64:
		if hasShort {
			fs.Uint64VarP(rawAddr[uint64](v), argName, short, uint64(intVal), help)
		} else {
			fs.Uint64Var(rawAddr[uint64](v), argName, uint64(intVal), help)
		}
	case reflect.Uint32:
		if hasShort {
			fs.Uint32VarP(rawAddr[uint32](v), argName, short, uint32(intVal), help)
		} else {
			fs.Uint32Var(rawAddr[uint32](v), argName, uint32(intVal), help)
		}
	case reflect.Uint16:
		if hasShort {
			fs.Uint16VarP(rawAddr[uint16](v), argName, short, uint16(intVal), help)
		} else {
			fs.Uint16Var(rawAddr[uint16](v), argName, uint16(intVal), help)
		}
	case reflect.Uint8:
		if hasShort {
			fs.Uint8VarP(rawAddr[uint8](v), argName, short, uint8(intVal), help)
		} else {
			fs.Uint8Var(rawAddr[uint8](v), argName, uint8(intVal), help)
		}
	case reflect.Float32:
		if hasShort {
			fs.Float32VarP(rawAddr[float32](v), argName, short, float32(floatVal), help)
		} else {
			fs.Float32Var(rawAddr[float32](v), argName, float32(floatVal), help)
		}
	case reflect.Float64:
		if hasShort {
			fs.Float64VarP(rawAddr[float64](v), argName, short, float64(floatVal), help)
		} else {
			fs.Float64Var(rawAddr[float64](v), argName, float64(floatVal), help)
		}
	case reflect.String:
		if hasShort {
			fs.StringVarP(rawAddr[string](v), argName, short, strVal, help)
		} else {
			fs.StringVar(rawAddr[string](v), argName, strVal, help)
		}

	default:
		log.Panicf("Unable to handle type set flags for %v", f)
	}
}

func getFlags(name string, c any) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	v := reflect.Indirect(reflect.ValueOf(c))
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return fs
	}
	v.FieldByNameFunc(func(s string) bool {
		f := v.FieldByName(s)
		sf, ok := t.FieldByName(s)
		// Skip over embedded structs from other packages. Their fields will come later in the traversal.
		if sf.Anonymous {
			return false
		}
		if !ok {
			panic("wtf")
		}

		// check if unexported
		if sf.PkgPath != "" {
			return false
		}

		// respect the ignore flag
		if val, ok := sf.Tag.Lookup(ignoreTag); ok {
			if val != "false" {
				return false
			}
		}

		setFlag(f, sf, fs)
		return false
	})
	return fs
}

func fmtHelp(name string, u any) string {
	var b strings.Builder

	switch u := u.(type) {
	case Command:
		fmt.Fprintf(&b, "Usage:    %s [args]\n", name)
		f := getFlags(name, u)
		if h, ok := u.(Helper); ok {
			b.WriteByte('\t')
			b.WriteString(h.Help())
			b.WriteByte('\n')
		}
		fmtUsage(&b, f)

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
			fmt.Fprintf(&b, "      %s ", name)
			c := cmds[name]
			if h, ok := c.(Helper); ok {
				fmt.Fprint(&b, h.Help())
			}
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func printHelp(name string, u any) {
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

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
