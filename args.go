package quack

import (
	"log"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
)

// filedNameToArg formats a filed name to it's name as cli arg.
func fieldNameToArg(f string) string {
	return strcase.ToKebab(f)
}

func rawAddr[T any](v reflect.Value) *T {
	return (*T)(unsafe.Pointer(v.UnsafeAddr()))
}
func (o *option) setFlag(fs *pflag.FlagSet) {
	if o.Ignore {
		return
	}
	addr := o.Target.Addr().Interface()
	hasShort := o.Short != ""
	short := o.Short
	strVal := o.Default
	intVal, _ := strconv.Atoi(strVal)
	floatVal, _ := strconv.ParseFloat(strVal, 64)
	durationVal, _ := time.ParseDuration(strVal)
	boolVal := strVal == "true"
	argName := o.Name
	help := o.Help
	v := o.Target

	switch o.Target.Kind() {
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
			fs.Uint64VarP(
				rawAddr[uint64](v),
				argName,
				short,
				uint64(intVal),
				help,
			)
		} else {
			fs.Uint64Var(rawAddr[uint64](v), argName, uint64(intVal), help)
		}
	case reflect.Uint32:
		if hasShort {
			fs.Uint32VarP(
				rawAddr[uint32](v),
				argName,
				short,
				uint32(intVal),
				help,
			)
		} else {
			fs.Uint32Var(rawAddr[uint32](v), argName, uint32(intVal), help)
		}
	case reflect.Uint16:
		if hasShort {
			fs.Uint16VarP(
				rawAddr[uint16](v),
				argName,
				short,
				uint16(intVal),
				help,
			)
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
			fs.Float32VarP(
				rawAddr[float32](v),
				argName,
				short,
				float32(floatVal),
				help,
			)
		} else {
			fs.Float32Var(rawAddr[float32](v), argName, float32(floatVal), help)
		}
	case reflect.Float64:
		if hasShort {
			fs.Float64VarP(
				rawAddr[float64](v),
				argName,
				short,
				float64(floatVal),
				help,
			)
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
		log.Panicf("Unable to handle type set flags for %v", o.Target)
	}
}
