package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type ModelDefinition struct {
	TableName  string
	FieldTypes map[string]string
}

// ParseModelDefinitionsFromFile 解析给定 Go 文件中的 modelDefinitions 信息
func ParseModelDefinitionsFromFile(filename string) ([]ModelDefinition, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	var result []ModelDefinition

	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if name.Name != "modelDefinitions" {
					continue
				}

				compLit, ok := valueSpec.Values[i].(*ast.CompositeLit)
				if !ok {
					continue
				}

				for _, elt := range compLit.Elts {
					modelLit, ok := elt.(*ast.CompositeLit)
					if !ok {
						continue
					}

					var info ModelDefinition
					info.FieldTypes = make(map[string]string)

					for _, field := range modelLit.Elts {
						kv, ok := field.(*ast.KeyValueExpr)
						if !ok {
							continue
						}

						key := getIdentName(kv.Key)

						switch key {
						case "TableName":
							if val, ok := kv.Value.(*ast.BasicLit); ok {
								info.TableName = trimQuotes(val.Value)
							}
						case "FieldTypes":
							if fieldMap, ok := kv.Value.(*ast.CompositeLit); ok {
								for _, m := range fieldMap.Elts {
									mapKV, ok := m.(*ast.KeyValueExpr)
									if !ok {
										continue
									}
									mapKey := trimQuotes(getBasicLitValue(mapKV.Key))
									mapVal := trimQuotes(getBasicLitValue(mapKV.Value))
									info.FieldTypes[mapKey] = mapVal
								}
							}
						}
					}

					result = append(result, info)
				}
			}
		}
		return true
	})

	return result, nil
}

// 获取标识符名称
func getIdentName(expr ast.Expr) string {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// 获取 BasicLit 字面值
func getBasicLitValue(expr ast.Expr) string {
	if bl, ok := expr.(*ast.BasicLit); ok {
		return bl.Value
	}
	return ""
}

// 去掉字符串两边引号
func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
