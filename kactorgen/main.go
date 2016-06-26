package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	actor_namespace string = "koalanet.Actor"
	pkg_name        string = ""
	actor_list      map[string]*AstStruct
	dir_path        string = "."
)

type AstMethod struct {
	Name string

	args  string
	reply string
}

type AstStruct struct {
	Name    string
	Methods []AstMethod
}

func parseStructFields(fields []*ast.Field) bool {
	field := fields[0]

	var expr *ast.SelectorExpr

	switch typ := field.Type.(type) {
	case (*ast.SelectorExpr):
		expr = typ
		break
	default:
		return false
	}

	nsIdent := expr.X.(*ast.Ident)
	if nsIdent == nil {
		return false
	}

	nsName := nsIdent.Name + "." + expr.Sel.Name
	if nsName != actor_namespace {
		return false
	}

	return true
}

func parseStruct(decl *ast.GenDecl) {
	if decl.Tok != token.TYPE {
		return
	}

	structSpec := decl.Specs[0].(*ast.TypeSpec)

	typ := structSpec.Type.(*ast.StructType)
	if typ == nil {
		return
	}

	if typ.Fields == nil {
		return
	}

	if len(typ.Fields.List) <= 0 {
		return
	}

	if !parseStructFields(typ.Fields.List) {
		return
	}

	actorStruct := &AstStruct{}
	actorStruct.Name = structSpec.Name.Name
	actor_list[actorStruct.Name] = actorStruct
	log.Printf("find class: %s.", actorStruct.Name)
}

func parseFunc(decl *ast.FuncDecl) {
	if decl.Recv == nil {
		return
	}

	if len(decl.Recv.List) <= 0 {
		return
	}

	thisField := decl.Recv.List[0]

	classField := thisField.Type.(*ast.StarExpr)
	if classField == nil {
		return
	}

	classIdent := classField.X.(*ast.Ident)
	if classIdent == nil {
		return
	}

	// 判断首字母大写
	c := decl.Name.Name[0]
	if c < 65 || c > 90 {
		return
	}

	args, reply, ok := parseMethodArgs(decl.Type)
	if !ok {
		return
	}

	className := classIdent.Name
	methodName := decl.Name.Name
	astActor, ok := actor_list[className]

	if !ok {
		astActor = &AstStruct{className, make([]AstMethod, 0)}
		actor_list[className] = astActor
	}

	astActor.Methods = append(astActor.Methods, AstMethod{methodName, args, reply})

	log.Printf("find method: %s:%s(%s %s) error", className, methodName, args, reply)
}

func parseMethodArgs(typ *ast.FuncType) (args string, reply string, ok bool) {
	ok = false
	if typ.Params == nil {
		return
	}

	for _, v := range typ.Params.List {
		if v.Names[0].Name == "args" {
			args = v.Names[0].Name

			argType := v.Type.(*ast.Ident)
			if argType == nil {
				return
			}

			args += " "
			args += argType.Name

			// log.Printf(args)

			continue
		}

		if v.Names[0].Name == "reply" {
			reply = v.Names[0].Name

			argType := v.Type.(*ast.StarExpr)
			if argType == nil {
				return
			}

			reply += " *"
			reply += argType.X.(*ast.Ident).Name

			// log.Printf(reply)
		}
	}

	ok = true

	return
}

func parseFile(fileName string) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		log.Println(err)
		return
	}

	pkg_name = f.Name.Name
	log.Printf("package:%s", pkg_name)

	for _, v := range f.Decls {
		switch t := v.(type) {
		case *ast.GenDecl:
			parseStruct(t)
			break
		case *ast.FuncDecl:
			parseFunc(t)
			break
		}
	}
}

func nameFilter(filename string) bool {
	if filename == "actor_wrap.go" {
		return false
	}

	ext := filepath.Ext(filename)
	return ext == ".go"
}

func dirFilter(f os.FileInfo) bool { return nameFilter(f.Name()) }

func parseDir() {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dir_path, dirFilter, 0)
	if err != nil {
		return
	}

	for _, pkg := range pkgs {
		for fname := range pkg.Files {
			parseFile(fname)
		}
	}
}

var impl_temp_def string = `package %s

import "github.com/taodev/koalanet"
%s
`

var impl_struct_def string = `
type %sImpl struct {
	%s
}
`

var impl_method_def string = `
func (actor *%sImpl) %sWrap(args interface{}, reply interface{}) error {
	return actor.%s(%s)
}
`

var impl_init_def string = `
func init() {%s}`

var impl_reg_head string = `
	koalanet.RegActor("%s", func() koalanet.IActor {
		actor := &%sImpl{}
		actor.InitActor()`

var impl_reg_method string = `
		actor.RegMethod("%s", actor.%sWrap)`

var impl_reg_end string = `
		return actor
	})
`

var wrap_struct_def string = `
type %sWrap struct {
	Handle uint32
}
`

var wrap_method_def string = `
func (actor *%sWrap) %s(isSync bool%s) error {%s}
`

var wrap_method_code string = `
	if isSync {
		return koalanet.Call(actor.Handle, "%s", %s, %s)
	}
	
	return koalanet.Send(actor.Handle, "%s", %s)
`

func genWrap() string {
	body := ""

	for _, actor := range actor_list {
		body += fmt.Sprintf(wrap_struct_def, actor.Name)

		for _, method := range actor.Methods {
			args := ""
			args1 := "nil"

			if len(method.args) > 0 {
				args += ", "
				args += method.args

				args1 = "args"
			}

			reply1 := "nil"
			if len(method.reply) > 0 {
				args += ", "
				args += method.reply

				reply1 = "reply"
			}

			code := fmt.Sprintf(wrap_method_code, method.Name, args1, reply1, method.Name, args1)

			body += fmt.Sprintf(wrap_method_def, actor.Name, method.Name, args, code)
		}
	}

	return body
}

func genImpl() string {
	body := fmt.Sprintf(impl_temp_def, pkg_name, genWrap())

	reg_code := ""

	for _, actor := range actor_list {
		body += fmt.Sprintf(impl_struct_def, actor.Name, actor.Name)
		reg_code += fmt.Sprintf(impl_reg_head, actor.Name, actor.Name)

		for _, method := range actor.Methods {
			args := ""
			if len(method.args) > 0 {
				typ := strings.Split(method.args, " ")[1]
				args += fmt.Sprintf("args.(%s)", typ)
			}

			if len(method.reply) > 0 {
				if len(method.args) > 0 {
					args += ", "
				}
				typ := strings.Split(method.reply, " ")[1]
				args += fmt.Sprintf("reply.(%s)", typ)
			}

			body += fmt.Sprintf(impl_method_def, actor.Name, method.Name, method.Name, args)

			reg_code += fmt.Sprintf(impl_reg_method, method.Name, method.Name)
		}

		reg_code += impl_reg_end
	}

	body += fmt.Sprintf(impl_init_def, reg_code)

	// fmt.Println(body)

	ioutil.WriteFile(dir_path+"/actor_wrap.go", []byte(body), os.ModeAppend)

	return body
}

func main() {
	actor_list = make(map[string]*AstStruct)

	if len(os.Args) > 1 {
		dir_path = os.Args[1]
	}

	parseDir()

	genImpl()
}
