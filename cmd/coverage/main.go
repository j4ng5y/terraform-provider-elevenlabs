package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type specOperation struct {
	Method  string
	Path    string
	Summary string
	Tags    []string
}

type clientOperation struct {
	Method string
	Path   string
	Func   string
}

type key struct {
	Method string
	Path   string
}

type filter struct {
	includeTags   map[string]struct{}
	excludeTags   map[string]struct{}
	includeMethod map[string]struct{}
}

var bracePattern = regexp.MustCompile(`\{[^/}]+\}`)

func main() {
	var specPath, clientPath, includeTagsRaw, excludeTagsRaw, includeMethodsRaw string
	var showDetails bool

	flag.StringVar(&specPath, "spec", "openapi_pretty.json", "Path to the ElevenLabs OpenAPI spec")
	flag.StringVar(&clientPath, "client", filepath.Join("internal", "client", "client.go"), "Path to the ElevenLabs client implementation")
	flag.StringVar(&includeTagsRaw, "include-tags", "", "Comma-separated list of tags to include (defaults to all)")
	flag.StringVar(&excludeTagsRaw, "exclude-tags", "", "Comma-separated list of tags to exclude")
	flag.StringVar(&includeMethodsRaw, "methods", "", "Comma-separated list of HTTP methods to include (e.g., GET,POST)")
	flag.BoolVar(&showDetails, "details", true, "Show detailed missing operation list")
	flag.Parse()

	filt := newFilter(includeTagsRaw, excludeTagsRaw, includeMethodsRaw)

	specOps, err := loadSpecOperations(specPath, filt)
	if err != nil {
		fatalf("load spec: %v", err)
	}

	clientOps, err := parseClientOperations(clientPath)
	if err != nil {
		fatalf("parse client: %v", err)
	}

	reportCoverage(specOps, clientOps, showDetails)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func newFilter(includeTagsRaw, excludeTagsRaw, includeMethodsRaw string) filter {
	return filter{
		includeTags:   buildSet(includeTagsRaw),
		excludeTags:   buildSet(excludeTagsRaw),
		includeMethod: buildSet(strings.ToUpper(includeMethodsRaw)),
	}
}

func buildSet(csv string) map[string]struct{} {
	if strings.TrimSpace(csv) == "" {
		return nil
	}

	items := strings.Split(csv, ",")
	set := make(map[string]struct{}, len(items))
	for _, raw := range items {
		item := strings.TrimSpace(raw)
		if item == "" {
			continue
		}
		set[strings.ToLower(item)] = struct{}{}
	}
	return set
}

func (f filter) allow(tags []string, method string) bool {
	method = strings.ToLower(strings.TrimSpace(method))
	if len(f.includeMethod) > 0 {
		if _, ok := f.includeMethod[method]; !ok {
			return false
		}
	}

	normalizedTags := tags
	if len(normalizedTags) == 0 {
		normalizedTags = []string{"untagged"}
	}

	seenIncluded := len(f.includeTags) == 0

	for _, tag := range normalizedTags {
		lt := strings.ToLower(tag)
		if _, excluded := f.excludeTags[lt]; excluded {
			return false
		}
		if !seenIncluded {
			if _, ok := f.includeTags[lt]; ok {
				seenIncluded = true
			}
		}
	}

	return seenIncluded
}

func loadSpecOperations(path string, filt filter) ([]specOperation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var spec struct {
		Paths map[string]map[string]struct {
			Summary string   `json:"summary"`
			Tags    []string `json:"tags"`
		} `json:"paths"`
	}

	if err := json.NewDecoder(file).Decode(&spec); err != nil {
		return nil, err
	}

	var ops []specOperation
	for rawPath, methods := range spec.Paths {
		for method, op := range methods {
			if strings.HasPrefix(method, "x-") {
				continue
			}
			upperMethod := strings.ToUpper(method)
			if !filt.allow(op.Tags, upperMethod) {
				continue
			}
			ops = append(ops, specOperation{
				Method:  upperMethod,
				Path:    canonicalizePath(rawPath),
				Summary: op.Summary,
				Tags:    op.Tags,
			})
		}
	}

	return ops, nil
}

func parseClientOperations(path string) ([]clientOperation, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	consts := extractConstStrings(file)

	var ops []clientOperation
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkg, ok := sel.X.(*ast.Ident)
			if !ok || pkg.Name != "http" || sel.Sel.Name != "NewRequest" {
				return true
			}

			if len(call.Args) < 2 {
				return true
			}

			method := methodFromExpr(call.Args[0])
			if method == "" {
				return true
			}

			rawPath := exprToString(call.Args[1], consts)
			if rawPath == "" {
				return true
			}

			ops = append(ops, clientOperation{
				Method: method,
				Path:   canonicalizePath(rawPath),
				Func:   fn.Name.Name,
			})

			return true
		})
	}

	return ops, nil
}

func extractConstStrings(file *ast.File) map[string]string {
	consts := map[string]string{}

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}

		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Values) != 1 {
				continue
			}

			literal, ok := vs.Values[0].(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				continue
			}

			val, err := strconv.Unquote(literal.Value)
			if err != nil {
				continue
			}

			for _, name := range vs.Names {
				consts[name.Name] = val
			}
		}
	}

	return consts
}

func methodFromExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			val, err := strconv.Unquote(e.Value)
			if err == nil {
				return strings.ToUpper(val)
			}
		}
	case *ast.SelectorExpr:
		if ident, ok := e.X.(*ast.Ident); ok && ident.Name == "http" && strings.HasPrefix(e.Sel.Name, "Method") {
			return strings.ToUpper(strings.TrimPrefix(e.Sel.Name, "Method"))
		}
	}
	return ""
}

func exprToString(expr ast.Expr, consts map[string]string) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			val, err := strconv.Unquote(e.Value)
			if err == nil {
				return val
			}
			return e.Value
		}
	case *ast.Ident:
		if val, ok := consts[e.Name]; ok {
			return val
		}
		return fmt.Sprintf("{%s}", e.Name)
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return exprToString(e.X, consts) + exprToString(e.Y, consts)
		}
	case *ast.CallExpr:
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "fmt" && sel.Sel.Name == "Sprintf" {
				if len(e.Args) == 0 {
					return ""
				}
				format := exprToString(e.Args[0], consts)
				args := make([]interface{}, 0, len(e.Args)-1)
				for _, argExpr := range e.Args[1:] {
					val := exprToString(argExpr, consts)
					if val == "" {
						val = "*"
					}
					args = append(args, val)
				}
				return fmt.Sprintf(format, args...)
			}
		}
	case *ast.ParenExpr:
		return exprToString(e.X, consts)
	}
	return ""
}

func canonicalizePath(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}

	if idx := strings.Index(raw, "://"); idx >= 0 {
		raw = raw[idx+3:]
		if slash := strings.Index(raw, "/"); slash >= 0 {
			raw = raw[slash:]
		} else {
			raw = "/"
		}
	}

	raw = strings.TrimPrefix(raw, "api.elevenlabs.io")

	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}

	if idx := strings.Index(raw, "?"); idx >= 0 {
		raw = raw[:idx]
	}

	raw = strings.ReplaceAll(raw, "//", "/")

	if len(raw) > 1 {
		raw = strings.TrimSuffix(raw, "/")
	}

	raw = bracePattern.ReplaceAllString(raw, "*")
	return raw
}

func reportCoverage(specOps []specOperation, clientOps []clientOperation, showDetails bool) {
	specKeys := make(map[key][]specOperation)
	for _, op := range specOps {
		k := key{Method: op.Method, Path: op.Path}
		specKeys[k] = append(specKeys[k], op)
	}

	clientKeys := make(map[key][]clientOperation)
	for _, op := range clientOps {
		k := key{Method: op.Method, Path: op.Path}
		clientKeys[k] = append(clientKeys[k], op)
	}

	totalSpec := len(specKeys)
	covered := 0
	for k := range specKeys {
		if _, ok := clientKeys[k]; ok {
			covered++
		}
	}

	fmt.Printf("Spec operations considered: %d\n", totalSpec)
	fmt.Printf("Client operations detected: %d\n", len(clientKeys))
	coverage := 0.0
	if totalSpec > 0 {
		coverage = float64(covered) / float64(totalSpec) * 100
	}
	fmt.Printf("Coverage: %d/%d (%.1f%%)\n", covered, totalSpec, coverage)

	var missing []specOperation
	for k, ops := range specKeys {
		if _, ok := clientKeys[k]; ok {
			continue
		}
		missing = append(missing, ops...)
	}

	if len(missing) == 0 {
		fmt.Println("All operations covered! 🥳")
		return
	}

	sort.Slice(missing, func(i, j int) bool {
		if missing[i].Method == missing[j].Method {
			return missing[i].Path < missing[j].Path
		}
		return missing[i].Method < missing[j].Method
	})

	tagCounts := map[string]int{}
	for _, op := range missing {
		if len(op.Tags) == 0 {
			tagCounts["untagged"]++
			continue
		}
		for _, tag := range op.Tags {
			tagCounts[tag]++
		}
	}

	type tagCount struct {
		Tag   string
		Count int
	}

	var sortedTags []tagCount
	for tag, count := range tagCounts {
		sortedTags = append(sortedTags, tagCount{Tag: tag, Count: count})
	}

	sort.Slice(sortedTags, func(i, j int) bool {
		if sortedTags[i].Count == sortedTags[j].Count {
			return sortedTags[i].Tag < sortedTags[j].Tag
		}
		return sortedTags[i].Count > sortedTags[j].Count
	})

	fmt.Println("Missing operations by tag:")
	for _, tc := range sortedTags {
		fmt.Printf("  %-30s %d\n", tc.Tag, tc.Count)
	}

	if !showDetails {
		return
	}

	fmt.Println("\nMissing operations:")
	for _, op := range missing {
		tag := "untagged"
		if len(op.Tags) > 0 {
			tag = op.Tags[0]
		}
		fmt.Printf("  %-6s %-50s (%s) %s\n", op.Method, op.Path, tag, op.Summary)
	}
}
