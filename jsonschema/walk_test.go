package jsonschema

import (
	"encoding/json"
	"reflect"
	"testing"
)

type testWalker struct {
	calls    map[string]Instance
	rawCalls map[string]json.RawMessage
}

func (tw *testWalker) walkFn(path string, i Instance, value json.RawMessage) error {
	tw.calls[path] = i
	tw.rawCalls[path] = value
	return nil
}

func TestWalkJSONSchema(t *testing.T) {
	tests := []struct {
		description string
		oneOfType   string
		schemaPath  string
		want        map[string]Instance
		wantErr     bool
	}{
		{
			description: "Basic walk, no allOf, no oneOf",
			schemaPath:  "./test_data/image.json",
			want: map[string]Instance{
				"$.type": {Type: "string"},
				"$.crops": {Type: "array", Items: []byte(`{
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
			}`)},
				"$.crops[*]": {Type: "object", Properties: map[string]json.RawMessage{
					"name":         []byte(`{"type": "string", "default": "name"}`),
					"width":        []byte(`{"type": "number" }`),
					"height":       []byte(`{"type": "number" }`),
					"path":         []byte(`{"type": "string" }`),
					"relativePath": []byte(`{"type": "string" }`),
				},
					Required: []string{"name", "width", "height", "path", "relativePath"},
				},
				"$.crops[*].name":         {Type: "string"},
				"$.crops[*].width":        {Type: "number"},
				"$.crops[*].height":       {Type: "number"},
				"$.crops[*].path":         {Type: "string"},
				"$.crops[*].relativePath": {Type: "string"},
				"$.URL": {Type: "object", Properties: map[string]json.RawMessage{
					"publish": []byte(`{"type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }}`),
					"absolute": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }}`),
				},
					Required: []string{"publish", "absolute"},
				},
				"$.URL.publish":  {Type: "string"},
				"$.URL.absolute": {Type: "string"},
			},
		},
		{
			description: "Walk with allOf, no oneOf",
			schemaPath:  "./test_data/embed_parent.json",
			want:        map[string]Instance{"$.type": {Type: "string"}},
		},
		{
			description: "Walk with oneOf, no allOf",
			oneOfType:   "image",
			schemaPath:  "./test_data/image_parent.json",
			want: map[string]Instance{
				"$.type": {Type: "string"},
				"$.crops": {Type: "array", Items: []byte(`{
			        "type": "object",
			        "properties": {
			          "name": {
			            "type": "string",
			            "default": "name"
			          },
			          "width": {
			            "type": "number"
			          },
			          "height": {
			            "type": "number"
			          },
			          "path": {
			            "type": "string"
			          },
			          "relativePath": {
			            "type": "string"
			          }
			        },
			        "required":[
			          "name",
			          "width",
			          "height",
			          "path",
			          "relativePath"
			        ]
			    }`)},
				"$.crops[*]": {Type: "object", Properties: map[string]json.RawMessage{
					"name": []byte(`{
			            "type": "string",
			            "default": "name"
			          }`),
					"width": []byte(`{
			            "type": "number"
			          }`),
					"height": []byte(`{
			            "type": "number"
			          }`),
					"path": []byte(`{
			            "type": "string"
			          }`),
					"relativePath": []byte(`{
			            "type": "string"
			          }`),
				},
					Required: []string{"name", "width", "height", "path", "relativePath"},
				},
				"$.crops[*].name":         {Type: "string"},
				"$.crops[*].width":        {Type: "number"},
				"$.crops[*].height":       {Type: "number"},
				"$.crops[*].path":         {Type: "string"},
				"$.crops[*].relativePath": {Type: "string"},
				"$.URL": {Type: "object", Properties: map[string]json.RawMessage{
					"publish": []byte(`{
			          "type": "string",
			          "transform": {
			            "cumulo": {
			              "from" : [
			                {
			                  "jsonPath": "$.publishUrl"
			                }
			              ]
			            }
			          }
			        }`),
					"absolute": []byte(`{
			          "type": "string",
			          "transform": {
			            "cumulo": {
			              "from" : [
			                {
			                  "jsonPath": "$.absoluteUrl"
			                }
			              ]
			            }
			          }
			      }`)},
					Required: []string{"publish", "absolute"},
				},
				"$.URL.publish":  {Type: "string"},
				"$.URL.absolute": {Type: "string"},
			},
		},
		{
			description: "Advanced walk does it all",
			oneOfType:   "array-of-array",
			schemaPath:  "./test_data/parent.json",
			want: map[string]Instance{
				"$.type": {Type: "string"},
				"$.crops": {Type: "array", Items: []byte(`{
			        "type": "array",
			        "items": {
			          "type": "object",
			          "properties": {
			            "name": {
			              "type": "string"
			            }
			          }
			        }
			    }`)},
				"$.crops[*]": {Type: "array", Items: []byte(`{
			          "type": "object",
			          "properties": {
			            "name": {
			              "type": "string"
			            }
			          }
			      }`)},
				"$.crops[*][*]": {Type: "object", Properties: map[string]json.RawMessage{
					"name": []byte(`{
			              "type": "string"
			        }`)},
				},
				"$.crops[*][*].name": {Type: "string"},
			},
		},
		{
			description: "Object with missing properties",
			schemaPath:  "./test_data/bad-object.json",
			wantErr:     true,
		},
		{
			description: "Array with missing Items",
			schemaPath:  "./test_data/bad-array.json",
			wantErr:     true,
		},
	}

	for _, test := range tests {
		walker := testWalker{calls: make(map[string]Instance), rawCalls: make(map[string]json.RawMessage)}
		schema, err := SchemaFromFile(test.schemaPath, test.oneOfType)
		if err != nil {
			t.Fatal(err)
		}
		err = Walk(schema, walker.walkFn)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil error want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error: %v", test.description, err)
			continue
		}
		if got, want := len(walker.calls), len(test.want); got != want {
			t.Errorf("Test %q - got %d calls, want %d", test.description, got, want)
		}
		for key, call := range walker.calls {
			got, err := json.Marshal(call)
			if err != nil {
				t.Fatal(err)
			}
			want, err := json.Marshal(test.want[key])
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Test %q - at got key %q got call\n%s\n\twant\n%s", test.description, key, call, want)
			}
		}
		for key, call := range test.want {
			got, err := json.Marshal(walker.calls[key])
			if err != nil {
				t.Fatal(err)
			}
			want, err := json.Marshal(call)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Test %q - at want key %q got call\n%s\n\twant\n%s", test.description, key, got, want)
			}
		}
	}
}
func TestWalkJSONSchemaRaw(t *testing.T) {
	tests := []struct {
		description string
		oneOfType   string
		schemaPath  string
		want        map[string]json.RawMessage
		wantErr     bool
	}{
		{
			description: "Basic walk, no allOf, no oneOf",
			schemaPath:  "./test_data/image.json",
			want: map[string]json.RawMessage{
				"$.type": []byte(`{
      "type": "string",
      "enum": [
        "image"
      ]
    }`),
				"$.crops": []byte(`{
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
      }
    }`),
				"$.crops[*]": []byte(`{
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
      }`),
				"$.crops[*].name": []byte(`{
            "type": "string",
            "default": "name"
          }`),
				"$.crops[*].width": []byte(`{
            "type": "number"
          }`),
				"$.crops[*].height": []byte(`{
            "type": "number"
          }`),
				"$.crops[*].path": []byte(`{
            "type": "string"
          }`),
				"$.crops[*].relativePath": []byte(`{
            "type": "string"
          }`),
				"$.URL": []byte(`{
      "type": "object",
      "properties": {
        "publish": {
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }
        },
        "absolute": {
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }
        }
      },
      "required":[
        "publish",
        "absolute"
      ]
    }`),
				"$.URL.publish": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }
        }`),
				"$.URL.absolute": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }
        }`),
			},
		},
		{
			description: "Walk with allOf, no oneOf",
			schemaPath:  "./test_data/embed_parent.json",
			want: map[string]json.RawMessage{"$.type": []byte(`{
      "type": "string",
      "enum": [
        "embed"
      ]
    }`),
			},
		},
		{
			description: "Walk with oneOf, no allOf",
			oneOfType:   "image",
			schemaPath:  "./test_data/image_parent.json",
			want: map[string]json.RawMessage{
				"$.type": []byte(`{
      "type": "string",
      "enum": [
        "image"
      ]
    }`),
				"$.crops": []byte(`{
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
      }
    }`),
				"$.crops[*]": []byte(`{
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "default": "name"
          },
          "width": {
            "type": "number"
          },
          "height": {
            "type": "number"
          },
          "path": {
            "type": "string"
          },
          "relativePath": {
            "type": "string"
          }
        },
        "required":[
          "name",
          "width",
          "height",
          "path",
          "relativePath"
        ]
      }`),
				"$.crops[*].name": []byte(`{
            "type": "string",
            "default": "name"
          }`),
				"$.crops[*].width": []byte(`{
            "type": "number"
          }`),
				"$.crops[*].height": []byte(`{
            "type": "number"
          }`),
				"$.crops[*].path": []byte(`{
            "type": "string"
          }`),
				"$.crops[*].relativePath": []byte(`{
            "type": "string"
          }`),
				"$.URL": []byte(`{
      "type": "object",
      "properties": {
        "publish": {
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }
        },
        "absolute": {
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }
        }
      },
      "required":[
        "publish",
        "absolute"
      ]
    }`),
				"$.URL.publish": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.publishUrl"
                }
              ]
            }
          }
        }`),
				"$.URL.absolute": []byte(`{
          "type": "string",
          "transform": {
            "cumulo": {
              "from" : [
                {
                  "jsonPath": "$.absoluteUrl"
                }
              ]
            }
          }
        }`),
			},
		},
		{
			description: "Advanced walk does it all",
			oneOfType:   "array-of-array",
			schemaPath:  "./test_data/parent.json",
			want: map[string]json.RawMessage{
				"$.type": []byte(`{
      "type": "string",
      "enum": [
        "array-of-array"
      ]
    }`),
				"$.crops": []byte(`{
      "type": "array",
      "items": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            }
          }
        }
      }
    }`),
				"$.crops[*]": []byte(`{
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            }
          }
        }
      }`),
				"$.crops[*][*]": []byte(`{
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            }
          }
        }`),
				"$.crops[*][*].name": []byte(`{
              "type": "string"
            }`),
			},
		},
		{
			description: "Object with missing properties",
			schemaPath:  "./test_data/bad-object.json",
			wantErr:     true,
		},
		{
			description: "Array with missing Items",
			schemaPath:  "./test_data/bad-array.json",
			wantErr:     true,
		},
	}

	for _, test := range tests {
		walker := testWalker{calls: make(map[string]Instance), rawCalls: make(map[string]json.RawMessage)}
		schema, err := SchemaFromFile(test.schemaPath, test.oneOfType)
		if err != nil {
			t.Fatal(err)
		}
		err = Walk(schema, walker.walkFn)

		switch {
		case test.wantErr && err != nil:
			continue
		case test.wantErr && err == nil:
			t.Errorf("Test %q - got nil error want error", test.description)
		case !test.wantErr && err != nil:
			t.Errorf("Test %q - got error: %v", test.description, err)
			continue
		}
		if got, want := len(walker.rawCalls), len(test.want); got != want {
			t.Errorf("Test %q - got %d calls, want %d", test.description, got, want)
		}
		for key, call := range walker.rawCalls {
			if !reflect.DeepEqual(call, test.want[key]) {
				t.Errorf("Test %q - at got key %q got call\n%s\n\twant\n%s", test.description, key, call, test.want[key])
			}
		}
		for key, call := range test.want {
			if !reflect.DeepEqual(call, walker.rawCalls[key]) {
				t.Errorf("Test %q - at want key %q got call\n%s\n\twant\n%s", test.description, key, walker.rawCalls[key], call)
			}
		}
	}
}

func BenchmarkWalk(b *testing.B) {
	schema, err := SchemaFromFile("./test_data/image_parent.json", "image")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		if err := Walk(schema, func(path string, i Instance, value json.RawMessage) error { return nil }); err != nil {
			b.Fatal(err)
		}
	}
}
