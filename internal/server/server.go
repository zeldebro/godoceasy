package server

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/godoceasy/internal/docs"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

// Server is the documentation web server
type Server struct {
	docs        []*docs.FriendlyDoc
	packagePath string
	templates   *template.Template
}

// SearchResult holds a single search result
type SearchResult struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Package     string `json:"package"`
	Description string `json:"description"`
	Anchor      string `json:"anchor"`
}

// New creates a new documentation server
func New(friendlyDocs []*docs.FriendlyDoc, packagePath string) *Server {
	funcMap := template.FuncMap{
		"lower":    strings.ToLower,
		"contains": strings.Contains,
		"join":     strings.Join,
		"add": func(a, b int) int {
			return a + b
		},
		"hasContent": func(s string) bool {
			return strings.TrimSpace(s) != ""
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	return &Server{
		docs:        friendlyDocs,
		packagePath: packagePath,
		templates:   tmpl,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	// Static files
	http.Handle("/static/", http.FileServer(http.FS(staticFS)))

	// Routes
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/pkg/", s.handlePackage)
	http.HandleFunc("/search", s.handleSearch)
	http.HandleFunc("/api/search", s.handleAPISearch)
	http.HandleFunc("/api/diagram/", s.handleDiagram)

	return http.ListenAndServe(addr, nil)
}

// handleIndex shows the main documentation page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title":       "godoceasy - " + s.packagePath,
		"PackagePath": s.packagePath,
		"Packages":    s.docs,
	}

	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Template error: %v", err)
	}
}

// handlePackage shows documentation for a specific package
func (s *Server) handlePackage(w http.ResponseWriter, r *http.Request) {
	pkgName := strings.TrimPrefix(r.URL.Path, "/pkg/")
	pkgName = strings.TrimSuffix(pkgName, "/")

	var pkg *docs.FriendlyDoc
	for _, p := range s.docs {
		if p.PackageName == pkgName || p.ImportPath == pkgName {
			pkg = p
			break
		}
	}

	if pkg == nil && len(s.docs) > 0 {
		pkg = s.docs[0]
	}

	if pkg == nil {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Title":       "godoceasy - " + pkg.PackageName,
		"PackagePath": s.packagePath,
		"Package":     pkg,
		"Packages":    s.docs,
	}

	if err := s.templates.ExecuteTemplate(w, "package.html", data); err != nil {
		log.Printf("Template error: %v", err)
	}
}

// handleSearch shows search results page
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	results := s.search(query)

	data := map[string]interface{}{
		"Title":       "Search - " + query,
		"PackagePath": s.packagePath,
		"Query":       query,
		"Results":     results,
		"Packages":    s.docs,
	}

	if err := s.templates.ExecuteTemplate(w, "search.html", data); err != nil {
		log.Printf("Template error: %v", err)
	}
}

// handleAPISearch returns JSON search results
func (s *Server) handleAPISearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	results := s.search(query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// search searches across ALL content in ALL packages — every word is searchable
func (s *Server) search(query string) []SearchResult {
	if query == "" {
		return nil
	}

	var results []SearchResult
	q := strings.ToLower(query)

	for _, pkg := range s.docs {
		// Search functions
		for _, f := range pkg.Functions {
			if deepMatchFunc(f, q) {
				results = append(results, SearchResult{
					Name:        f.Name,
					Type:        "Function",
					Package:     pkg.PackageName,
					Description: truncate(stripHTML(f.WhatItDoes), 120),
					Anchor:      "func-" + f.Name,
				})
			}
		}

		// Search structs
		for _, st := range pkg.Structs {
			if deepMatchStruct(st, q) {
				results = append(results, SearchResult{
					Name:        st.Name,
					Type:        "Struct",
					Package:     pkg.PackageName,
					Description: truncate(stripHTML(st.WhatItDoes), 120),
					Anchor:      "struct-" + st.Name,
				})
			}
			// Search struct fields individually
			for _, field := range st.Fields {
				if containsAny(q, field.Name, field.Type, field.Description, field.Tag) {
					results = append(results, SearchResult{
						Name:        st.Name + "." + field.Name,
						Type:        "Field",
						Package:     pkg.PackageName,
						Description: truncate(stripHTML(field.Description), 120),
						Anchor:      "struct-" + st.Name,
					})
				}
			}
			// Search struct methods
			for _, m := range st.Methods {
				if deepMatchFunc(m, q) {
					results = append(results, SearchResult{
						Name:        st.Name + "." + m.Name,
						Type:        "Method",
						Package:     pkg.PackageName,
						Description: truncate(stripHTML(m.WhatItDoes), 120),
						Anchor:      "method-" + st.Name + "-" + m.Name,
					})
				}
			}
		}

		// Search interfaces
		for _, iface := range pkg.Interfaces {
			if deepMatchInterface(iface, q) {
				results = append(results, SearchResult{
					Name:        iface.Name,
					Type:        "Interface",
					Package:     pkg.PackageName,
					Description: truncate(stripHTML(iface.WhatItDoes), 120),
					Anchor:      "interface-" + iface.Name,
				})
			}
			// Search interface methods
			for _, m := range iface.Methods {
				if deepMatchFunc(m, q) {
					results = append(results, SearchResult{
						Name:        iface.Name + "." + m.Name,
						Type:        "Method",
						Package:     pkg.PackageName,
						Description: truncate(stripHTML(m.WhatItDoes), 120),
						Anchor:      "interface-" + iface.Name,
					})
				}
			}
		}

		// Search types
		for _, t := range pkg.Types {
			if containsAny(q, t.Name, t.Overview, t.WhatItDoes, t.WhenToUse, t.InSimpleWords, t.Underlying, t.Example, t.Explanation) {
				results = append(results, SearchResult{
					Name:        t.Name,
					Type:        "Type",
					Package:     pkg.PackageName,
					Description: truncate(stripHTML(t.WhatItDoes), 120),
					Anchor:      "type-" + t.Name,
				})
			}
		}

		// Search constants
		for _, c := range pkg.Constants {
			if containsAny(q, c.Name, c.Value, c.Type, c.Description) {
				results = append(results, SearchResult{
					Name:        c.Name,
					Type:        "Constant",
					Package:     pkg.PackageName,
					Description: truncate(c.Description, 120),
					Anchor:      "constants",
				})
			}
		}

		// Search variables
		for _, v := range pkg.Variables {
			if containsAny(q, v.Name, v.Type, v.Description) {
				results = append(results, SearchResult{
					Name:        v.Name,
					Type:        "Variable",
					Package:     pkg.PackageName,
					Description: truncate(v.Description, 120),
					Anchor:      "variables",
				})
			}
		}

		// Search package-level text (overview, what it does, when to use)
		if containsAny(q, pkg.PackageName, pkg.ImportPath, pkg.Overview, pkg.WhatItDoes, pkg.WhenToUse) {
			// Only add if not already found via a child item
			found := false
			for _, r := range results {
				if r.Package == pkg.PackageName && r.Name == pkg.PackageName {
					found = true
					break
				}
			}
			if !found {
				results = append(results, SearchResult{
					Name:        pkg.PackageName,
					Type:        "Package",
					Package:     pkg.PackageName,
					Description: truncate(stripHTML(pkg.Overview), 120),
					Anchor:      "overview",
				})
			}
		}
	}

	// Deduplicate: if same Name+Package+Type exists, keep only the first
	seen := map[string]bool{}
	var deduped []SearchResult
	for _, r := range results {
		key := r.Name + "|" + r.Package + "|" + r.Type
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, r)
		}
	}

	return deduped
}

// deepMatchFunc searches all text in a function: name, doc, params, returns, signature, example, explanation
func deepMatchFunc(f docs.FriendlyFunction, q string) bool {
	if containsAny(q, f.Name, f.Overview, f.WhatItDoes, f.WhenToUse, f.Signature, f.Example, f.Explanation, f.Receiver) {
		return true
	}
	for _, p := range f.Params {
		if containsAny(q, p.Name, p.Type, p.Description) {
			return true
		}
	}
	for _, r := range f.Returns {
		if strings.Contains(strings.ToLower(r), q) {
			return true
		}
	}
	return false
}

// deepMatchStruct searches all text in a struct: name, doc, fields, methods, example, explanation
func deepMatchStruct(st docs.FriendlyStruct, q string) bool {
	if containsAny(q, st.Name, st.Overview, st.WhatItDoes, st.WhenToUse, st.Example, st.Explanation) {
		return true
	}
	for _, f := range st.Fields {
		if containsAny(q, f.Name, f.Type, f.Description, f.Tag) {
			return true
		}
	}
	for _, m := range st.Methods {
		if deepMatchFunc(m, q) {
			return true
		}
	}
	return false
}

// deepMatchInterface searches all text in an interface
func deepMatchInterface(iface docs.FriendlyInterface, q string) bool {
	if containsAny(q, iface.Name, iface.Overview, iface.WhatItDoes, iface.WhenToUse, iface.Example, iface.Explanation) {
		return true
	}
	for _, m := range iface.Methods {
		if deepMatchFunc(m, q) {
			return true
		}
	}
	return false
}

// containsAny checks if query appears in any of the given strings (case-insensitive)
func containsAny(query string, texts ...string) bool {
	for _, t := range texts {
		if t != "" && strings.Contains(strings.ToLower(t), query) {
			return true
		}
	}
	return false
}

// stripHTML removes HTML tags for clean display in search results
func stripHTML(s string) string {
	var out strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			out.WriteRune(r)
		}
	}
	return strings.TrimSpace(out.String())
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// DiagramNode represents a node in the package structure diagram
type DiagramNode struct {
	ID       string         `json:"id"`
	Label    string         `json:"label"`
	Type     string         `json:"type"` // "package", "struct", "interface", "function", "type"
	Children []DiagramNode  `json:"children,omitempty"`
	Fields   []DiagramField `json:"fields,omitempty"`
	Methods  []string       `json:"methods,omitempty"`
	Params   string         `json:"params,omitempty"`
	Returns  string         `json:"returns,omitempty"`
	Count    map[string]int `json:"count,omitempty"`
}

// DiagramField is a simplified field for diagram display
type DiagramField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// handleDiagram returns JSON data for rendering a package structure diagram
func (s *Server) handleDiagram(w http.ResponseWriter, r *http.Request) {
	pkgName := strings.TrimPrefix(r.URL.Path, "/api/diagram/")
	pkgName = strings.TrimSuffix(pkgName, "/")

	var pkg *docs.FriendlyDoc
	for _, p := range s.docs {
		if p.PackageName == pkgName || p.ImportPath == pkgName {
			pkg = p
			break
		}
	}

	if pkg == nil && len(s.docs) > 0 {
		// Build an overview diagram showing ALL packages
		root := DiagramNode{
			ID:    "root",
			Label: s.packagePath,
			Type:  "package",
			Count: map[string]int{"packages": len(s.docs)},
		}
		for _, p := range s.docs {
			child := DiagramNode{
				ID:    "pkg-" + p.PackageName,
				Label: p.PackageName,
				Type:  "package",
				Count: map[string]int{
					"structs":    len(p.Structs),
					"interfaces": len(p.Interfaces),
					"functions":  len(p.Functions),
					"types":      len(p.Types),
				},
			}
			root.Children = append(root.Children, child)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(root)
		return
	}

	if pkg == nil {
		http.NotFound(w, r)
		return
	}

	// Build detailed diagram for a single package
	root := DiagramNode{
		ID:    "pkg-" + pkg.PackageName,
		Label: pkg.PackageName,
		Type:  "package",
		Count: map[string]int{
			"structs":    len(pkg.Structs),
			"interfaces": len(pkg.Interfaces),
			"functions":  len(pkg.Functions),
			"types":      len(pkg.Types),
			"constants":  len(pkg.Constants),
			"variables":  len(pkg.Variables),
		},
	}

	// Structs
	for _, st := range pkg.Structs {
		node := DiagramNode{
			ID:    "struct-" + st.Name,
			Label: st.Name,
			Type:  "struct",
		}
		for _, f := range st.Fields {
			node.Fields = append(node.Fields, DiagramField{Name: f.Name, Type: f.Type})
		}
		for _, m := range st.Methods {
			node.Methods = append(node.Methods, m.Name)
		}
		root.Children = append(root.Children, node)
	}

	// Interfaces
	for _, iface := range pkg.Interfaces {
		node := DiagramNode{
			ID:    "interface-" + iface.Name,
			Label: iface.Name,
			Type:  "interface",
		}
		for _, m := range iface.Methods {
			node.Methods = append(node.Methods, m.Name)
		}
		root.Children = append(root.Children, node)
	}

	// Functions
	for _, fn := range pkg.Functions {
		node := DiagramNode{
			ID:      "func-" + fn.Name,
			Label:   fn.Name,
			Type:    "function",
			Returns: strings.Join(fn.Returns, ", "),
		}
		params := []string{}
		for _, p := range fn.Params {
			if p.Name != "" {
				params = append(params, p.Name)
			} else {
				params = append(params, p.Type)
			}
		}
		node.Params = strings.Join(params, ", ")
		root.Children = append(root.Children, node)
	}

	// Types
	for _, t := range pkg.Types {
		node := DiagramNode{
			ID:    "type-" + t.Name,
			Label: t.Name,
			Type:  "type",
		}
		root.Children = append(root.Children, node)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(root)
}
