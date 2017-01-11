package maps

import (
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func TestParseSingleLine(t *testing.T) {
	cases := []struct {
		name     string
		raw      string
		expected TextMapTree
	}{
		{
			"Only one root line",
			"1. Linelineline",
			TextMapTree{
				NodesCollection{
					MapNode{
						"first_level",
						1,
						"Linelineline",
						NodesCollection{},
					},
				},
			},
		},
		{
			"Two root lines",
			`
1. Linelineline
2. Line2Line2Line
`,
			TextMapTree{
				NodesCollection{
					MapNode{
						"first_level",
						1,
						"Linelineline",
						NodesCollection{},
					},
					MapNode{
						"first_level",
						1,
						"Line2Line2Line",
						NodesCollection{},
					},
				},
			},
		},
		{
			"OneRootWithTwoEmbed",
			`
1. Linelineline
    * Line 1
    * Line 2
`,
			TextMapTree{
				NodesCollection{
					MapNode{
						"first_level",
						1,
						"Linelineline",
						NodesCollection{
							{
								"embed",
								2,
								"Line 1",
								NodesCollection{},
							},
							{
								"embed",
								2,
								"Line 2",
								NodesCollection{},
							},
						},
					},
				},
			},
		},
		{
			"TwoRootsWithTwoEmbed",
			`
1. Linelineline
    * Line 1
    * Line 2
2. Linelineline2
    * Line 3
    * Line 4
`,
			TextMapTree{
				NodesCollection{
					MapNode{
						"first_level",
						1,
						"Linelineline",
						NodesCollection{
							{
								"embed",
								2,
								"Line 1",
								NodesCollection{},
							},
							{
								"embed",
								2,
								"Line 2",
								NodesCollection{},
							},
						},
					},
					MapNode{
						"first_level",
						1,
						"Linelineline2",
						NodesCollection{
							{
								"embed",
								2,
								"Line 3",
								NodesCollection{},
							},
							{
								"embed",
								2,
								"Line 4",
								NodesCollection{},
							},
						},
					},
				},
			},
		},
		{
			"RecursiveEmbed",
			`
1. Linelineline
    * Line 1
    	* Line 1 - 1
    	* Line 1 - 2

`,
			TextMapTree{
				NodesCollection{
					MapNode{
						"first_level",
						1,
						"Linelineline",
						NodesCollection{
							{
								"embed",
								2,
								"Line 1",
								NodesCollection{
									{
										"embed",
										3,
										"Line 1 - 1",
										NodesCollection{},
									},
									{"embed",
										3,
										"Line 1 - 2",
										NodesCollection{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"ComplexTreeWithoutDescriptions",
			`
1. Linelineline
    * Line 1
    	* Line 1 - 1
    * Line 2
    	* Line 2 - 1
    	    * Line 2 - 2 - 1
2. Linelineline2
    * Line 3
    	* Line 3 - 1
    	* Line 3 - 2
    * Line 4
    	* Line 4 - 1
    	* Line 4 - 2
`,
			TextMapTree{
				NodesCollection{
					MapNode{
						"first_level",
						1,
						"Linelineline",
						NodesCollection{
							{
								"embed",
								2,
								"Line 1",
								NodesCollection{
									MapNode{
										"embed",
										3,
										"Line 1 - 1",
										NodesCollection{},
									},
								},
							},
							{
								"embed",
								2,
								"Line 2",
								NodesCollection{
									MapNode{
										"embed",
										3,
										"Line 2 - 1",
										NodesCollection{
											MapNode{
												"embed",
												4,
												"Line 2 - 2 - 1",
												NodesCollection{},
											},
										},
									},
								},
							},
						},
					},
					MapNode{
						"first_level",
						1,
						"Linelineline2",
						NodesCollection{
							{
								"embed",
								2,
								"Line 3",
								NodesCollection{
									MapNode{
										"embed",
										3,
										"Line 3 - 1",
										NodesCollection{},
									},
									MapNode{
										"embed",
										3,
										"Line 3 - 2",
										NodesCollection{},
									},
								},
							},
							{
								"embed",
								2,
								"Line 4",
								NodesCollection{
									MapNode{
										"embed",
										3,
										"Line 4 - 1",
										NodesCollection{},
									},
									MapNode{
										"embed",
										3,
										"Line 4 - 2",
										NodesCollection{},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			res, err := parserFactory{}.Get(tcase.raw).Parse()
			if err != nil {
				t.Fatalf("Have error %#v", err)
			}
			if !reflect.DeepEqual(res, tcase.expected) {
				t.Fatalf("Expected %v, have %v, diff %v", tcase.expected, res, pretty.Diff(tcase.expected, res))
			}
		})
	}
}
