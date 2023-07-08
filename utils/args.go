package utils

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

// ParseArgs 解析命令行
// type Args struct {
//	LogDir string `arg:"dir" txt:"日志所在文件夹"`
//	Model  int    `arg:"model" txt:"日志模型"`
//	Flow   bool   `arg:"flow" txt:"是否滚动显示"`
// }
// 支持的类型: string, int, int64, uint, uint64, float64, bool
func ParseArgs(ins interface{}) {
	ParseArgsEx(os.Args[1:], ins)
}

func ParseArgsEx(args []string, ins interface{}) {
	argPtr := reflect.ValueOf(ins)
	if false == argPtr.IsValid() ||
		argPtr.Kind() != reflect.Pointer ||
		argPtr.Elem().Kind() != reflect.Struct {
		panic("'ins' paramter error")
	}

	argElem := argPtr.Elem()
	argType := argElem.Type()
	for i := 0; i < argElem.NumField(); i++ {
		tp := argType.Field(i)
		name := tp.Tag.Get("arg")
		usage := tp.Tag.Get("txt")
		if "" == name {
			continue
		}

		field := argElem.Field(i)
		switch tp.Type.Kind() {
		case reflect.String:
			flag.StringVar(field.Addr().Interface().(*string), name, field.String(), usage)
		case reflect.Int:
			flag.IntVar(field.Addr().Interface().(*int), name, int(field.Int()), usage)
		case reflect.Int64:
			flag.Int64Var(field.Addr().Interface().(*int64), name, field.Int(), usage)
		case reflect.Uint:
			flag.UintVar(field.Addr().Interface().(*uint), name, uint(field.Uint()), usage)
		case reflect.Uint64:
			flag.Uint64Var(field.Addr().Interface().(*uint64), name, field.Uint(), usage)
		case reflect.Float64:
			flag.Float64Var(field.Addr().Interface().(*float64), name, field.Float(), usage)
		case reflect.Bool:
			flag.BoolVar(field.Addr().Interface().(*bool), name, field.Bool(), usage)
		default:
			panic(fmt.Sprintf("can't support %s", tp.Type.Name()))
		}
	}

	flag.CommandLine.Parse(args)
}
