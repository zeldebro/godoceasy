package docs

import (
	"fmt"
	"strings"

	"github.com/godoceasy/internal/parser"
)

// FriendlyDoc represents beginner-friendly documentation for a package
type FriendlyDoc struct {
	PackageName string
	ImportPath  string
	Overview    string
	WhatItDoes  string
	WhenToUse   string
	Structs     []FriendlyStruct
	Interfaces  []FriendlyInterface
	Functions   []FriendlyFunction
	Types       []FriendlyType
	Constants   []FriendlyConst
	Variables   []FriendlyVar
	Files       []string
}

// FriendlyStruct is a beginner-friendly struct doc
type FriendlyStruct struct {
	Name        string
	Overview    string
	WhatItDoes  string
	WhenToUse   string
	Fields      []FriendlyField
	Methods     []FriendlyFunction
	Example     string
	Explanation string
}

// FriendlyField is a beginner-friendly field doc
type FriendlyField struct {
	Name        string
	Type        string
	Tag         string
	Description string
}

// FriendlyInterface is a beginner-friendly interface doc
type FriendlyInterface struct {
	Name        string
	Overview    string
	WhatItDoes  string
	WhenToUse   string
	Methods     []FriendlyFunction
	Example     string
	Explanation string
}

// FriendlyFunction is a beginner-friendly function doc
type FriendlyFunction struct {
	Name        string
	Signature   string
	Overview    string
	WhatItDoes  string
	WhenToUse   string
	Params      []FriendlyParam
	Returns     []string
	Receiver    string
	Example     string
	Explanation string
}

// FriendlyParam is a beginner-friendly parameter doc
type FriendlyParam struct {
	Name        string
	Type        string
	Description string
}

// FriendlyType is a beginner-friendly type doc
type FriendlyType struct {
	Name          string
	Underlying    string
	Overview      string
	WhatItDoes    string
	WhenToUse     string
	InSimpleWords string
	Example       string
	Explanation   string
}

// FriendlyConst is a beginner-friendly constant doc
type FriendlyConst struct {
	Name        string
	Value       string
	Type        string
	Description string
}

// FriendlyVar is a beginner-friendly variable doc
type FriendlyVar struct {
	Name        string
	Value       string
	Type        string
	Description string
}

// Generate creates beginner-friendly documentation from parsed package docs
func Generate(pkgDocs []*parser.PackageDoc) []*FriendlyDoc {
	var result []*FriendlyDoc

	for _, pkg := range pkgDocs {
		fd := &FriendlyDoc{
			PackageName: pkg.Name,
			ImportPath:  pkg.ImportPath,
			Overview:    generateOverview(pkg),
			WhatItDoes:  generateWhatItDoes(pkg),
			WhenToUse:   generateWhenToUse(pkg),
			Files:       pkg.Files,
		}

		// Process structs
		for _, s := range pkg.Structs {
			fd.Structs = append(fd.Structs, generateFriendlyStruct(s, pkg.Name, pkg.ImportPath))
		}

		// Process interfaces
		for _, i := range pkg.Interfaces {
			fd.Interfaces = append(fd.Interfaces, generateFriendlyInterface(i, pkg.Name))
		}

		// Process functions
		for _, f := range pkg.Functions {
			fd.Functions = append(fd.Functions, generateFriendlyFunction(f, pkg.Name, pkg.ImportPath))
		}

		// Process types
		for _, t := range pkg.Types {
			fd.Types = append(fd.Types, generateFriendlyType(t, pkg.Name))
		}

		// Process constants
		for _, c := range pkg.Constants {
			fd.Constants = append(fd.Constants, FriendlyConst{
				Name:        c.Name,
				Value:       c.Value,
				Type:        c.Type,
				Description: cleanDoc(c.Doc),
			})
		}

		// Process variables
		for _, v := range pkg.Variables {
			fd.Variables = append(fd.Variables, FriendlyVar{
				Name:        v.Name,
				Value:       v.Value,
				Type:        v.Type,
				Description: cleanDoc(v.Doc),
			})
		}

		result = append(result, fd)
	}

	return result
}

func generateOverview(pkg *parser.PackageDoc) string {
	// PRIMARY SOURCE: Use actual doc comment from the source code
	if pkg.Doc != "" {
		return cleanDoc(pkg.Doc)
	}

	// LAST RESORT: no doc comment — summarize from contents
	var points []string
	points = append(points, fmt.Sprintf("Package <strong>%s</strong> — no doc comment found in source code.", pkg.Name))
	if len(pkg.Structs) > 0 {
		points = append(points, fmt.Sprintf("Provides <strong>%d</strong> struct type(s).", len(pkg.Structs)))
	}
	if len(pkg.Interfaces) > 0 {
		points = append(points, fmt.Sprintf("Defines <strong>%d</strong> interface(s).", len(pkg.Interfaces)))
	}
	if len(pkg.Functions) > 0 {
		points = append(points, fmt.Sprintf("Includes <strong>%d</strong> exported function(s).", len(pkg.Functions)))
	}
	if len(pkg.Constants) > 0 {
		points = append(points, fmt.Sprintf("Contains <strong>%d</strong> constant(s).", len(pkg.Constants)))
	}

	return buildPointList(points)
}

func generateWhatItDoes(pkg *parser.PackageDoc) string {
	var sb strings.Builder

	// PRIMARY SOURCE: Always use actual doc comment from the source code
	if pkg.Doc != "" {
		sb.WriteString(cleanDoc(pkg.Doc))
	}

	// SUPPLEMENTAL: If knowledge base has practical tasks, add them AFTER the real doc
	if k := lookupPackageKnowledge(pkg.ImportPath); k != nil && len(k.Tasks) > 0 {
		sb.WriteString("<p>🎯 <strong>What you can do with this package:</strong></p>")
		sb.WriteString(generateTaskList(k.Tasks))
		return sb.String()
	}

	// If doc comment was found, return it (with struct/func summary)
	if sb.Len() > 0 {
		// Add a brief summary of what's inside
		var extras []string
		if len(pkg.Structs) > 0 {
			names := collectNames(pkg.Structs, 4, func(s parser.StructDoc) string { return s.Name })
			extras = append(extras, fmt.Sprintf("Defines structs: %s", names))
		}
		if len(pkg.Interfaces) > 0 {
			names := collectNames(pkg.Interfaces, 4, func(i parser.InterfaceDoc) string { return i.Name })
			extras = append(extras, fmt.Sprintf("Defines interfaces: %s", names))
		}
		if len(pkg.Functions) > 0 {
			names := collectNames(pkg.Functions, 4, func(f parser.FunctionDoc) string { return f.Name })
			extras = append(extras, fmt.Sprintf("Provides functions: %s", names))
		}
		if len(extras) > 0 {
			sb.WriteString(buildPointList(extras))
		}
		return sb.String()
	}

	// LAST RESORT: no doc comment, no knowledge base — summarize from contents
	var points []string
	if len(pkg.Structs) > 0 {
		names := collectNames(pkg.Structs, 5, func(s parser.StructDoc) string { return s.Name })
		points = append(points, fmt.Sprintf("Defines structs: %s", names))
	}
	if len(pkg.Interfaces) > 0 {
		names := collectNames(pkg.Interfaces, 5, func(i parser.InterfaceDoc) string { return i.Name })
		points = append(points, fmt.Sprintf("Defines interfaces: %s", names))
	}
	if len(pkg.Functions) > 0 {
		names := collectNames(pkg.Functions, 5, func(f parser.FunctionDoc) string { return f.Name })
		points = append(points, fmt.Sprintf("Provides functions: %s", names))
	}
	if len(pkg.Constants) > 0 {
		points = append(points, fmt.Sprintf("Exports %d constant(s).", len(pkg.Constants)))
	}

	if len(points) == 0 {
		return fmt.Sprintf("Package %s is a Go package.", pkg.Name)
	}

	return buildPointList(points)
}

func generateWhenToUse(pkg *parser.PackageDoc) string {
	var points []string

	// PRIMARY SOURCE: Use actual doc comment to derive "when to use"
	if pkg.Doc != "" {
		sentences := splitSentences(pkg.Doc)
		if len(sentences) > 0 {
			first := strings.TrimSpace(sentences[0])
			first = strings.TrimSuffix(first, ".")
			points = append(points, fmt.Sprintf("Use this package when you need to <strong>%s</strong>.", strings.ToLower(first)))
		}
	}

	points = append(points, fmt.Sprintf("📝 <strong>Import it:</strong> <code>import \"%s\"</code>", pkg.ImportPath))

	// SUPPLEMENTAL: Knowledge base adds practical task tips
	if k := lookupPackageKnowledge(pkg.ImportPath); k != nil && len(k.Tasks) > 0 {
		points = append(points, "🎯 <strong>Common tasks:</strong>")
		for i, t := range k.Tasks {
			if i >= 4 {
				break
			}
			points = append(points, fmt.Sprintf("&nbsp;&nbsp;• <strong>%s</strong> → %s", t.Task, t.How))
		}
	}

	if len(pkg.Structs) > 0 || len(pkg.Interfaces) > 0 {
		structNames := []string{}
		for _, s := range pkg.Structs {
			if len(structNames) < 3 {
				structNames = append(structNames, "<code>"+s.Name+"</code>")
			}
		}
		if len(structNames) > 0 {
			points = append(points, fmt.Sprintf("📦 <strong>Main types:</strong> %s — create them, set fields, call methods.", strings.Join(structNames, ", ")))
		}
	}
	if len(pkg.Functions) > 0 {
		funcNames := []string{}
		for _, f := range pkg.Functions {
			if f.Receiver == "" && len(funcNames) < 3 {
				funcNames = append(funcNames, "<code>"+f.Name+"()</code>")
			}
		}
		if len(funcNames) > 0 {
			points = append(points, fmt.Sprintf("⚡ <strong>Key functions:</strong> %s — call them directly, no setup needed.", strings.Join(funcNames, ", ")))
		}
	}

	points = append(points, fmt.Sprintf("💡 <strong>Quick start:</strong> Pick a function or type from the sidebar and follow the example code."))

	if len(points) == 1 {
		return points[0]
	}

	var sb strings.Builder
	sb.WriteString("<ul class=\"doc-points\">")
	for _, p := range points {
		sb.WriteString("<li>")
		sb.WriteString(p)
		sb.WriteString("</li>")
	}
	sb.WriteString("</ul>")
	return sb.String()
}

func generateFriendlyStruct(s parser.StructDoc, pkgName string, importPath string) FriendlyStruct {
	// Check knowledge base for struct-specific knowledge
	sk := lookupStructKnowledge(importPath, s.Name)

	// Count exported vs total fields
	exportedFields := 0
	for _, f := range s.Fields {
		if f.Name != "" && len(f.Name) > 0 && f.Name[0] >= 'A' && f.Name[0] <= 'Z' {
			exportedFields++
		}
	}

	// Build rich WhatItDoes
	var whatPoints []string
	if s.Doc != "" {
		whatPoints = append(whatPoints, cleanDoc(s.Doc))
	} else {
		whatPoints = append(whatPoints, fmt.Sprintf("<code>%s</code> is a struct (data container) in the <code>%s</code> package.", s.Name, pkgName))
	}

	// Describe the shape of the struct
	if exportedFields > 0 {
		fieldNames := []string{}
		for _, f := range s.Fields {
			if f.Name != "" && len(f.Name) > 0 && f.Name[0] >= 'A' && f.Name[0] <= 'Z' {
				fieldNames = append(fieldNames, "<code>"+f.Name+"</code>")
				if len(fieldNames) >= 5 {
					break
				}
			}
		}
		if exportedFields <= 5 {
			whatPoints = append(whatPoints, fmt.Sprintf("Fields: %s", strings.Join(fieldNames, ", ")))
		} else {
			whatPoints = append(whatPoints, fmt.Sprintf("Has %d fields including %s and %d more.", exportedFields, strings.Join(fieldNames, ", "), exportedFields-5))
		}
	}
	if len(s.Methods) > 0 {
		methodNames := []string{}
		builderCount := 0
		for _, m := range s.Methods {
			if isBuilderMethod(m) {
				builderCount++
			}
			methodNames = append(methodNames, "<code>."+m.Name+"()</code>")
			if len(methodNames) >= 4 {
				break
			}
		}
		whatPoints = append(whatPoints, fmt.Sprintf("Methods: %s", strings.Join(methodNames, ", ")))
		if builderCount > 0 {
			whatPoints = append(whatPoints, fmt.Sprintf("🔗 <strong>%d builder method(s)</strong> — you can chain calls like <code>obj.%s(...).%s(...)</code>", builderCount, s.Methods[0].Name, s.Methods[min(1, len(s.Methods)-1)].Name))
		}
	}

	whatItDoes := buildPointList(whatPoints)

	// SUPPLEMENTAL: Knowledge base adds practical tips AFTER the real doc
	if sk != nil && sk.WhatItDoes != "" {
		whatItDoes = whatItDoes + "<p>💡 <em>" + sk.WhatItDoes + "</em></p>"
	}

	// Build smarter WhenToUse — doc comment first, knowledge base supplements
	var whenPoints []string
	if s.Doc != "" {
		// Derive "when to use" from the struct's own doc comment
		first := strings.TrimSpace(strings.Split(s.Doc, "\n")[0])
		first = strings.TrimSuffix(first, ".")
		if first != "" {
			whenPoints = append(whenPoints, fmt.Sprintf("Use <code>%s.%s</code> when you need to <strong>%s</strong>.", pkgName, s.Name, strings.ToLower(first)))
		}
	}
	// Add knowledge base tips as supplement
	if sk != nil && sk.WhenToUse != "" {
		whenPoints = append(whenPoints, "💡 "+sk.WhenToUse)
		if sk.HowToUse != "" {
			whenPoints = append(whenPoints, fmt.Sprintf("📝 <strong>How to use in your code:</strong><pre class=\"task-code\"><code>%s</code></pre>", sk.HowToUse))
		}
	}
	// Fallback if nothing from doc or knowledge base
	if len(whenPoints) == 0 {
		whenPoints = generateStructWhenToUse(s, pkgName, exportedFields)
	}

	// In Simple Words — real-world analogy even non-tech people understand
	simpleWords := generateStructAnalogy(s.Name, exportedFields, len(s.Methods))

	fs := FriendlyStruct{
		Name:       s.Name,
		Overview:   cleanDoc(s.Doc),
		WhatItDoes: whatItDoes,
		WhenToUse:  buildPointList(whenPoints),
	}

	// Fields with smarter descriptions — never leave Name or Description empty
	for _, f := range s.Fields {
		ff := FriendlyField{
			Name: f.Name,
			Type: f.Type,
			Tag:  f.Tag,
		}
		// Ensure Name is never empty
		if ff.Name == "" {
			ff.Name = f.Type // Use type as name for embedded fields
		}
		// Generate description
		if f.Doc != "" {
			ff.Description = cleanDoc(f.Doc)
		} else {
			ff.Description = describeField(ff.Name, f.Type)
		}
		// Fallback: if description is still empty, create a generic one
		if ff.Description == "" {
			ff.Description = fmt.Sprintf("<code>%s</code> — field of type <code>%s</code>.", ff.Name, f.Type)
		}
		fs.Fields = append(fs.Fields, ff)
	}

	// Methods
	for _, m := range s.Methods {
		fs.Methods = append(fs.Methods, generateFriendlyFunction(m, pkgName))
	}

	// Generate example
	fs.Example = generateStructExample(s, pkgName)
	fs.Explanation = simpleWords

	return fs
}

// generateStructWhenToUse creates descriptive, practical "When To Use" guidance for a struct
func generateStructWhenToUse(s parser.StructDoc, pkgName string, exportedFields int) []string {
	n := strings.ToLower(s.Name)
	var points []string

	// Context-specific practical guidance
	switch {
	case strings.Contains(n, "server"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> when you want to <strong>start a server</strong> with custom settings (port, timeouts, TLS, etc.).", pkgName, s.Name))
		points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Create it → set the fields you need → call a method like <code>.ListenAndServe()</code>"))
		points = append(points, fmt.Sprintf("💡 You can start with just one field: <code>%s.%s{Addr: \":8080\"}</code> — Go uses sensible defaults for the rest.", pkgName, s.Name))
	case strings.Contains(n, "client"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> when you need to <strong>make requests</strong> to other services (HTTP calls, API requests, etc.).", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> Create it once → reuse it for many requests (it's safe to share across goroutines).")
		points = append(points, fmt.Sprintf("💡 An empty <code>%s.%s{}</code> works fine — customize only when you need timeouts, proxies, or TLS settings.", pkgName, s.Name))
	case strings.Contains(n, "request") || strings.Contains(n, "req"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>describe what you want</strong> — the URL, method (GET/POST), headers, and body.", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> Usually created by a helper function, then passed to a client or handler.")
		points = append(points, "💡 You rarely create this from scratch — look for <code>NewRequest()</code> or similar factory functions.")
	case strings.Contains(n, "response") || strings.Contains(n, "resp"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>read the result</strong> of a request — status code, headers, and body.", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> You get this back from a function call — read its fields to see what happened.")
		points = append(points, "⚠️ <strong>Important:</strong> Always close the body when done: <code>defer resp.Body.Close()</code>")
	case strings.Contains(n, "config") || strings.Contains(n, "option") || strings.Contains(n, "setting"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>configure behavior</strong> before passing it to another function or struct.", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> Create it → set the options you care about → pass it as an argument.")
		points = append(points, "💡 Think of it as a settings file — you only fill in what you want to change from the defaults.")
	case strings.Contains(n, "handler"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>handle incoming requests or events</strong>.", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> Create it → register it with a router or server → it gets called automatically when events arrive.")
	case strings.Contains(n, "error") || strings.Contains(n, "err"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>represent and inspect errors</strong> with more detail than a plain string.", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> Check if an error is this type using <code>errors.As()</code>, then read its fields for details.")
	case strings.Contains(n, "pool"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>reuse expensive resources</strong> (connections, buffers) instead of creating new ones.", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> Create the pool once at startup → use <code>.Get()</code> and <code>.Put()</code> to borrow and return items.")
	case strings.Contains(n, "writer"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>write data</strong> to a destination (file, network, buffer).", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> Create it with a destination → call <code>.Write()</code> to send data.")
	case strings.Contains(n, "reader"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>read data</strong> from a source (file, network, string).", pkgName, s.Name))
		points = append(points, "📝 <strong>How:</strong> Create it with a source → call <code>.Read()</code> to get data piece by piece.")
	default:
		// Generic but still descriptive
		if len(s.Methods) > 0 && exportedFields > 0 {
			points = append(points, fmt.Sprintf("Use <code>%s.%s</code> when you need to <strong>store data and perform actions</strong> on it.", pkgName, s.Name))
			points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Create it with <code>%s.%s{...}</code> → fill in the fields → call its methods.", pkgName, s.Name))
		} else if len(s.Methods) > 0 {
			points = append(points, fmt.Sprintf("Use <code>%s.%s</code> when you need the <strong>actions (methods)</strong> it provides.", pkgName, s.Name))
			points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Create it → call <code>.MethodName()</code> to do work."))
		} else {
			points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to <strong>hold related data together</strong> in one place.", pkgName, s.Name))
			points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Create it with <code>%s.%s{Field: value}</code> and pass it around.", pkgName, s.Name))
		}
	}

	if exportedFields > 3 {
		points = append(points, "💡 <strong>Tip:</strong> You don't need to set every field — only the ones you need. Go uses zero-values for the rest.")
	}

	return points
}

// generateFriendlyType creates a rich, beginner-friendly type documentation
func generateFriendlyType(t parser.TypeDoc, pkgName string) FriendlyType {
	ft := FriendlyType{
		Name:       t.Name,
		Underlying: t.Underlying,
		Overview:   cleanDoc(t.Doc),
	}

	// Build rich WhatItDoes
	var whatPoints []string
	if t.Doc != "" {
		whatPoints = append(whatPoints, cleanDoc(t.Doc))
	}
	whatPoints = append(whatPoints, fmt.Sprintf("<code>type %s %s</code> — a named type based on <code>%s</code>.", t.Name, t.Underlying, t.Underlying))

	// Explain the underlying type in simple terms
	underlyingExplain := explainType(t.Underlying)
	if underlyingExplain != "" {
		whatPoints = append(whatPoints, underlyingExplain)
	}

	whatPoints = append(whatPoints, fmt.Sprintf("Having a named type lets Go distinguish <code>%s</code> from a plain <code>%s</code>, which prevents mixing them up.", t.Name, t.Underlying))
	ft.WhatItDoes = buildPointList(whatPoints)

	// Build descriptive WhenToUse with practical guidance
	typeWhenPoints := generateTypeWhenToUse(t, pkgName)
	ft.WhenToUse = buildPointList(typeWhenPoints)

	// InSimpleWords — real-world analogy
	ft.InSimpleWords = generateTypeAnalogy(t)

	// Example
	ft.Example = generateTypeExample(t, pkgName)
	ft.Explanation = fmt.Sprintf("%s is just a %s with a meaningful name. Use it wherever the package expects it.", t.Name, t.Underlying)

	return ft
}

// generateTypeAnalogy creates a real-world analogy for a type
func generateTypeAnalogy(t parser.TypeDoc) string {
	n := strings.ToLower(t.Name)

	switch {
	case strings.Contains(n, "handler") && strings.HasPrefix(t.Underlying, "func"):
		return fmt.Sprintf("👋 <code>%s</code> is like a specific job title for a function. Any function with the right signature can wear this name tag and be used wherever a <code>%s</code> is needed.", t.Name, t.Name)
	case strings.Contains(n, "status"):
		return fmt.Sprintf("🚦 <code>%s</code> is like a traffic light color — it's just a %s underneath, but the name tells you it represents a status (healthy, error, pending, etc.).", t.Name, t.Underlying)
	case strings.Contains(n, "phase"):
		return fmt.Sprintf("🔄 <code>%s</code> is like a stage in a process — \"preparing\", \"running\", \"finished\". It's a %s with a name that says \"I represent a phase.\"", t.Name, t.Underlying)
	case strings.Contains(n, "method") && t.Underlying == "string":
		return fmt.Sprintf("📬 <code>%s</code> is like labeling a letter as \"Express\" or \"Standard\" — it's just a string, but the label tells you it represents a method (GET, POST, PUT...).", t.Name)
	case strings.Contains(n, "protocol") || strings.Contains(n, "scheme"):
		return fmt.Sprintf("🔌 <code>%s</code> is like choosing between \"phone\" and \"email\" for communication. The underlying %s holds the protocol name, but the type makes sure you don't accidentally mix it up with a random string.", t.Name, t.Underlying)
	}

	// Generic analogies based on underlying type
	switch t.Underlying {
	case "string":
		return fmt.Sprintf("🏷️ <code>%s</code> is like a labeled sticky note. Underneath it's just text (a string), but the label <code>%s</code> tells Go — and you — what the text means. You can't accidentally mix it up with a random string.", t.Name, t.Name)
	case "int", "int32", "int64":
		return fmt.Sprintf("🔢 <code>%s</code> is like a numbered jersey. It's just a number underneath, but the jersey tells you what the number represents (a player, not a score).", t.Name)
	case "bool":
		return fmt.Sprintf("🔘 <code>%s</code> is like a labeled light switch. It's just on/off (true/false), but the label tells you what it controls.", t.Name)
	case "float32", "float64":
		return fmt.Sprintf("📏 <code>%s</code> is like a measurement with units. It's a decimal number underneath, but the name tells you what's being measured.", t.Name)
	default:
		if strings.HasPrefix(t.Underlying, "func(") {
			return fmt.Sprintf("🎯 <code>%s</code> is like a job title for a function. Any function with the right skills (matching signature) can do this job. The name just makes it official.", t.Name)
		}
		if strings.HasPrefix(t.Underlying, "[]") {
			return fmt.Sprintf("📚 <code>%s</code> is like a named shelf. It holds a list of items, and the name tells you what kind of items belong on this shelf.", t.Name)
		}
		if strings.HasPrefix(t.Underlying, "map[") {
			return fmt.Sprintf("🗂️ <code>%s</code> is like a labeled filing cabinet. It stores key-value pairs, and the label tells you what kind of files go inside.", t.Name)
		}
		return fmt.Sprintf("🏷️ <code>%s</code> is like a label on top of <code>%s</code>. It works exactly the same way underneath, but the label tells Go (and you) what the value means — so you don't mix things up.", t.Name, t.Underlying)
	}
}

// explainType returns a simple explanation of what a Go type is
func explainType(typeName string) string {
	switch typeName {
	case "string":
		return "The underlying type is <code>string</code> — a text value."
	case "int", "int32", "int64":
		return fmt.Sprintf("The underlying type is <code>%s</code> — a whole number.", typeName)
	case "float32", "float64":
		return fmt.Sprintf("The underlying type is <code>%s</code> — a decimal number.", typeName)
	case "bool":
		return "The underlying type is <code>bool</code> — true or false."
	case "byte":
		return "The underlying type is <code>byte</code> — a single byte (alias for uint8)."
	case "rune":
		return "The underlying type is <code>rune</code> — a single Unicode character."
	}
	if strings.HasPrefix(typeName, "[]") {
		inner := strings.TrimPrefix(typeName, "[]")
		return fmt.Sprintf("The underlying type is a <strong>slice</strong> (list) of <code>%s</code> values.", inner)
	}
	if strings.HasPrefix(typeName, "map[") {
		return fmt.Sprintf("The underlying type is a <strong>map</strong> (key-value dictionary): <code>%s</code>.", typeName)
	}
	if strings.HasPrefix(typeName, "func(") {
		return "The underlying type is a <strong>function type</strong> — a variable that holds a function."
	}
	if strings.HasPrefix(typeName, "chan") {
		return "The underlying type is a <strong>channel</strong> — used for communication between goroutines."
	}
	return ""
}

// generateTypeExample creates a Hello-World simple example for a type
func generateTypeExample(t parser.TypeDoc, pkgName string) string {
	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString(fmt.Sprintf("import \"%s\"\n", pkgName))
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")

	switch t.Underlying {
	case "string":
		sb.WriteString(fmt.Sprintf("    // Create a %s (it's just a string with a special name)\n", t.Name))
		sb.WriteString(fmt.Sprintf("    myValue := %s.%s(\"hello-world\")\n\n", pkgName, t.Name))
		sb.WriteString("    // Print it — works just like a string!\n")
		sb.WriteString(fmt.Sprintf("    fmt.Println(\"My %s is:\", myValue)\n", t.Name))
		sb.WriteString(fmt.Sprintf("    // Output: My %s is: hello-world\n", t.Name))
	case "int", "int32", "int64":
		sb.WriteString(fmt.Sprintf("    // Create a %s (it's just a number with a special name)\n", t.Name))
		sb.WriteString(fmt.Sprintf("    myValue := %s.%s(42)\n\n", pkgName, t.Name))
		sb.WriteString("    // Print it — works just like a number!\n")
		sb.WriteString(fmt.Sprintf("    fmt.Println(\"My %s is:\", myValue)\n", t.Name))
		sb.WriteString(fmt.Sprintf("    // Output: My %s is: 42\n", t.Name))
	case "bool":
		sb.WriteString(fmt.Sprintf("    // Create a %s (it's just true/false with a special name)\n", t.Name))
		sb.WriteString(fmt.Sprintf("    myValue := %s.%s(true)\n\n", pkgName, t.Name))
		sb.WriteString("    // Print it\n")
		sb.WriteString(fmt.Sprintf("    fmt.Println(\"%s:\", myValue)\n", t.Name))
		sb.WriteString(fmt.Sprintf("    // Output: %s: true\n", t.Name))
	default:
		if strings.HasPrefix(t.Underlying, "func(") {
			sb.WriteString(fmt.Sprintf("    // Create a %s (it's a function with a special name)\n", t.Name))
			sb.WriteString(fmt.Sprintf("    var myFunc %s.%s\n\n", pkgName, t.Name))
			sb.WriteString("    // Assign it a function that does something\n")
			sb.WriteString("    // myFunc = func(...) { fmt.Println(\"Hello!\") }\n\n")
			sb.WriteString("    fmt.Println(\"Function ready:\", myFunc != nil)\n")
		} else {
			sb.WriteString(fmt.Sprintf("    // Create a %s\n", t.Name))
			sb.WriteString(fmt.Sprintf("    var myValue %s.%s\n\n", pkgName, t.Name))
			sb.WriteString("    // Use it wherever the package expects this type\n")
			sb.WriteString(fmt.Sprintf("    fmt.Println(\"%s:\", myValue)\n", t.Name))
		}
	}

	sb.WriteString("}")
	return sb.String()
}

// describeField generates a smart description for a struct field based on name and type
func describeField(name string, typeName string) string {
	n := strings.ToLower(name)

	// Context-aware descriptions
	switch {
	case strings.Contains(n, "name"):
		return fmt.Sprintf("<code>%s</code> — the name identifier (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "namespace"):
		return fmt.Sprintf("<code>%s</code> — the Kubernetes namespace (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "label"):
		return fmt.Sprintf("<code>%s</code> — labels for selection/filtering (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "annotation"):
		return fmt.Sprintf("<code>%s</code> — metadata annotations (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "spec"):
		return fmt.Sprintf("<code>%s</code> — the desired state/configuration (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "status"):
		return fmt.Sprintf("<code>%s</code> — the current observed state (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "error") || strings.Contains(n, "err"):
		return fmt.Sprintf("<code>%s</code> — error information if something went wrong (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "timeout"):
		return fmt.Sprintf("<code>%s</code> — how long to wait before timing out (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "port"):
		return fmt.Sprintf("<code>%s</code> — network port number (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "host"):
		return fmt.Sprintf("<code>%s</code> — hostname or IP address (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "path") || strings.Contains(n, "url") || strings.Contains(n, "uri"):
		return fmt.Sprintf("<code>%s</code> — file path or URL (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "image"):
		return fmt.Sprintf("<code>%s</code> — container image reference (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "replica"):
		return fmt.Sprintf("<code>%s</code> — number of replicas/copies (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "secret"):
		return fmt.Sprintf("<code>%s</code> — reference to a secret/credential (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "config"):
		return fmt.Sprintf("<code>%s</code> — configuration settings (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "create") || strings.Contains(n, "created"):
		return fmt.Sprintf("<code>%s</code> — when this was created (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "update") || strings.Contains(n, "updated"):
		return fmt.Sprintf("<code>%s</code> — when this was last updated (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "delete") || strings.Contains(n, "deleted"):
		return fmt.Sprintf("<code>%s</code> — when this was deleted (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "enabled") || strings.Contains(n, "disabled"):
		return fmt.Sprintf("<code>%s</code> — whether this feature is on or off (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "count") || strings.Contains(n, "size") || strings.Contains(n, "length"):
		return fmt.Sprintf("<code>%s</code> — a numeric count or size (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "type") || strings.Contains(n, "kind"):
		return fmt.Sprintf("<code>%s</code> — the type or kind of this object (<code>%s</code>).", name, typeName)
	case strings.Contains(n, "version"):
		return fmt.Sprintf("<code>%s</code> — version identifier (<code>%s</code>).", name, typeName)
	}

	// Generic but still useful
	typeExplain := ""
	switch {
	case typeName == "string":
		typeExplain = "text value"
	case typeName == "bool":
		typeExplain = "true/false flag"
	case typeName == "int" || typeName == "int32" || typeName == "int64":
		typeExplain = "number"
	case typeName == "error":
		typeExplain = "error value"
	case strings.HasPrefix(typeName, "[]"):
		typeExplain = "list of " + strings.TrimPrefix(typeName, "[]")
	case strings.HasPrefix(typeName, "map["):
		typeExplain = "key-value map"
	case strings.HasPrefix(typeName, "*"):
		typeExplain = "pointer to " + strings.TrimPrefix(typeName, "*")
	default:
		typeExplain = typeName
	}

	return fmt.Sprintf("<code>%s</code> — %s (<code>%s</code>).", name, typeExplain, typeName)
}

// generateStructAnalogy creates a real-world analogy for a struct
func generateStructAnalogy(name string, fieldCount int, methodCount int) string {
	n := strings.ToLower(name)

	// Context-specific analogies
	switch {
	case strings.Contains(n, "server"):
		return fmt.Sprintf("🏪 Think of <code>%s</code> like a restaurant. The fields are its settings (address, menu, hours), and the methods are what it does (open, serve customers, close). You configure it first, then tell it to start.", name)
	case strings.Contains(n, "client"):
		return fmt.Sprintf("📱 Think of <code>%s</code> like a phone. You set it up once (phone number, settings), then use it to make calls (send requests). The fields are the phone settings, the methods are the actions.", name)
	case strings.Contains(n, "request") || strings.Contains(n, "req"):
		return fmt.Sprintf("✉️ Think of <code>%s</code> like a letter. It has an address (where to send it), a message (the body), and extra info (headers). You fill in the details, then hand it to the mail service.", name)
	case strings.Contains(n, "response") || strings.Contains(n, "resp"):
		return fmt.Sprintf("📬 Think of <code>%s</code> like a reply letter you receive. It contains a status (delivered? rejected?), headers (metadata), and the actual message body.", name)
	case strings.Contains(n, "config") || strings.Contains(n, "option") || strings.Contains(n, "setting"):
		return fmt.Sprintf("🎛️ Think of <code>%s</code> like the settings panel on your phone. Each field is one setting you can tweak. You don't have to change all of them — just the ones that matter to you.", name)
	case strings.Contains(n, "handler"):
		return fmt.Sprintf("👋 Think of <code>%s</code> like a receptionist. When someone arrives (a request comes in), the handler decides what to do with it — answer a question, redirect them, or send them away.", name)
	case strings.Contains(n, "error") || strings.Contains(n, "err"):
		return fmt.Sprintf("🚨 Think of <code>%s</code> like an error report. It captures what went wrong, when, and where — so you can figure out how to fix it.", name)
	case strings.Contains(n, "pool"):
		return fmt.Sprintf("🏊 Think of <code>%s</code> like a swimming pool of reusable resources. Instead of creating something new every time, you borrow one from the pool and return it when done.", name)
	case strings.Contains(n, "cache"):
		return fmt.Sprintf("🗃️ Think of <code>%s</code> like a shelf near your desk. Instead of walking to the warehouse (slow database) every time, you keep frequently-used items on the shelf (cache) for quick access.", name)
	case strings.Contains(n, "writer"):
		return fmt.Sprintf("✍️ Think of <code>%s</code> like a pen. You create it, point it at a notebook (destination), and then write stuff to it. The fields configure where and how it writes.", name)
	case strings.Contains(n, "reader"):
		return fmt.Sprintf("📖 Think of <code>%s</code> like reading glasses. You point it at text (a data source), and it helps you read the content piece by piece.", name)
	case strings.Contains(n, "transport"):
		return fmt.Sprintf("🚚 Think of <code>%s</code> like a delivery truck. It handles the low-level work of moving data from point A to point B. The fields configure how fast, how many trips, timeouts, etc.", name)
	case strings.Contains(n, "token"):
		return fmt.Sprintf("🎫 Think of <code>%s</code> like a movie ticket. It proves you have permission to access something, and it might expire after a while.", name)
	}

	// Generic but still friendly analogies based on shape
	if methodCount == 0 {
		if fieldCount <= 3 {
			return fmt.Sprintf("📋 Think of <code>%s</code> as a sticky note with %d piece(s) of info written on it. Simple and to-the-point — just the data you need.", name, fieldCount)
		}
		return fmt.Sprintf("📝 Think of <code>%s</code> as a paper form with %d fields to fill in. You write down the info, and pass the form along to something that needs it.", name, fieldCount)
	}
	if fieldCount == 0 && methodCount > 0 {
		return fmt.Sprintf("🔧 Think of <code>%s</code> as a tool with %d button(s). There's no data to fill in — you just pick it up and press the buttons (call the methods).", name, methodCount)
	}
	return fmt.Sprintf("📦 Think of <code>%s</code> as a toolbox. It has %d compartment(s) (fields) to store things, and %d tool(s) (methods) to do work with. Set it up, then use it.", name, fieldCount, methodCount)
}

func generateFriendlyInterface(i parser.InterfaceDoc, pkgName string) FriendlyInterface {
	// Build point-wise WhatItDoes
	var whatPoints []string
	whatPoints = append(whatPoints, fmt.Sprintf("<code>%s</code> is an interface in the <code>%s</code> package.", i.Name, pkgName))
	whatPoints = append(whatPoints, "It defines a contract that types must implement.")
	if len(i.Methods) > 0 {
		names := []string{}
		for _, m := range i.Methods {
			names = append(names, "<code>"+m.Name+"</code>")
		}
		whatPoints = append(whatPoints, fmt.Sprintf("Required methods: %s", strings.Join(names, ", ")))
	}

	whatItDoes := buildPointList(whatPoints)
	if i.Doc != "" {
		whatItDoes = cleanDoc(i.Doc)
	}

	// Build descriptive WhenToUse with practical guidance
	ifaceWhenPoints := generateInterfaceWhenToUse(i, pkgName)

	fi := FriendlyInterface{
		Name:       i.Name,
		Overview:   cleanDoc(i.Doc),
		WhatItDoes: whatItDoes,
		WhenToUse:  buildPointList(ifaceWhenPoints),
	}

	for _, m := range i.Methods {
		fi.Methods = append(fi.Methods, generateFriendlyFunction(m, pkgName))
	}

	fi.Example = generateInterfaceExample(i, pkgName)

	// Build a clear "In Simple Words" explanation with real-world analogy
	fi.Explanation = generateInterfaceAnalogy(i)

	return fi
}

// generateInterfaceAnalogy creates a real-world analogy for an interface
func generateInterfaceAnalogy(i parser.InterfaceDoc) string {
	n := strings.ToLower(i.Name)

	// Context-specific analogies
	switch {
	case strings.Contains(n, "reader"):
		return fmt.Sprintf("📖 Think of <code>%s</code> like the ability to read. Books, newspapers, and phones are all different things, but they all satisfy the \"readable\" requirement. Any type with a <code>Read()</code> method counts.", i.Name)
	case strings.Contains(n, "writer"):
		return fmt.Sprintf("✍️ Think of <code>%s</code> like the ability to write. A pen, a keyboard, and a printer all do it differently, but they all \"write\". Any type with a <code>Write()</code> method counts.", i.Name)
	case strings.Contains(n, "closer"):
		return fmt.Sprintf("🚪 Think of <code>%s</code> like a door — anything that can be closed. Files, connections, and windows are all different, but they all need closing when you're done.", i.Name)
	case strings.Contains(n, "handler"):
		return fmt.Sprintf("👋 Think of <code>%s</code> like a job description: \"must be able to handle requests.\" Any employee (type) that can do this job qualifies — doesn't matter if they're a robot or a human.", i.Name)
	case strings.Contains(n, "stringer") || strings.Contains(n, "string"):
		return fmt.Sprintf("🏷️ Think of <code>%s</code> like a name tag. Anything that can introduce itself (turn itself into text) satisfies this — \"Hi, I'm a Server\" or \"Hi, I'm Error #404\".", i.Name)
	case strings.Contains(n, "error"):
		return fmt.Sprintf("🚨 Think of <code>%s</code> like the ability to explain what went wrong. A flat tire and a dead battery are different problems, but both can tell you \"here's what happened.\"", i.Name)
	case strings.Contains(n, "formatter") || strings.Contains(n, "format"):
		return fmt.Sprintf("🎨 Think of <code>%s</code> like the ability to dress up for different occasions. Same person, different outfits — it controls how something presents itself.", i.Name)
	case strings.Contains(n, "scanner") || strings.Contains(n, "scan"):
		return fmt.Sprintf("🔍 Think of <code>%s</code> like the ability to scan a barcode. Different scanners look different, but they all do the same job — read information from input.", i.Name)
	}

	// Generic analogies
	if len(i.Methods) == 0 {
		return fmt.Sprintf("📭 <code>%s</code> is an empty interface — like saying \"anyone is welcome.\" Every type in Go automatically satisfies it. It's the universal container.", i.Name)
	}
	if len(i.Methods) == 1 {
		return fmt.Sprintf("📋 Think of <code>%s</code> as a job requirement with one skill: <code>%s()</code>. Any type that has this skill automatically qualifies — no resume needed, no paperwork. Go checks at compile time.", i.Name, i.Methods[0].Name)
	}
	methodNames := []string{}
	for _, m := range i.Methods {
		methodNames = append(methodNames, m.Name+"()")
	}
	return fmt.Sprintf("📋 Think of <code>%s</code> as a job posting with %d required skills: %s. Any type that has ALL these skills qualifies. Go checks automatically — if it walks like a duck and quacks like a duck, it's a duck.", i.Name, len(i.Methods), strings.Join(methodNames, ", "))
}

func generateFriendlyFunction(f parser.FunctionDoc, pkgName string, importPath ...string) FriendlyFunction {
	impPath := ""
	if len(importPath) > 0 {
		impPath = importPath[0]
	}
	// Check knowledge base for function-specific knowledge
	fk := lookupFuncKnowledge(impPath, f.Name)

	ff := FriendlyFunction{
		Name:      f.Name,
		Signature: f.Signature,
		Receiver:  f.Receiver,
		Returns:   f.Returns,
	}

	// Detect builder pattern: method returns pointer to receiver type
	isBuilder := isBuilderMethod(f)

	if f.Doc != "" {
		ff.Overview = cleanDoc(f.Doc)
		ff.WhatItDoes = cleanDoc(f.Doc)
	} else {
		var whatPoints []string
		if isBuilder {
			whatPoints = append(whatPoints, fmt.Sprintf("<code>%s</code> is a <strong>builder method</strong> — it sets a value and returns the same object so you can chain calls.", f.Name))
		} else if f.Receiver != "" {
			whatPoints = append(whatPoints, fmt.Sprintf("<code>%s</code> is a method on <code>%s</code> in the <code>%s</code> package.", f.Name, f.Receiver, pkgName))
		} else {
			whatPoints = append(whatPoints, fmt.Sprintf("<code>%s</code> is a function in the <code>%s</code> package.", f.Name, pkgName))
		}
		if len(f.Params) > 0 {
			paramDescs := []string{}
			for _, p := range f.Params {
				if p.Name != "" {
					paramDescs = append(paramDescs, fmt.Sprintf("<code>%s</code> (%s)", p.Name, p.Type))
				} else {
					paramDescs = append(paramDescs, fmt.Sprintf("<code>%s</code>", p.Type))
				}
			}
			whatPoints = append(whatPoints, fmt.Sprintf("Accepts: %s", strings.Join(paramDescs, ", ")))
		}
		if len(f.Returns) > 0 && !isBuilder {
			whatPoints = append(whatPoints, fmt.Sprintf("Returns: <code>%s</code>", strings.Join(f.Returns, ", ")))
		}
		if isBuilder {
			whatPoints = append(whatPoints, "Returns the same object (<strong>builder pattern</strong>) — you can chain multiple calls together.")
		}
		ff.Overview = buildPointList(whatPoints)
		ff.WhatItDoes = buildPointList(whatPoints)
	}

	// SUPPLEMENTAL: Knowledge base adds tips AFTER the real doc, never replaces it
	if fk != nil && fk.WhatItDoes != "" {
		ff.WhatItDoes = ff.WhatItDoes + "<p>💡 <em>" + fk.WhatItDoes + "</em></p>"
	}

	// Build descriptive WhenToUse — doc-based first, knowledge base supplements
	var whenPoints []string
	if f.Doc != "" {
		// Derive "when to use" from the function's own doc comment
		first := strings.TrimSpace(strings.Split(f.Doc, "\n")[0])
		first = strings.TrimSuffix(first, ".")
		if first != "" && len(first) > 10 {
			whenPoints = append(whenPoints, fmt.Sprintf("Use <code>%s</code> when you need to <strong>%s</strong>.", f.Name, strings.ToLower(first)))
		}
	}
	// Add knowledge base tips as supplement
	if fk != nil && fk.WhenToUse != "" {
		whenPoints = append(whenPoints, "💡 "+fk.WhenToUse)
		if fk.HowToUse != "" {
			whenPoints = append(whenPoints, fmt.Sprintf("📝 <strong>How to use:</strong><pre class=\"task-code\"><code>%s</code></pre>", fk.HowToUse))
		}
	}
	// Fallback if nothing from doc or knowledge base
	if len(whenPoints) == 0 {
		whenPoints = generateFuncWhenToUse(f, pkgName, isBuilder)
	}
	ff.WhenToUse = buildPointList(whenPoints)

	// Params
	for _, p := range f.Params {
		fp := FriendlyParam{
			Name: p.Name,
			Type: p.Type,
		}
		if p.Name != "" {
			fp.Description = fmt.Sprintf("<code>%s</code> — a value of type <code>%s</code>.", p.Name, p.Type)
		} else {
			fp.Description = fmt.Sprintf("A value of type <code>%s</code>.", p.Type)
		}
		ff.Params = append(ff.Params, fp)
	}

	// Generate example
	ff.Example = generateFuncExample(f, pkgName)

	// In Simple Words — real-world analogy
	ff.Explanation = generateFuncAnalogy(f)

	return ff
}

// generateFuncWhenToUse creates descriptive, practical "When To Use" guidance for a function
func generateFuncWhenToUse(f parser.FunctionDoc, pkgName string, isBuilder bool) []string {
	n := strings.ToLower(f.Name)
	var points []string

	if isBuilder {
		points = append(points, fmt.Sprintf("Call <code>.%s(value)</code> when <strong>building or configuring</strong> an object step by step.", f.Name))
		points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Chain it with other builder methods: <code>obj.%s(val1).AnotherMethod(val2)</code>", f.Name))
		points = append(points, "💡 Each call sets one setting and returns the object — like filling in a form field by field.")
		return points
	}

	// Context-specific guidance based on function name
	switch {
	case strings.HasPrefix(n, "new") || strings.HasPrefix(n, "create"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>create a new instance</strong> ready to use.", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Call it once, store the result, then use the returned object.")
		points = append(points, "💡 This is the recommended way to create this type — don't use <code>{}</code> directly unless the docs say so.")
	case strings.HasPrefix(n, "get") || strings.HasPrefix(n, "fetch") || strings.HasPrefix(n, "read") || strings.HasPrefix(n, "load"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>retrieve data</strong> — it looks something up and gives it back to you.", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Pass in what you're looking for → get back the result (check for errors if returned).")
	case strings.HasPrefix(n, "set") || strings.HasPrefix(n, "update") || strings.HasPrefix(n, "put"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>change or update</strong> a value.", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Pass the new value → the old one gets replaced.")
	case strings.HasPrefix(n, "delete") || strings.HasPrefix(n, "remove") || strings.HasPrefix(n, "clear"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>remove something</strong>.", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Pass what to remove → it's gone. Check for errors if the operation can fail.")
	case strings.HasPrefix(n, "list") || strings.HasPrefix(n, "find") || strings.HasPrefix(n, "search") || strings.HasPrefix(n, "filter"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>search or list items</strong> that match your criteria.", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Pass your search criteria → get back a list of results.")
	case strings.HasPrefix(n, "close") || strings.HasPrefix(n, "shutdown") || strings.HasPrefix(n, "stop") || strings.HasPrefix(n, "cancel"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>clean up and release resources</strong>.", pkgName, f.Name))
		points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Call it when you're done. Best practice: use <code>defer obj.%s()</code> right after creation.", f.Name))
		points = append(points, "⚠️ <strong>Important:</strong> Forgetting to close can leak memory, connections, or file handles.")
	case strings.HasPrefix(n, "start") || strings.HasPrefix(n, "run") || strings.HasPrefix(n, "serve") || strings.HasPrefix(n, "listen"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>start a long-running process</strong> (server, listener, worker).", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> This usually blocks — put it at the end of main, or run it in a goroutine with <code>go</code>.")
		points = append(points, "💡 It keeps running until stopped or an error occurs.")
	case strings.HasPrefix(n, "write") || strings.HasPrefix(n, "print") || strings.HasPrefix(n, "log") || strings.HasPrefix(n, "emit"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>output data</strong> (to screen, file, or network).", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Pass the data you want to write → it goes to the destination.")
	case strings.HasPrefix(n, "parse") || strings.HasPrefix(n, "decode") || strings.HasPrefix(n, "unmarshal"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>convert raw data into a Go value</strong> (JSON → struct, string → time, etc.).", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Pass the raw input → get back a usable Go value. Always check the error.")
	case strings.HasPrefix(n, "encode") || strings.HasPrefix(n, "marshal") || strings.HasPrefix(n, "format") || strings.HasPrefix(n, "sprint"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>convert a Go value into a formatted output</strong> (struct → JSON, values → string, etc.).", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Pass your Go value → get back formatted bytes or string.")
	case strings.HasPrefix(n, "handle") || strings.HasPrefix(n, "register") || strings.HasPrefix(n, "add"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>register a handler or add an item</strong>.", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Call it during setup (before starting). The registered handler gets called later automatically.")
	case strings.HasPrefix(n, "is") || strings.HasPrefix(n, "has") || strings.HasPrefix(n, "can") || strings.HasPrefix(n, "check") || strings.HasPrefix(n, "valid"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>check a condition</strong> — returns true or false.", pkgName, f.Name))
		points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Use it in an <code>if</code> statement: <code>if %s.%s(...) { ... }</code>", pkgName, f.Name))
	case strings.HasPrefix(n, "with"):
		points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>create a copy with a modified setting</strong>.", pkgName, f.Name))
		points = append(points, "📝 <strong>How:</strong> Pass the new value → get back a modified version (original stays unchanged).")
	default:
		// Still descriptive for unknown functions
		if f.Receiver != "" {
			receiverBase := strings.TrimPrefix(f.Receiver, "*")
			if len(f.Returns) > 0 && len(f.Params) > 0 {
				points = append(points, fmt.Sprintf("Call <code>.%s()</code> on a <code>%s</code> to <strong>process input and get a result</strong>.", f.Name, receiverBase))
				points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>result := my%s.%s(args)</code>", receiverBase, f.Name))
			} else if len(f.Returns) > 0 {
				points = append(points, fmt.Sprintf("Call <code>.%s()</code> on a <code>%s</code> to <strong>get information</strong> from it.", f.Name, receiverBase))
				points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>value := my%s.%s()</code>", receiverBase, f.Name))
			} else if len(f.Params) > 0 {
				points = append(points, fmt.Sprintf("Call <code>.%s()</code> on a <code>%s</code> to <strong>perform an action</strong> with the given input.", f.Name, receiverBase))
				points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>my%s.%s(args)</code>", receiverBase, f.Name))
			} else {
				points = append(points, fmt.Sprintf("Call <code>.%s()</code> on a <code>%s</code> to <strong>trigger an action</strong>.", f.Name, receiverBase))
				points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>my%s.%s()</code>", receiverBase, f.Name))
			}
		} else {
			if len(f.Returns) > 0 && len(f.Params) > 0 {
				points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>process input and get a result</strong>.", pkgName, f.Name))
				points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>result := %s.%s(args)</code>", pkgName, f.Name))
			} else if len(f.Returns) > 0 {
				points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>get a value</strong>.", pkgName, f.Name))
				points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>value := %s.%s()</code>", pkgName, f.Name))
			} else {
				points = append(points, fmt.Sprintf("Call <code>%s.%s()</code> to <strong>perform an action</strong>.", pkgName, f.Name))
				points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>%s.%s(args)</code>", pkgName, f.Name))
			}
		}
	}

	// Add error handling tip if function returns error
	for _, ret := range f.Returns {
		if ret == "error" {
			points = append(points, "⚠️ <strong>Always check the error:</strong> <code>if err != nil { log.Fatal(err) }</code>")
			break
		}
	}

	return points
}

// generateInterfaceWhenToUse creates descriptive "When To Use" guidance for an interface
func generateInterfaceWhenToUse(i parser.InterfaceDoc, pkgName string) []string {
	n := strings.ToLower(i.Name)
	var points []string

	switch {
	case strings.Contains(n, "reader"):
		points = append(points, fmt.Sprintf("Implement <code>%s</code> when your type can <strong>provide data to be read</strong> (files, network, buffers).", i.Name))
		points = append(points, "📝 <strong>How:</strong> Add a <code>Read(p []byte) (n int, err error)</code> method to your type.")
		points = append(points, "💡 Then pass your type anywhere a <code>Reader</code> is accepted — Go figures it out automatically.")
	case strings.Contains(n, "writer"):
		points = append(points, fmt.Sprintf("Implement <code>%s</code> when your type can <strong>accept data to be written</strong> (files, HTTP responses, loggers).", i.Name))
		points = append(points, "📝 <strong>How:</strong> Add a <code>Write(p []byte) (n int, err error)</code> method to your type.")
		points = append(points, "💡 Then pass your type anywhere a <code>Writer</code> is accepted.")
	case strings.Contains(n, "closer"):
		points = append(points, fmt.Sprintf("Implement <code>%s</code> when your type <strong>holds resources that need cleanup</strong> (connections, files, handles).", i.Name))
		points = append(points, "📝 <strong>How:</strong> Add a <code>Close() error</code> method → users call <code>defer x.Close()</code>.")
	case strings.Contains(n, "handler"):
		points = append(points, fmt.Sprintf("Implement <code>%s</code> to <strong>handle incoming requests or events</strong> with your own logic.", i.Name))
		points = append(points, "📝 <strong>How:</strong> Create a struct → add the required handler method → register it with the server/router.")
		points = append(points, "💡 This is how you add custom behavior — each handler does something different with the same type of request.")
	case strings.Contains(n, "stringer") || strings.Contains(n, "string"):
		points = append(points, fmt.Sprintf("Implement <code>%s</code> so your type can <strong>describe itself as text</strong> (useful for printing and logging).", i.Name))
		points = append(points, "📝 <strong>How:</strong> Add a <code>String() string</code> method → <code>fmt.Println(myObj)</code> will use it automatically.")
	case strings.Contains(n, "error"):
		points = append(points, fmt.Sprintf("Implement <code>%s</code> to create <strong>custom error types</strong> with extra details (codes, context, nested errors).", i.Name))
		points = append(points, "📝 <strong>How:</strong> Add an <code>Error() string</code> method → your type works with <code>if err != nil</code> checks.")
	case strings.Contains(n, "sort") || strings.Contains(n, "less"):
		points = append(points, fmt.Sprintf("Implement <code>%s</code> to make your type <strong>sortable</strong>.", i.Name))
		points = append(points, "📝 <strong>How:</strong> Add the required methods → pass your collection to <code>sort.Sort()</code>.")
	case strings.Contains(n, "marshal"):
		points = append(points, fmt.Sprintf("Implement <code>%s</code> to <strong>control how your type converts to/from data formats</strong> (JSON, XML, etc.).", i.Name))
		points = append(points, "📝 <strong>How:</strong> Add the marshal/unmarshal method → the encoder/decoder calls it automatically.")
	default:
		// Generic but descriptive
		if len(i.Methods) == 0 {
			points = append(points, fmt.Sprintf("Use <code>%s</code> (empty interface) to accept <strong>any type</strong> — like a universal container.", i.Name))
			points = append(points, "📝 <strong>How:</strong> Use it as a parameter type when you don't know the exact type in advance.")
		} else if len(i.Methods) == 1 {
			m := i.Methods[0]
			points = append(points, fmt.Sprintf("Implement <code>%s</code> by adding a <code>%s()</code> method to your type — that's all it takes.", i.Name, m.Name))
			points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>func (t YourType) %s(...) { ... }</code> → now your type works anywhere <code>%s</code> is expected.", m.Name, i.Name))
			points = append(points, "💡 Go checks this automatically at compile time — no need to explicitly declare that you implement it.")
		} else {
			methodNames := []string{}
			for _, m := range i.Methods {
				if len(methodNames) < 3 {
					methodNames = append(methodNames, "<code>"+m.Name+"()</code>")
				}
			}
			points = append(points, fmt.Sprintf("Implement <code>%s</code> by adding these methods to your type: %s.", i.Name, strings.Join(methodNames, ", ")))
			points = append(points, "📝 <strong>How:</strong> Create a struct → add each required method → your type automatically satisfies the interface.")
			points = append(points, "💡 You don't need to declare <code>implements</code> — Go checks it for you at compile time.")
		}
	}

	return points
}

// generateTypeWhenToUse creates descriptive "When To Use" guidance for a named type
func generateTypeWhenToUse(t parser.TypeDoc, pkgName string) []string {
	n := strings.ToLower(t.Name)
	u := strings.ToLower(t.Underlying)
	var points []string

	switch {
	case strings.Contains(n, "status") || strings.Contains(n, "phase") || strings.Contains(n, "state"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to represent a <strong>status or state</strong> value. The package uses this instead of raw strings/ints to prevent typos.", pkgName, t.Name))
		points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Use the predefined constants (e.g., <code>%s.StatusRunning</code>) — don't make up your own values.", pkgName))
	case strings.Contains(n, "method") || strings.Contains(n, "verb") || strings.Contains(n, "action"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to specify an <strong>action or method</strong> (GET, POST, etc.). Safer than raw strings.", pkgName, t.Name))
		points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Use the package constants: <code>%s.MethodGet</code>, <code>%s.MethodPost</code>, etc.", pkgName, pkgName))
	case strings.Contains(n, "type") || strings.Contains(n, "kind"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> to identify the <strong>kind or category</strong> of something. Prevents mixing up with unrelated values.", pkgName, t.Name))
		points = append(points, "📝 <strong>How:</strong> Use the predefined constants — check the docs for available values.")
	case strings.Contains(n, "func") || strings.Contains(u, "func"):
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> as a <strong>callback or handler function</strong>. Define your own function matching this signature and pass it.", pkgName, t.Name))
		points = append(points, fmt.Sprintf("📝 <strong>How:</strong> <code>var myHandler %s.%s = func(...) { ... }</code>", pkgName, t.Name))
	default:
		points = append(points, fmt.Sprintf("Use <code>%s.%s</code> instead of raw <code>%s</code> when the package expects this specific type — it adds <strong>type safety</strong>.", pkgName, t.Name, t.Underlying))
		points = append(points, fmt.Sprintf("📝 <strong>How:</strong> Convert with <code>%s.%s(yourValue)</code> or use the predefined constants.", pkgName, t.Name))
		points = append(points, fmt.Sprintf("💡 <strong>Why not just use <code>%s</code>?</strong> Because Go won't let you accidentally pass a wrong value — the compiler catches mistakes.", t.Underlying))
	}

	return points
}

// generateFuncAnalogy creates a real-world analogy for a function
func generateFuncAnalogy(f parser.FunctionDoc) string {
	n := strings.ToLower(f.Name)
	isBuilder := isBuilderMethod(f)

	if isBuilder {
		return fmt.Sprintf("🔗 <code>.%s()</code> is like filling in one field on an order form. You write something down, get the form back, and can fill in the next field — that's method chaining.", f.Name)
	}

	// Name-based analogies
	switch {
	case strings.Contains(n, "new") || strings.HasPrefix(n, "create"):
		return fmt.Sprintf("🏭 <code>%s</code> is like a factory — you call it, and it builds a new thing for you, ready to use.", f.Name)
	case strings.HasPrefix(n, "get") || strings.HasPrefix(n, "fetch") || strings.HasPrefix(n, "read"):
		return fmt.Sprintf("📬 <code>%s</code> is like checking your mailbox — you ask for something, and it gives you what's there (or tells you it's empty).", f.Name)
	case strings.HasPrefix(n, "set") || strings.HasPrefix(n, "update"):
		return fmt.Sprintf("✏️ <code>%s</code> is like updating a contact in your phone — you pick the field and change its value.", f.Name)
	case strings.HasPrefix(n, "delete") || strings.HasPrefix(n, "remove"):
		return fmt.Sprintf("🗑️ <code>%s</code> is like throwing something in the trash — once called, the item is gone.", f.Name)
	case strings.HasPrefix(n, "list") || strings.HasPrefix(n, "find") || strings.HasPrefix(n, "search"):
		return fmt.Sprintf("🔍 <code>%s</code> is like searching through a filing cabinet — it looks through everything and returns what matches.", f.Name)
	case strings.HasPrefix(n, "close") || strings.HasPrefix(n, "shutdown") || strings.HasPrefix(n, "stop"):
		return fmt.Sprintf("🚪 <code>%s</code> is like closing a store for the night — it wraps up any ongoing work and locks the doors.", f.Name)
	case strings.HasPrefix(n, "start") || strings.HasPrefix(n, "run") || strings.HasPrefix(n, "serve") || strings.HasPrefix(n, "listen"):
		return fmt.Sprintf("🚀 <code>%s</code> is like pressing the power button — it starts things up and keeps running until you tell it to stop.", f.Name)
	case strings.HasPrefix(n, "write") || strings.HasPrefix(n, "print") || strings.HasPrefix(n, "log"):
		return fmt.Sprintf("📝 <code>%s</code> is like writing on a whiteboard — it takes your data and puts it somewhere visible (screen, file, network).", f.Name)
	case strings.HasPrefix(n, "parse") || strings.HasPrefix(n, "decode"):
		return fmt.Sprintf("🔓 <code>%s</code> is like translating a foreign language — it takes raw data (text, JSON, etc.) and turns it into something your code understands.", f.Name)
	case strings.HasPrefix(n, "encode") || strings.HasPrefix(n, "marshal"):
		return fmt.Sprintf("📦 <code>%s</code> is like packing a suitcase — it takes your data and packs it into a format that can be sent or stored.", f.Name)
	case strings.HasPrefix(n, "format") || strings.HasPrefix(n, "sprint"):
		return fmt.Sprintf("🎨 <code>%s</code> is like formatting a document — it takes raw information and makes it look nice as text.", f.Name)
	case strings.HasPrefix(n, "handle"):
		return fmt.Sprintf("👋 <code>%s</code> is like a receptionist — when something arrives, this function decides what to do with it.", f.Name)
	case strings.HasPrefix(n, "register"):
		return fmt.Sprintf("📝 <code>%s</code> is like signing up for a service — you register something so the system knows about it later.", f.Name)
	case strings.HasPrefix(n, "validate") || strings.HasPrefix(n, "check") || strings.HasPrefix(n, "verify"):
		return fmt.Sprintf("✅ <code>%s</code> is like checking IDs at the door — it looks at the data and tells you if everything is valid.", f.Name)
	case strings.HasPrefix(n, "convert") || strings.HasPrefix(n, "to"):
		return fmt.Sprintf("🔄 <code>%s</code> is like currency exchange — it takes one form of data and converts it to another.", f.Name)
	case strings.HasPrefix(n, "init") || strings.HasPrefix(n, "setup"):
		return fmt.Sprintf("🔧 <code>%s</code> is like setting up a new phone — it runs once at the beginning to get everything ready.", f.Name)
	case strings.HasPrefix(n, "wait") || strings.HasPrefix(n, "sleep"):
		return fmt.Sprintf("⏳ <code>%s</code> is like waiting in line — it pauses until something happens or a condition is met.", f.Name)
	case strings.HasPrefix(n, "send") || strings.HasPrefix(n, "emit") || strings.HasPrefix(n, "notify"):
		return fmt.Sprintf("📤 <code>%s</code> is like sending a text message — it delivers data to someone else.", f.Name)
	}

	// Receiver-based fallback
	if f.Receiver != "" {
		return fmt.Sprintf("🔧 First create a <code>%s</code> object, then call <code>.%s()</code> on it — like pressing a button on a device you've already set up.", f.Receiver, f.Name)
	}

	// Generic
	if len(f.Returns) > 0 && len(f.Params) > 0 {
		return fmt.Sprintf("⚡ You give <code>%s</code> some input, and it gives you something back — like a vending machine.", f.Name)
	}
	if len(f.Returns) > 0 {
		return fmt.Sprintf("📦 Call <code>%s</code> and it gives you something back, ready to use.", f.Name)
	}
	if len(f.Params) > 0 {
		return fmt.Sprintf("📥 Give <code>%s</code> some data and it takes care of the rest.", f.Name)
	}
	return fmt.Sprintf("⚡ <code>%s</code> performs an action when you call it.", f.Name)
}

// isBuilderMethod detects if a method follows the builder pattern:
// - has a receiver
// - returns a pointer to the receiver type
// - typically named With*, Set*, Add*
func isBuilderMethod(f parser.FunctionDoc) bool {
	if f.Receiver == "" {
		return false
	}
	// Check if return type matches receiver (pointer to receiver)
	receiverBase := strings.TrimPrefix(f.Receiver, "*")
	for _, ret := range f.Returns {
		retBase := strings.TrimPrefix(ret, "*")
		if retBase == receiverBase {
			return true
		}
	}
	// Also detect by naming convention
	if strings.HasPrefix(f.Name, "With") || strings.HasPrefix(f.Name, "Set") {
		return true
	}
	return false
}

func generateStructExample(s parser.StructDoc, pkgName string) string {
	var sb strings.Builder

	// Collect exported fields
	exportedFields := []parser.FieldDoc{}
	for _, f := range s.Fields {
		if f.Name != "" && len(f.Name) > 0 && f.Name[0] >= 'A' && f.Name[0] <= 'Z' {
			exportedFields = append(exportedFields, f)
		}
	}

	// Hello-World style: numbered steps, plain English, print at end
	sb.WriteString("package main\n\n")
	sb.WriteString(fmt.Sprintf("import \"%s\"\n", pkgName))
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString(fmt.Sprintf("    // STEP 1: Create a %s (like filling out a form)\n", s.Name))
	sb.WriteString(fmt.Sprintf("    my%s := %s.%s{\n", s.Name, pkgName, s.Name))

	// Show max 3 fields with friendly values and comments
	shown := 0
	for _, f := range exportedFields {
		if shown >= 3 {
			sb.WriteString("        // ... you can set more fields if needed\n")
			break
		}
		sample := sampleValue(f.Name, f.Type)
		comment := friendlyFieldComment(f.Name, f.Type)
		sb.WriteString(fmt.Sprintf("        %s: %s, %s\n", f.Name, sample, comment))
		shown++
	}
	sb.WriteString("    }\n\n")

	// If has methods, show calling one with explanation
	if len(s.Methods) > 0 {
		m := s.Methods[0]
		callArgs := []string{}
		for _, p := range m.Params {
			callArgs = append(callArgs, sampleValue(p.Name, p.Type))
		}
		sb.WriteString(fmt.Sprintf("    // STEP 2: Use it! Call .%s()\n", m.Name))
		if isBuilderMethod(m) {
			sb.WriteString(fmt.Sprintf("    my%s = my%s.%s(%s)  // set a value and get it back\n", s.Name, s.Name, m.Name, strings.Join(callArgs, ", ")))
		} else if len(m.Returns) > 0 {
			sb.WriteString(fmt.Sprintf("    result := my%s.%s(%s)\n", s.Name, m.Name, strings.Join(callArgs, ", ")))
			sb.WriteString("    fmt.Println(\"Result:\", result)  // see what we got\n")
		} else {
			sb.WriteString(fmt.Sprintf("    my%s.%s(%s)\n", s.Name, m.Name, strings.Join(callArgs, ", ")))
		}
	} else {
		sb.WriteString(fmt.Sprintf("    // STEP 2: Print it to see what's inside\n"))
		sb.WriteString(fmt.Sprintf("    fmt.Println(my%s)\n", s.Name))
	}

	sb.WriteString("}")

	return sb.String()
}

// friendlyFieldComment returns a very simple comment explaining what a field is
func friendlyFieldComment(name string, typeName string) string {
	n := strings.ToLower(name)
	switch {
	case strings.Contains(n, "addr") || strings.Contains(n, "host"):
		return "// where to connect"
	case strings.Contains(n, "name"):
		return "// give it a name"
	case strings.Contains(n, "port"):
		return "// which port to use"
	case strings.Contains(n, "timeout"):
		return "// how long to wait"
	case strings.Contains(n, "handler"):
		return "// who handles the work"
	case strings.Contains(n, "path") || strings.Contains(n, "url"):
		return "// the address/path"
	case strings.Contains(n, "size") || strings.Contains(n, "max") || strings.Contains(n, "limit"):
		return "// how big/how many"
	case strings.Contains(n, "enabled") || strings.Contains(n, "disabled"):
		return "// turn on or off"
	case strings.Contains(n, "config"):
		return "// settings"
	case typeName == "string":
		return "// text value"
	case typeName == "int" || typeName == "int32" || typeName == "int64":
		return "// a number"
	case typeName == "bool":
		return "// true or false"
	default:
		return ""
	}
}

func generateInterfaceExample(i parser.InterfaceDoc, pkgName string) string {
	var sb strings.Builder

	sb.WriteString("package main\n\n")
	sb.WriteString(fmt.Sprintf("import \"%s\"\n", pkgName))
	sb.WriteString("import \"fmt\"\n\n")

	sb.WriteString(fmt.Sprintf("// STEP 1: Create your own type (like creating a new character)\n"))
	sb.WriteString(fmt.Sprintf("type My%s struct {\n", i.Name))
	sb.WriteString("    Name string  // give it some data\n")
	sb.WriteString("}\n\n")

	sb.WriteString(fmt.Sprintf("// STEP 2: Add the required method(s) — this makes it count as a %s\n", i.Name))
	for _, m := range i.Methods {
		params := []string{}
		for _, p := range m.Params {
			if p.Name != "" {
				params = append(params, p.Name+" "+p.Type)
			} else {
				params = append(params, p.Type)
			}
		}
		retStr := ""
		if len(m.Returns) > 0 {
			if len(m.Returns) == 1 {
				retStr = " " + m.Returns[0]
			} else {
				retStr = " (" + strings.Join(m.Returns, ", ") + ")"
			}
		}
		sb.WriteString(fmt.Sprintf("func (m My%s) %s(%s)%s {\n", i.Name, m.Name, strings.Join(params, ", "), retStr))
		sb.WriteString(fmt.Sprintf("    fmt.Println(\"Hello from\", m.Name)  // do your thing here\n"))
		if len(m.Returns) > 0 {
			if m.Returns[0] == "error" {
				sb.WriteString("    return nil  // no error\n")
			} else if m.Returns[0] == "int" {
				sb.WriteString("    return 0\n")
			} else if m.Returns[0] == "string" {
				sb.WriteString("    return m.Name\n")
			} else if m.Returns[0] == "bool" {
				sb.WriteString("    return true\n")
			}
		}
		sb.WriteString("}\n\n")
	}

	sb.WriteString(fmt.Sprintf("// STEP 3: Use it! Your type now works wherever %s is needed\n", i.Name))
	sb.WriteString("func main() {\n")
	sb.WriteString(fmt.Sprintf("    mine := My%s{Name: \"World\"}\n", i.Name))
	if len(i.Methods) > 0 {
		m := i.Methods[0]
		args := []string{}
		for _, p := range m.Params {
			args = append(args, sampleValue(p.Name, p.Type))
		}
		if len(m.Returns) > 0 {
			sb.WriteString(fmt.Sprintf("    result := mine.%s(%s)  // call the method!\n", m.Name, strings.Join(args, ", ")))
			sb.WriteString("    fmt.Println(\"Got:\", result)\n")
		} else {
			sb.WriteString(fmt.Sprintf("    mine.%s(%s)  // call the method!\n", m.Name, strings.Join(args, ", ")))
		}
	}
	sb.WriteString("}")

	return sb.String()
}

func generateFuncExample(f parser.FunctionDoc, pkgName string) string {
	var sb strings.Builder

	// Generate sample argument values
	callArgs := []string{}
	for _, p := range f.Params {
		callArgs = append(callArgs, sampleValue(p.Name, p.Type))
	}
	argStr := strings.Join(callArgs, ", ")

	isBuilder := isBuilderMethod(f)

	// Hello-World style: complete mini-program with comments on every line
	sb.WriteString("package main\n\n")
	sb.WriteString(fmt.Sprintf("import \"%s\"\n", pkgName))
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")

	if isBuilder {
		receiverBase := strings.TrimPrefix(f.Receiver, "*")
		sb.WriteString(fmt.Sprintf("    // STEP 1: Create an empty %s\n", receiverBase))
		sb.WriteString(fmt.Sprintf("    obj := &%s.%s{}\n\n", pkgName, receiverBase))
		sb.WriteString(fmt.Sprintf("    // STEP 2: Fill in a value using .%s()\n", f.Name))
		sb.WriteString(fmt.Sprintf("    obj = obj.%s(%s)\n", f.Name, argStr))
		sb.WriteString("    // ↑ This returns the same object, so you can keep going:\n")
		sb.WriteString(fmt.Sprintf("    // obj.%s(%s).AnotherMethod(value)\n\n", f.Name, argStr))
		sb.WriteString(fmt.Sprintf("    fmt.Println(\"Done! %s is set.\")  // see the result\n", f.Name))
	} else if f.Receiver != "" {
		receiverBase := strings.TrimPrefix(f.Receiver, "*")
		sb.WriteString(fmt.Sprintf("    // STEP 1: Create a %s first\n", receiverBase))
		sb.WriteString(fmt.Sprintf("    obj := %s.%s{}  // start with defaults\n\n", pkgName, receiverBase))
		sb.WriteString(fmt.Sprintf("    // STEP 2: Call .%s() on it\n", f.Name))
		if len(f.Returns) > 0 {
			sb.WriteString(fmt.Sprintf("    result := obj.%s(%s)\n", f.Name, argStr))
			sb.WriteString("    fmt.Println(\"Result:\", result)  // see what we got\n")
		} else {
			sb.WriteString(fmt.Sprintf("    obj.%s(%s)\n", f.Name, argStr))
			sb.WriteString("    fmt.Println(\"Done!\")  // it worked\n")
		}
	} else {
		// Standalone function — simplest case
		sb.WriteString(fmt.Sprintf("    // Just call %s.%s() — that's it!\n", pkgName, f.Name))
		if len(f.Returns) > 0 {
			if len(f.Returns) == 1 && f.Returns[0] == "error" {
				sb.WriteString(fmt.Sprintf("    err := %s.%s(%s)\n", pkgName, f.Name, argStr))
				sb.WriteString("    if err != nil {\n")
				sb.WriteString("        fmt.Println(\"Oops:\", err)  // something went wrong\n")
				sb.WriteString("    }\n")
				sb.WriteString("    fmt.Println(\"Done!\")  // it worked\n")
			} else if len(f.Returns) >= 2 && f.Returns[len(f.Returns)-1] == "error" {
				sb.WriteString(fmt.Sprintf("    result, err := %s.%s(%s)\n", pkgName, f.Name, argStr))
				sb.WriteString("    if err != nil {\n")
				sb.WriteString("        fmt.Println(\"Oops:\", err)  // something went wrong\n")
				sb.WriteString("    }\n")
				sb.WriteString("    fmt.Println(\"Got:\", result)  // here's your answer\n")
			} else {
				sb.WriteString(fmt.Sprintf("    result := %s.%s(%s)\n", pkgName, f.Name, argStr))
				sb.WriteString("    fmt.Println(\"Got:\", result)  // here's your answer\n")
			}
		} else {
			sb.WriteString(fmt.Sprintf("    %s.%s(%s)\n", pkgName, f.Name, argStr))
			sb.WriteString("    fmt.Println(\"Done!\")  // it worked\n")
		}
	}

	sb.WriteString("}")
	return sb.String()
}

// sampleValue generates a realistic-looking sample value for examples
func sampleValue(name string, typeName string) string {
	// Use the parameter name for context
	n := strings.ToLower(name)

	switch typeName {
	case "string":
		switch {
		case strings.Contains(n, "name"):
			return `"my-resource"`
		case strings.Contains(n, "namespace") || strings.Contains(n, "ns"):
			return `"default"`
		case strings.Contains(n, "label") || strings.Contains(n, "key"):
			return `"app"`
		case strings.Contains(n, "value") || strings.Contains(n, "val"):
			return `"my-app"`
		case strings.Contains(n, "expr") || strings.Contains(n, "expression"):
			return `"self.name == 'test'"`
		case strings.Contains(n, "host"):
			return `"localhost"`
		case strings.Contains(n, "path") || strings.Contains(n, "url"):
			return `"/api/v1/resources"`
		case strings.Contains(n, "image"):
			return `"nginx:latest"`
		default:
			return `"example-value"`
		}
	case "int", "int32", "int64":
		if strings.Contains(n, "port") {
			return "8080"
		}
		if strings.Contains(n, "replica") || strings.Contains(n, "count") {
			return "3"
		}
		return "1"
	case "bool":
		return "true"
	case "float32", "float64":
		return "1.0"
	case "error":
		return "nil"
	default:
		if strings.HasPrefix(typeName, "*") {
			return "nil"
		}
		if strings.HasPrefix(typeName, "[]") {
			return typeName + "{}"
		}
		if strings.HasPrefix(typeName, "map[") {
			return typeName + "{}"
		}
		if name != "" {
			return "/* " + name + " */"
		}
		return "/* " + typeName + " */"
	}
}

// collectNames collects up to maxCount names from a slice, formatting as code tags
func collectNames[T any](items []T, maxCount int, getName func(T) string) string {
	names := []string{}
	for _, item := range items {
		name := getName(item)
		if name != "" {
			names = append(names, "<code>"+name+"</code>")
		}
		if len(names) >= maxCount {
			break
		}
	}
	result := strings.Join(names, ", ")
	if len(items) > maxCount {
		result += fmt.Sprintf(" and %d more", len(items)-maxCount)
	}
	return result
}

func cleanDoc(doc string) string {
	doc = strings.TrimSpace(doc)
	if doc == "" {
		return ""
	}

	// Remove excessive newlines
	for strings.Contains(doc, "\n\n\n") {
		doc = strings.ReplaceAll(doc, "\n\n\n", "\n\n")
	}

	// Split into paragraphs and format as clean bullet points
	return formatAsPoints(doc)
}

// formatAsPoints converts a block of text into clean HTML bullet points.
// Each paragraph becomes one bullet point. Short docs stay as plain text.
func formatAsPoints(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	// Split by double newlines (paragraphs)
	paragraphs := strings.Split(text, "\n\n")
	var points []string

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		// Join wrapped lines within a paragraph into a single line
		lines := strings.Split(para, "\n")
		var joined []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Skip section headers like "# Printing" — they add noise
			if strings.HasPrefix(line, "#") {
				continue
			}
			// Skip lines that look like format tables or code blocks
			if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "    ") {
				continue
			}
			joined = append(joined, line)
		}

		if len(joined) == 0 {
			continue
		}

		combined := strings.Join(joined, " ")

		// Clean up: remove leading bullets/dashes
		combined = strings.TrimLeft(combined, "•-* ")
		combined = strings.TrimSpace(combined)
		if combined == "" || len(combined) < 5 {
			continue
		}

		// Truncate very long paragraphs to keep it simple
		if len(combined) > 200 {
			// Take first sentence only
			idx := strings.Index(combined, ". ")
			if idx > 0 && idx < 200 {
				combined = combined[:idx+1]
			} else {
				combined = combined[:197] + "..."
			}
		}

		points = append(points, combined)
	}

	if len(points) == 0 {
		return text
	}

	// If only 1 point, return as simple text
	if len(points) == 1 {
		return points[0]
	}

	// Limit to max 6 bullet points for readability
	if len(points) > 6 {
		points = points[:6]
	}

	// Format as HTML bullet list
	var sb strings.Builder
	sb.WriteString("<ul class=\"doc-points\">")
	for _, p := range points {
		sb.WriteString("<li>")
		sb.WriteString(p)
		sb.WriteString("</li>")
	}
	sb.WriteString("</ul>")
	return sb.String()
}

// splitSentences splits text into sentences by newline, keeping each line intact
func splitSentences(text string) []string {
	lines := strings.Split(text, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}

	return result
}

// buildPointList creates an HTML unordered list from a slice of strings
func buildPointList(points []string) string {
	if len(points) == 0 {
		return ""
	}
	if len(points) == 1 {
		return points[0]
	}
	var sb strings.Builder
	sb.WriteString("<ul class=\"doc-points\">")
	for _, p := range points {
		sb.WriteString("<li>")
		sb.WriteString(p)
		sb.WriteString("</li>")
	}
	sb.WriteString("</ul>")
	return sb.String()
}
