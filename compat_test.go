package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func testTypes(
	t *testing.T,
	expected *Symbol,
	actual *Symbol) {

	if expected == nil && actual == nil {
		return
	}

	if expected.Name != actual.Name {
		t.Errorf("Type name is mistaken.\n"+
			"\tExpected: %v\n"+
			"\tActual: %v\n", expected.Name, actual.Name)
	}

	for index, _ := range expected.Types {
		testTypes(t, expected.Types[index], actual.Types[index])
	}
}

func testCompat(
	t *testing.T,
	source string,
	expected CompatContext) {

	fileSet := token.NewFileSet()
	file, _ := parser.ParseFile(fileSet, "source.go", source, parser.ParseComments)

	actual := &CompatContext{Packages: map[string]*Package{}}
	ProcessFile(fileSet, file, actual)

	for packageName, expectedPackage := range expected.Packages {
		if actualPackage, ok := actual.Packages[packageName]; ok {
			if expectedPackage.Name != actualPackage.Name {
				t.Errorf("Package name is mistaken.\n"+
					"\tExpected: %v\n"+
					"\tActual: %v\n", expectedPackage.Name, actualPackage.Name)
			}

			for index, expectedSymbol := range expectedPackage.Exported {
				actualSymbol := actualPackage.Exported[index]

				if expectedSymbol.Name != actualSymbol.Name {
					t.Errorf("Symbol name is mistaken.\n"+
						"\tExpected: %v\n"+
						"\tActual: %v\n", expectedSymbol.Name, actualSymbol.Name)
				}

				if len(expectedSymbol.Types) != len(actualSymbol.Types) {
					t.Errorf("Wrong number of types for %v.\n"+
						"\tExpected: %v\n"+
						"\tActual: %v\n", expectedSymbol.Name, len(expectedSymbol.Types), len(actualSymbol.Types))
				}

				for index, expectedType := range expectedSymbol.Types {
					actualType := actualSymbol.Types[index]
					testTypes(t, expectedType, actualType)
				}
			}
		} else {
			t.Errorf("Package %v was expected but not found in the context.", packageName)
		}
	}
}

func TestSimpleType(t *testing.T) {
	source := `
package p

type MyInt int
`

	expected := CompatContext{
		Packages: map[string]*Package{
			"p": &Package{
				Name: "p",
				Exported: []*Symbol{
					&Symbol{
						Name: "MyInt",
						Types: []*Symbol{
							&Symbol{Name: "int"},
						},
					},
				},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestStructType(t *testing.T) {
	source := `
package p

type MyInt struct {
	A	int
	B	float32
	C	string
}
`

	expected := CompatContext{
		Packages: map[string]*Package{
			"p": &Package{
				Name: "p",
				Exported: []*Symbol{
					&Symbol{
						Name: "MyInt",
						Types: []*Symbol{
							&Symbol{"A", []*Symbol{&Symbol{Name: "int"}}},
							&Symbol{"B", []*Symbol{&Symbol{Name: "float32"}}},
							&Symbol{"C", []*Symbol{&Symbol{Name: "string"}}},
						},
					},
				},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestNestedStructType(t *testing.T) {
	source := `
package p

type MyInt struct {
	A	int
	B	struct {
		C	float32
		D	string
	}
}
`

	expected := CompatContext{
		Packages: map[string]*Package{
			"p": &Package{
				Name: "p",
				Exported: []*Symbol{
					&Symbol{
						Name: "MyInt",
						Types: []*Symbol{
							&Symbol{"A", []*Symbol{&Symbol{Name: "int"}}},
							&Symbol{"B", []*Symbol{
								&Symbol{"C", []*Symbol{&Symbol{Name: "float32"}}},
								&Symbol{"D", []*Symbol{&Symbol{Name: "string"}}},
							}},
						},
					},
				},
			},
		},
	}

	testCompat(t, source, expected)
}