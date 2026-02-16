// DDD module code generator
//
// Usage:
//
//	go run cmd/gen/main.go -name order -cn Order
//	make gen name=order cn=Order
//
// Output:
//
//	internal/domain/order/entity.go          - Domain entity
//	internal/domain/order/repository.go      - Repository interface
//	internal/infrastructure/persistence/database/order_model.go  - Data model
//	internal/infrastructure/persistence/database/order_repo.go   - Repository impl
//	internal/application/dto/order_dto.go    - DTO
//	internal/application/service/order_service.go - Application service
//	internal/interfaces/http/handler/order_handler.go - HTTP handler
//	Also auto-registers in router.go and container.go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type ModuleData struct {
	Name        string // order
	PascalName  string // Order
	CamelName   string // order
	SnakeName   string // order
	KebabName   string // order
	PluralName  string // orders
	ChineseName string // Order
	ModulePath  string // go-ddd-scaffold
}

func main() {
	name := flag.String("name", "", "module name (lowercase, e.g. order)")
	cn := flag.String("cn", "", "display name (e.g. Order)")
	modulePath := flag.String("module", "", "Go module path (auto-detected)")
	flag.Parse()

	if *name == "" {
		fmt.Fprintln(os.Stderr, "error: -name is required")
		fmt.Fprintln(os.Stderr, "usage: go run cmd/gen/main.go -name order -cn Order")
		os.Exit(1)
	}

	if *cn == "" {
		*cn = *name
	}

	if *modulePath == "" {
		*modulePath = detectModulePath()
	}

	data := ModuleData{
		Name:        strings.ToLower(*name),
		PascalName:  toPascalCase(*name),
		CamelName:   toCamelCase(*name),
		SnakeName:   toSnakeCase(*name),
		KebabName:   toKebabCase(*name),
		PluralName:  toPlural(strings.ToLower(*name)),
		ChineseName: *cn,
		ModulePath:  *modulePath,
	}

	fmt.Printf("Generating DDD module: %s (%s)\n", data.PascalName, data.ChineseName)

	// Generate files
	files := []struct {
		tmpl string
		out  string
	}{
		{"templates/domain_entity.go.tmpl", fmt.Sprintf("internal/domain/%s/entity.go", data.SnakeName)},
		{"templates/domain_repository.go.tmpl", fmt.Sprintf("internal/domain/%s/repository.go", data.SnakeName)},
		{"templates/infra_model.go.tmpl", fmt.Sprintf("internal/infrastructure/persistence/database/%s_model.go", data.SnakeName)},
		{"templates/infra_repo.go.tmpl", fmt.Sprintf("internal/infrastructure/persistence/database/%s_repo.go", data.SnakeName)},
		{"templates/app_dto.go.tmpl", fmt.Sprintf("internal/application/dto/%s_dto.go", data.SnakeName)},
		{"templates/app_service.go.tmpl", fmt.Sprintf("internal/application/service/%s_service.go", data.SnakeName)},
		{"templates/handler.go.tmpl", fmt.Sprintf("internal/interfaces/http/handler/%s_handler.go", data.SnakeName)},
	}

	for _, f := range files {
		if err := generateFile(f.tmpl, f.out, data); err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate %s: %v\n", f.out, err)
			os.Exit(1)
		}
		fmt.Printf("  + %s\n", f.out)
	}

	// Auto-register route
	if err := appendRoute(data); err != nil {
		fmt.Fprintf(os.Stderr, "  ! failed to register route: %v (please add manually)\n", err)
	} else {
		fmt.Println("  + route registered in router.go")
	}

	// Auto-register container
	if err := appendContainer(data); err != nil {
		fmt.Fprintf(os.Stderr, "  ! failed to register container: %v (please add manually)\n", err)
	} else {
		fmt.Println("  + service registered in container.go")
	}

	// Auto-register migration
	if err := appendMigration(data); err != nil {
		fmt.Fprintf(os.Stderr, "  ! failed to register migration: %v (please add manually)\n", err)
	} else {
		fmt.Println("  + model migration registered in container.go")
	}

	fmt.Printf("\nModule %s generated!\n", data.PascalName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  1. Edit internal/domain/%s/entity.go - add domain fields and business methods\n", data.SnakeName)
	fmt.Printf("  2. Edit internal/application/service/%s_service.go - implement business orchestration\n", data.SnakeName)
	fmt.Printf("  3. Run make docs - update Swagger documentation\n")
}

func generateFile(tmplPath, outPath string, data ModuleData) error {
	if _, err := os.Stat(outPath); err == nil {
		return fmt.Errorf("file already exists: %s", outPath)
	}

	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	t, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}

func appendRoute(data ModuleData) error {
	routerFile := "internal/interfaces/http/router/router.go"
	content, err := os.ReadFile(routerFile)
	if err != nil {
		return err
	}

	marker := "// GEN:ROUTE_REGISTER - Code generator appends routes here, do not remove"
	routeCode := fmt.Sprintf(`// %s module
			%sHandler := handler.New%sHandler(c.%sService)
			%s := authorized.Group("/%s")
			{
				%s.GET("", %sHandler.List)
				%s.POST("", %sHandler.Create)
				%s.GET("/:id", %sHandler.Get)
				%s.PUT("/:id", %sHandler.Update)
				%s.DELETE("/:id", %sHandler.Delete)
			}

			`,
		data.PascalName,
		data.CamelName, data.PascalName, data.PascalName,
		data.PluralName, data.PluralName,
		data.PluralName, data.CamelName,
		data.PluralName, data.CamelName,
		data.PluralName, data.CamelName,
		data.PluralName, data.CamelName,
		data.PluralName, data.CamelName,
	)

	newContent := strings.Replace(string(content), marker, routeCode+marker, 1)
	if newContent == string(content) {
		return fmt.Errorf("route marker comment not found")
	}

	return os.WriteFile(routerFile, []byte(newContent), 0o644)
}

func appendContainer(data ModuleData) error {
	containerFile := "internal/container/container.go"
	content, err := os.ReadFile(containerFile)
	if err != nil {
		return err
	}

	// Append service field
	fieldMarker := "// GEN:SERVICE_REGISTER - Code generator appends services here, do not remove"
	fieldCode := fmt.Sprintf("%sService *service.%sAppService\n\t", data.PascalName, data.PascalName)
	newContent := strings.Replace(string(content), fieldMarker, fieldCode+fieldMarker, 1)

	// Append initialization
	initMarker := "// GEN:SERVICE_INIT - Code generator appends initialization here, do not remove"
	initCode := fmt.Sprintf(`%sRepo := database.New%sRepository(db)
	c.%sService = service.New%sAppService(%sRepo)
	`,
		data.CamelName, data.PascalName,
		data.PascalName, data.PascalName, data.CamelName,
	)
	newContent = strings.Replace(newContent, initMarker, initCode+initMarker, 1)

	if newContent == string(content) {
		return fmt.Errorf("container marker comment not found")
	}

	// Ensure domain package is imported
	domainImport := fmt.Sprintf(`"%s/internal/domain/%s"`, data.ModulePath, data.SnakeName)
	_ = domainImport // domain package is implicitly referenced via repository interface

	return os.WriteFile(containerFile, []byte(newContent), 0o644)
}

func appendMigration(data ModuleData) error {
	containerFile := "internal/container/container.go"
	content, err := os.ReadFile(containerFile)
	if err != nil {
		return err
	}

	marker := "// GEN:MODEL_MIGRATE - Code generator appends models here, do not remove"
	migrationCode := fmt.Sprintf("&database.%sModel{},\n\t\t", data.PascalName)

	newContent := strings.Replace(string(content), marker, migrationCode+marker, 1)
	if newContent == string(content) {
		return fmt.Errorf("migration marker comment not found")
	}

	return os.WriteFile(containerFile, []byte(newContent), 0o644)
}

func detectModulePath() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "go-ddd-scaffold"
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return "go-ddd-scaffold"
}

// ========================
// Naming utilities
// ========================

func toPascalCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		if len(p) > 0 {
			runes := []rune(strings.ToLower(p))
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, "")
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if pascal == "" {
		return ""
	}
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func toSnakeCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	return strings.Join(parts, "_")
}

func toKebabCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	return strings.Join(parts, "-")
}

func toPlural(s string) string {
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") && len(s) > 1 {
		prev := s[len(s)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return s[:len(s)-1] + "ies"
		}
	}
	return s + "s"
}

func splitWords(s string) []string {
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	var result []string
	for _, p := range parts {
		var current []rune
		for i, r := range p {
			if unicode.IsUpper(r) && i > 0 {
				result = append(result, string(current))
				current = nil
			}
			current = append(current, r)
		}
		if len(current) > 0 {
			result = append(result, string(current))
		}
	}
	return result
}
