package parser

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// PackageDoc holds parsed documentation for a single Go package
type PackageDoc struct {
	Name       string
	ImportPath string
	Doc        string
	Structs    []StructDoc
	Interfaces []InterfaceDoc
	Functions  []FunctionDoc
	Types      []TypeDoc
	Constants  []ConstDoc
	Variables  []VarDoc
	Examples   []ExampleDoc
	Files      []string
}

// StructDoc holds documentation for a struct
type StructDoc struct {
	Name    string
	Doc     string
	Fields  []FieldDoc
	Methods []FunctionDoc
}

// FieldDoc holds documentation for a struct field
type FieldDoc struct {
	Name string
	Type string
	Tag  string
	Doc  string
}

// InterfaceDoc holds documentation for an interface
type InterfaceDoc struct {
	Name    string
	Doc     string
	Methods []FunctionDoc
}

// FunctionDoc holds documentation for a function or method
type FunctionDoc struct {
	Name       string
	Doc        string
	Signature  string
	Receiver   string
	Params     []ParamDoc
	Returns    []string
	IsExported bool
}

// ParamDoc holds documentation for a function parameter
type ParamDoc struct {
	Name string
	Type string
}

// TypeDoc holds documentation for a type alias or definition
type TypeDoc struct {
	Name       string
	Doc        string
	Underlying string
}

// ConstDoc holds documentation for constants
type ConstDoc struct {
	Name  string
	Doc   string
	Value string
	Type  string
}

// VarDoc holds documentation for variables
type VarDoc struct {
	Name  string
	Doc   string
	Type  string
	Value string
}

// ExampleDoc holds an example
type ExampleDoc struct {
	Name   string
	Doc    string
	Code   string
	Output string
}

// ParsePackage parses all Go files in a directory and returns documentation.
// It walks the ENTIRE directory tree — no depth limit, no package count limit.
// Every sub-package is parsed so nothing is missed.
func ParsePackage(srcDir string, importPath string) ([]*PackageDoc, error) {
	var allDocs []*PackageDoc

	// Directories to skip — only truly useless folders
	skipDirs := map[string]bool{
		"vendor":      true,
		"testdata":    true,
		"third_party": true,
		"hack":        true,
		"_output":     true,
	}

	// Walk the ENTIRE tree using filepath.WalkDir — no depth limit
	_ = filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors, keep walking
		}

		// Only process directories
		if !d.IsDir() {
			return nil
		}

		name := d.Name()

		// Skip hidden directories and known junk
		if strings.HasPrefix(name, ".") || skipDirs[name] {
			return filepath.SkipDir
		}

		// Compute import path for this directory
		rel, relErr := filepath.Rel(srcDir, path)
		if relErr != nil {
			return nil
		}
		subImport := importPath
		if rel != "." {
			subImport = importPath + "/" + filepath.ToSlash(rel)
		}

		// Parse ALL Go packages in this directory (a dir can have multiple packages)
		pkgDocs := parseDirAll(path, subImport)
		allDocs = append(allDocs, pkgDocs...)

		return nil
	})

	if len(allDocs) == 0 {
		return nil, fmt.Errorf("no Go packages found in %s", srcDir)
	}

	// Sort: most useful packages first (more exported items + priority keywords)
	sortPackagesByImportance(allDocs)

	return allDocs, nil
}

// sortPackagesByImportance sorts packages so the most useful ones appear first
func sortPackagesByImportance(docs []*PackageDoc) {
	priorityKeywords := []string{
		"kubernetes", "clientcmd", "clientset", "rest",
		"informers", "cache", "tools", "scheme",
		"v1", "v1beta1", "v1alpha1",
		"core", "apps", "batch", "policy", "networking",
		"client", "config", "util", "runtime",
	}

	for i := 0; i < len(docs); i++ {
		for j := i + 1; j < len(docs); j++ {
			if packageScore(docs[j], priorityKeywords) > packageScore(docs[i], priorityKeywords) {
				docs[i], docs[j] = docs[j], docs[i]
			}
		}
	}
}

// packageScore calculates an importance score for a package
func packageScore(pkg *PackageDoc, priorityKeywords []string) int {
	score := len(pkg.Structs)*3 + len(pkg.Functions)*2 + len(pkg.Interfaces)*3 + len(pkg.Types) + len(pkg.Constants)

	lower := strings.ToLower(pkg.ImportPath)
	for _, kw := range priorityKeywords {
		if strings.Contains(lower, kw) {
			score += 50
		}
	}

	if pkg.Doc != "" {
		score += 10
	}

	return score
}

// parseDirAll parses a single directory and returns ALL Go packages found in it.
// A directory can contain multiple packages — this returns every one of them.
func parseDirAll(dir string, importPath string) []*PackageDoc {
	fset := token.NewFileSet()

	// Parse all Go files in the directory
	pkgs, err := goparser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		// Skip test files
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, goparser.ParseComments)
	if err != nil {
		return nil
	}

	var results []*PackageDoc

	for pkgName, pkg := range pkgs {
		// Skip test packages
		if strings.HasSuffix(pkgName, "_test") {
			continue
		}

		doc := &PackageDoc{
			Name:       pkgName,
			ImportPath: importPath,
		}

		// Collect files
		for fileName := range pkg.Files {
			doc.Files = append(doc.Files, filepath.Base(fileName))
		}

		// Parse each file
		for _, file := range pkg.Files {
			// Package doc
			if file.Doc != nil && doc.Doc == "" {
				doc.Doc = file.Doc.Text()
			}

			// Walk the AST
			ast.Inspect(file, func(n ast.Node) bool {
				switch decl := n.(type) {
				case *ast.GenDecl:
					parseGenDecl(decl, doc, fset)
				case *ast.FuncDecl:
					parseFuncDecl(decl, doc, fset)
				}
				return true
			})
		}

		// Only include if it has any exported content
		if len(doc.Structs) > 0 || len(doc.Functions) > 0 || len(doc.Interfaces) > 0 ||
			len(doc.Types) > 0 || len(doc.Constants) > 0 || len(doc.Variables) > 0 {
			results = append(results, doc)
		}
	}

	return results
}

// parseGenDecl parses type, const, and var declarations
func parseGenDecl(decl *ast.GenDecl, doc *PackageDoc, fset *token.FileSet) {
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if !s.Name.IsExported() {
				continue
			}
			docText := ""
			if decl.Doc != nil {
				docText = decl.Doc.Text()
			}
			if s.Comment != nil && docText == "" {
				docText = s.Comment.Text()
			}

			switch t := s.Type.(type) {
			case *ast.StructType:
				structDoc := StructDoc{
					Name: s.Name.Name,
					Doc:  docText,
				}
				if t.Fields != nil {
					for _, field := range t.Fields.List {
						fd := FieldDoc{
							Type: exprToString(field.Type),
						}
						if len(field.Names) > 0 {
							fd.Name = field.Names[0].Name
						} else {
							// Embedded/anonymous field — use type name as field name
							fd.Name = embeddedFieldName(field.Type)
						}
						if field.Tag != nil {
							fd.Tag = field.Tag.Value
						}
						// Capture doc from both Doc (above) and Comment (inline)
						if field.Doc != nil {
							fd.Doc = field.Doc.Text()
						} else if field.Comment != nil {
							fd.Doc = field.Comment.Text()
						}
						structDoc.Fields = append(structDoc.Fields, fd)
					}
				}
				doc.Structs = append(doc.Structs, structDoc)

			case *ast.InterfaceType:
				ifaceDoc := InterfaceDoc{
					Name: s.Name.Name,
					Doc:  docText,
				}
				if t.Methods != nil {
					for _, method := range t.Methods.List {
						if len(method.Names) > 0 {
							mDoc := ""
							if method.Doc != nil {
								mDoc = method.Doc.Text()
							}
							ifaceDoc.Methods = append(ifaceDoc.Methods, FunctionDoc{
								Name:      method.Names[0].Name,
								Doc:       mDoc,
								Signature: exprToString(method.Type),
							})
						}
					}
				}
				doc.Interfaces = append(doc.Interfaces, ifaceDoc)

			default:
				doc.Types = append(doc.Types, TypeDoc{
					Name:       s.Name.Name,
					Doc:        docText,
					Underlying: exprToString(s.Type),
				})
			}

		case *ast.ValueSpec:
			docText := ""
			if decl.Doc != nil {
				docText = decl.Doc.Text()
			}
			if s.Doc != nil {
				docText = s.Doc.Text()
			}

			for i, name := range s.Names {
				if !name.IsExported() {
					continue
				}
				typeName := ""
				if s.Type != nil {
					typeName = exprToString(s.Type)
				}
				valStr := ""
				if i < len(s.Values) {
					valStr = exprToString(s.Values[i])
				}

				switch decl.Tok {
				case token.CONST:
					doc.Constants = append(doc.Constants, ConstDoc{
						Name:  name.Name,
						Doc:   docText,
						Type:  typeName,
						Value: valStr,
					})
				case token.VAR:
					doc.Variables = append(doc.Variables, VarDoc{
						Name:  name.Name,
						Doc:   docText,
						Type:  typeName,
						Value: valStr,
					})
				}
			}
		}
	}
}

// parseFuncDecl parses function and method declarations
func parseFuncDecl(decl *ast.FuncDecl, doc *PackageDoc, fset *token.FileSet) {
	if !decl.Name.IsExported() {
		return
	}

	funcDoc := FunctionDoc{
		Name:       decl.Name.Name,
		IsExported: true,
	}

	if decl.Doc != nil {
		funcDoc.Doc = decl.Doc.Text()
	}

	// Build signature
	funcDoc.Signature = buildSignature(decl)

	// Receiver (method)
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		funcDoc.Receiver = exprToString(decl.Recv.List[0].Type)
	}

	// Parameters
	if decl.Type.Params != nil {
		for _, param := range decl.Type.Params.List {
			typeName := exprToString(param.Type)
			if len(param.Names) > 0 {
				for _, name := range param.Names {
					funcDoc.Params = append(funcDoc.Params, ParamDoc{
						Name: name.Name,
						Type: typeName,
					})
				}
			} else {
				funcDoc.Params = append(funcDoc.Params, ParamDoc{
					Type: typeName,
				})
			}
		}
	}

	// Return types
	if decl.Type.Results != nil {
		for _, result := range decl.Type.Results.List {
			funcDoc.Returns = append(funcDoc.Returns, exprToString(result.Type))
		}
	}

	// If it's a method, attach to struct
	if funcDoc.Receiver != "" {
		receiverName := cleanReceiverName(funcDoc.Receiver)
		attached := false
		for i, s := range doc.Structs {
			if s.Name == receiverName {
				doc.Structs[i].Methods = append(doc.Structs[i].Methods, funcDoc)
				attached = true
				break
			}
		}
		if !attached {
			// Add as a standalone function with receiver info
			doc.Functions = append(doc.Functions, funcDoc)
		}
	} else {
		doc.Functions = append(doc.Functions, funcDoc)
	}
}

// buildSignature creates a human-readable function signature
func buildSignature(decl *ast.FuncDecl) string {
	var sb strings.Builder

	sb.WriteString("func ")

	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		sb.WriteString("(")
		recv := decl.Recv.List[0]
		if len(recv.Names) > 0 {
			sb.WriteString(recv.Names[0].Name)
			sb.WriteString(" ")
		}
		sb.WriteString(exprToString(recv.Type))
		sb.WriteString(") ")
	}

	sb.WriteString(decl.Name.Name)
	sb.WriteString("(")

	if decl.Type.Params != nil {
		params := []string{}
		for _, param := range decl.Type.Params.List {
			typeName := exprToString(param.Type)
			if len(param.Names) > 0 {
				for _, name := range param.Names {
					params = append(params, name.Name+" "+typeName)
				}
			} else {
				params = append(params, typeName)
			}
		}
		sb.WriteString(strings.Join(params, ", "))
	}

	sb.WriteString(")")

	if decl.Type.Results != nil {
		results := []string{}
		for _, result := range decl.Type.Results.List {
			typeName := exprToString(result.Type)
			if len(result.Names) > 0 {
				for _, name := range result.Names {
					results = append(results, name.Name+" "+typeName)
				}
			} else {
				results = append(results, typeName)
			}
		}
		if len(results) == 1 {
			sb.WriteString(" ")
			sb.WriteString(results[0])
		} else if len(results) > 1 {
			sb.WriteString(" (")
			sb.WriteString(strings.Join(results, ", "))
			sb.WriteString(")")
		}
	}

	return sb.String()
}

// exprToString converts an AST expression to a string representation
func exprToString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.ArrayType:
		return "[]" + exprToString(e.Elt)
	case *ast.MapType:
		return "map[" + exprToString(e.Key) + "]" + exprToString(e.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)"
	case *ast.Ellipsis:
		return "..." + exprToString(e.Elt)
	case *ast.ChanType:
		return "chan " + exprToString(e.Value)
	case *ast.BasicLit:
		return e.Value
	case *ast.ParenExpr:
		return "(" + exprToString(e.X) + ")"
	case *ast.UnaryExpr:
		return e.Op.String() + exprToString(e.X)
	case *ast.BinaryExpr:
		return exprToString(e.X) + " " + e.Op.String() + " " + exprToString(e.Y)
	case *ast.CallExpr:
		return exprToString(e.Fun) + "(...)"
	case *ast.CompositeLit:
		return exprToString(e.Type) + "{...}"
	case *ast.IndexExpr:
		return exprToString(e.X) + "[" + exprToString(e.Index) + "]"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// cleanReceiverName extracts the type name from a receiver
func cleanReceiverName(receiver string) string {
	r := strings.TrimPrefix(receiver, "*")
	r = strings.TrimSpace(r)
	return r
}

// embeddedFieldName extracts a readable name from an embedded (anonymous) field type.
// e.g. *http.Handler → "Handler", sync.Mutex → "Mutex", io.Reader → "Reader"
func embeddedFieldName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return e.Sel.Name
	case *ast.StarExpr:
		return embeddedFieldName(e.X)
	case *ast.IndexExpr:
		return embeddedFieldName(e.X)
	default:
		s := exprToString(expr)
		// Last resort: take the last part after any "."
		if idx := strings.LastIndex(s, "."); idx >= 0 {
			return s[idx+1:]
		}
		return s
	}
}
