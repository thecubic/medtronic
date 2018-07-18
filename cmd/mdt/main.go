package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/thecubic/medtronic"
)

type (
	// Printer represents a function that prints an arbitrary value.
	Printer func(interface{})
)

var (
	formatFlag = flag.String("f", "openaps", "print result in specified `format`")

	format = map[string]Printer{
		"internal": showInternal,
		"json":     showJSON,
		"openaps":  showOpenAPS,
	}

	openAPSMode bool
)

func usage() {
	eprintf("usage: %s [options] command [ arg ... ]\n", os.Args[0])
	eprintf("   or: %s [options] command [ args.json ]\n", os.Args[0])
	flag.PrintDefaults()
	fmts := ""
	for k := range format {
		fmts += " " + k
	}
	eprintf("output formats:%s\n", fmts)
	keys := make([]string, len(command))
	i := 0
	for k := range command {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	cmds := ""
	for _, k := range keys {
		cmds += " " + k
	}
	eprintf("commands:%s\n", cmds)
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	printFn := format[*formatFlag]
	if printFn == nil {
		eprintf("%s: unknown format\n", *formatFlag)
		usage()
	}
	openAPSMode = *formatFlag == "openaps"
	if flag.NArg() == 0 {
		usage()
	}
	name := flag.Arg(0)
	cmd, found := command[name]
	if !found {
		eprintf("%s: unknown command\n", name)
		usage()
	}
	args := getArgs(name, cmd)
	pump := medtronic.Open()
	exitOnError(pump)
	defer pump.Close()
	pump.Wakeup()
	exitOnError(pump)
	result := cmd.Cmd(pump, args)
	exitOnError(pump)
	if result == nil {
		return
	}
	printFn(result)
}

func exitOnError(pump *medtronic.Pump) {
	err := pump.Error()
	if err == nil {
		return
	}
	if pump.NoResponse() {
		log.Print(err)
		os.Exit(2)
	}
	log.Fatal(err)
}

type (
	// Arguments represents the formal and actual parameters for a command.
	Arguments map[string]interface{}
)

// String returns the string value associated with the given key.
func (args Arguments) String(key string) (string, error) {
	arg := args[key]
	s, ok := arg.(string)
	if !ok {
		return s, fmt.Errorf("%q argument must be a string", key)
	}
	return s, nil
}

// Float returns the float64 value associated with the given key.
func (args Arguments) Float(key string) (float64, error) {
	arg := args[key]
	if openAPSMode {
		f, ok := arg.(float64)
		if !ok {
			return f, fmt.Errorf("%q parameter must be a number", key)
		}
		return f, nil
	}
	return strconv.ParseFloat(arg.(string), 64)
}

// Int returns the int value associated with the given key.
func (args Arguments) Int(key string) (int, error) {
	arg := args[key]
	if openAPSMode {
		f, ok := arg.(float64)
		if !ok {
			return int(f), fmt.Errorf("%q argument must be a number", key)
		}
		return int(f), nil
	}
	return strconv.Atoi(arg.(string))
}

// Strings returns the []string value associated with the given key.
func (args Arguments) Strings(key string) ([]string, error) {
	arg := args[key]
	if openAPSMode {
		v, ok := arg.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%q argument must be an array", key)
		}
		a := make([]string, len(v))
		for i, si := range v {
			s, ok := si.(string)
			if !ok {
				return nil, fmt.Errorf("%q argument must be a list of strings", key)
			}
			a[i] = s
		}
		return a, nil
	}
	return arg.([]string), nil
}

func getArgs(name string, cmd Command) Arguments {
	params := cmd.Params
	argv := flag.Args()[1:]
	if len(params) == 0 {
		if len(argv) != 0 {
			log.Fatalf("%s does not take any arguments", name)
		}
		return nil
	}
	if openAPSMode {
		return openAPSArgs(name, params, argv, cmd.Variadic)
	}
	return cliArgs(name, params, argv, cmd.Variadic)
}

// Parse an openaps JSON file for arguments.
func openAPSArgs(name string, params []string, argv []string, variadic bool) Arguments {
	if len(argv) != 1 || !strings.HasSuffix(argv[0], ".json") {
		log.Fatalf("%s: openaps format requires single JSON argument file", name)
	}
	// Unmarshal the JSON argument file.
	file := argv[0]
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("%s: %v", name, err)
	}
	args := make(Arguments)
	err = json.NewDecoder(f).Decode(&args)
	if err != nil {
		log.Fatalf("%s: %v", name, err)
	}
	_ = f.Close()
	// Check that all parameters are present.
	for _, k := range params {
		_, present := args[k]
		if !present {
			log.Fatalf("%s: argument file %s is missing %q parameter", name, file, k)
		}
	}
	return args
}

// Collect command-line arguments.
func cliArgs(name string, params []string, argv []string, variadic bool) Arguments {
	if !variadic && len(argv) != len(params) {
		var p string
		if len(params) != 1 {
			p = "s"
		}
		log.Fatalf("%s requires %d argument%s", name, len(params), p)
	}
	args := make(Arguments)
	for i, k := range params {
		if variadic && i == len(params)-1 {
			// Bind all remaining args to this parameter.
			args[k] = argv[i:]
			break
		}
		if i >= len(argv) {
			// Bind remaining parameters to "".
			args[k] = ""
			continue
		}
		args[k] = argv[i]
	}
	return args
}

func eprintf(format string, arg ...interface{}) {
	fmt.Fprintf(os.Stderr, format, arg...)
}
