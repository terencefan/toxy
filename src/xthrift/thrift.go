package xthrift

import (
	tp "github.com/stdrickforce/go-thrift/parser"
)

type Thrift struct {
	Typedefs   map[string]*tp.Typedef
	Constants  map[string]*tp.Constant
	Enums      map[string]*tp.Enum
	Structs    map[string]*tp.Struct
	Exceptions map[string]*tp.Struct
	Unions     map[string]*tp.Struct
	Services   map[string]*tp.Service
}

var parser = &tp.Parser{}

func Parse(filename string) (thrift *Thrift, files []string, err error) {
	thrift = &Thrift{
		Typedefs:   make(map[string]*tp.Typedef),
		Constants:  make(map[string]*tp.Constant),
		Enums:      make(map[string]*tp.Enum),
		Structs:    make(map[string]*tp.Struct),
		Exceptions: make(map[string]*tp.Struct),
		Unions:     make(map[string]*tp.Struct),
		Services:   make(map[string]*tp.Service),
	}

	parsedThrift, _, err := parser.ParseFile(filename)
	if err != nil {
		return
	}

	for file, t := range parsedThrift {
		files = append(files, file)
		for k, v := range t.Typedefs {
			thrift.Typedefs[k] = v
		}
		for k, v := range t.Constants {
			thrift.Constants[k] = v
		}
		for k, v := range t.Enums {
			thrift.Enums[k] = v
		}
		for k, v := range t.Structs {
			thrift.Structs[k] = v
		}
		for k, v := range t.Exceptions {
			thrift.Exceptions[k] = v
		}
		for k, v := range t.Unions {
			thrift.Unions[k] = v
		}
		for k, v := range t.Services {
			thrift.Services[k] = v
		}
	}
	// fmt.Println(thrift)
	return
}
