package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

var (
	mode     = flag.Int("mode", 0, "mode must be set\n 0: add debug statements\n 1: remove debug statements\n")
	filename = flag.String("filename", "", "target filename must be set\n")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("goeasydebug: ")
	flag.Usage = Usage
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(2)
	}

	switch *mode {
	case 0:
		processAddDebugStatementsMode()
	case 1:
		processRemoveDebugStatementsMode()
	}
}

func processAddDebugStatementsMode() {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, *filename, nil, 0)
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}

	f, err := ioutil.TempFile("", "goeasydebug.go")
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	for _, d := range file.Decls {
		switch w := d.(type) {
		case *ast.FuncDecl:
			newBodyList := []ast.Stmt{}
			for _, l := range w.Body.List {
				newBodyList = append(newBodyList, l)
				switch a := l.(type) {
				case *ast.AssignStmt:
					for _, lh := range a.Lhs {
						switch i := lh.(type) {
						case *ast.Ident:
							debug := fmt.Sprintf(`dmp("%s",%s)`, string(i.Name), string(i.Name))
							ex, err := parser.ParseExpr(debug)
							if err != nil {
								log.Fatalln("Error:", err)
								return
							}
							switch as := ex.(type) {
							case *ast.CallExpr:
								as.Args[1].(*ast.Ident).NamePos = token.Pos(int(as.Args[0].(*ast.BasicLit).ValuePos) + len(i.Name))
							}
							debugStmt := &ast.ExprStmt{
								X: ex,
							}
							newBodyList = append(newBodyList, debugStmt)
						}
					}
				}
			}
			w.Body.List = newBodyList
		}
	}

	if err := format.Node(writer, fset, file); err != nil {
		log.Fatalln("Error:", err)
		return
	}

	goeasydebugDefinition := `
// generated from goeasydebug
// function for data dump
func dmp(valueName string, v ...interface{}) {
  for _, vv := range(v) {
      fmt.Printf("%s: %#v\n",valueName, vv)
  }
}
`

	writer.WriteString(goeasydebugDefinition)

	writer.Flush()

	if err := os.Rename(f.Name(), *filename); err != nil {
		log.Fatalln("Error:", err)
		return
	}

	return
}

func processRemoveDebugStatementsMode() {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, *filename, nil, 0)
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}

	f, err := ioutil.TempFile("", "goeasydebug.go")
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	newDeclList := []ast.Decl{}
	for _, d := range file.Decls {
		switch w := d.(type) {
		case *ast.FuncDecl:
			if w.Name.Name == "dmp" {
				continue
			}
			newBodyList := []ast.Stmt{}
			for _, l := range w.Body.List {
				newBodyList = append(newBodyList, l)
				switch a := l.(type) {
				case *ast.ExprStmt:
					switch x := a.X.(type) {
					case *ast.CallExpr:
						switch f := x.Fun.(type) {
						case *ast.Ident:
							if f.Name == "dmp" {
								newBodyList = newBodyList[:len(newBodyList)-1]
							}
						}
					}
				}
			}
			w.Body.List = newBodyList
		}
		newDeclList = append(newDeclList, d)
	}

	file.Decls = newDeclList

	if err := format.Node(writer, fset, file); err != nil {
		log.Fatalln("Error:", err)
		return
	}

	writer.Flush()

	if err := os.Rename(f.Name(), *filename); err != nil {
		log.Fatalln("Error:", err)
		return
	}

	return
}
