package genopenapi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/genproto/googleapis/api/visibility"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestIsVisible(t *testing.T) {
	t.Parallel()

	reg := descriptor.NewRegistry()

	// No annotation → always visible.
	if !isVisible(nil, reg) {
		t.Error("nil rule should be visible")
	}

	// Annotation present but no selectors configured (empty map) → hidden.
	if isVisible(&visibility.VisibilityRule{Restriction: "INTERNAL"}, reg) {
		t.Error("INTERNAL rule with no selectors should be hidden")
	}

	// Configure PREVIEW selector.
	reg.SetVisibilityRestrictionSelectors([]string{"PREVIEW"})

	// No annotation → still visible.
	if !isVisible(nil, reg) {
		t.Error("nil rule should be visible regardless of selectors")
	}

	// INTERNAL only → hidden (PREVIEW doesn't match).
	if isVisible(&visibility.VisibilityRule{Restriction: "INTERNAL"}, reg) {
		t.Error("INTERNAL rule with PREVIEW selector should be hidden")
	}

	// PREVIEW → visible.
	if !isVisible(&visibility.VisibilityRule{Restriction: "PREVIEW"}, reg) {
		t.Error("PREVIEW rule with PREVIEW selector should be visible")
	}

	// Multi-label with partial match → visible.
	if !isVisible(&visibility.VisibilityRule{Restriction: "INTERNAL,PREVIEW"}, reg) {
		t.Error("INTERNAL,PREVIEW rule with PREVIEW selector should be visible")
	}

	// Configure both selectors.
	reg.SetVisibilityRestrictionSelectors([]string{"INTERNAL", "PREVIEW"})
	if !isVisible(&visibility.VisibilityRule{Restriction: "INTERNAL"}, reg) {
		t.Error("INTERNAL rule with INTERNAL+PREVIEW selectors should be visible")
	}
}

// runWithSelectors loads req under the given selectors, runs Generate, and
// returns the emitted document as a decoded map.
func runWithSelectors(t *testing.T, req *pluginpb.CodeGeneratorRequest, selectors []string) map[string]any {
	t.Helper()
	reg := descriptor.NewRegistry()
	reg.SetVisibilityRestrictionSelectors(selectors)
	if err := reg.Load(req); err != nil {
		t.Fatalf("registry load: %v", err)
	}
	var targets []*descriptor.File
	for _, name := range req.FileToGenerate {
		f, err := reg.LookupFile(name)
		if err != nil {
			t.Fatalf("lookup %s: %v", name, err)
		}
		targets = append(targets, f)
	}
	out, err := Generate(reg, targets)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 output file, got %d", len(out))
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out[0].GetContent()), &doc); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, out[0].GetContent())
	}
	return doc
}

func TestGenerate_Visibility(t *testing.T) {
	t.Parallel()

	req := loadRequest(t, "testdata/visibility.prototext")

	t.Run("preview_selector", func(t *testing.T) {
		t.Parallel()
		doc := runWithSelectors(t, req, []string{"PREVIEW"})
		paths := doc["paths"].(map[string]any)

		// PublicMethod (no annotation) → visible.
		if _, ok := paths["/v1/public"]; !ok {
			t.Error("PublicMethod should be visible (no annotation)")
		}

		// InternalMethod (INTERNAL only) → hidden.
		if _, ok := paths["/v1/internal"]; ok {
			t.Error("InternalMethod should be hidden (INTERNAL, selector is PREVIEW)")
		}

		// PreviewMethod (PREVIEW) → visible.
		if _, ok := paths["/v1/preview"]; !ok {
			t.Error("PreviewMethod should be visible (PREVIEW matches selector)")
		}

		// InternalService (INTERNAL) → entirely hidden.
		if _, ok := paths["/v1/secret"]; ok {
			t.Error("InternalService.SecretMethod should be hidden (service is INTERNAL)")
		}

		// Tags: only PublicService should appear (it has visible methods).
		tags := doc["tags"].([]any)
		for _, tag := range tags {
			m := tag.(map[string]any)
			if m["name"] == "InternalService" {
				t.Error("InternalService tag should not appear")
			}
		}

		// Field visibility in the body="*" request body for PublicMethod:
		publicPost := paths["/v1/public"].(map[string]any)
		postOp := publicPost["post"].(map[string]any)
		reqBody := postOp["requestBody"].(map[string]any)
		content := reqBody["content"].(map[string]any)
		appJSON := content["application/json"].(map[string]any)
		bodySchema := appJSON["schema"].(map[string]any)
		props := bodySchema["properties"].(map[string]any)

		if _, ok := props["publicField"]; !ok {
			t.Error("publicField should be visible (no annotation)")
		}
		if _, ok := props["internalField"]; ok {
			t.Error("internalField should be hidden (INTERNAL, selector is PREVIEW)")
		}
		if _, ok := props["previewField"]; !ok {
			t.Error("previewField should be visible (PREVIEW matches)")
		}
		if _, ok := props["multiField"]; !ok {
			t.Error("multiField should be visible (INTERNAL,PREVIEW — PREVIEW matches)")
		}

		// Enum value visibility: STATUS_INTERNAL_ONLY should be hidden.
		components := doc["components"].(map[string]any)
		schemas := components["schemas"].(map[string]any)
		statusEnum := schemas["vis.v1.Status"].(map[string]any)
		enumVals := statusEnum["enum"].([]any)
		for _, v := range enumVals {
			if v == "STATUS_INTERNAL_ONLY" {
				t.Error("STATUS_INTERNAL_ONLY should be hidden (INTERNAL, selector is PREVIEW)")
			}
		}
	})

	t.Run("no_selectors", func(t *testing.T) {
		t.Parallel()
		// No selectors: annotated elements are hidden, unannotated are visible.
		doc := runWithSelectors(t, req, []string{})
		paths := doc["paths"].(map[string]any)

		// PublicMethod (no annotation) → visible.
		if _, ok := paths["/v1/public"]; !ok {
			t.Error("PublicMethod should be visible (no annotation)")
		}

		// All annotated methods → hidden.
		if _, ok := paths["/v1/internal"]; ok {
			t.Error("InternalMethod should be hidden (has annotation, no matching selectors)")
		}
		if _, ok := paths["/v1/preview"]; ok {
			t.Error("PreviewMethod should be hidden (has annotation, no matching selectors)")
		}
		if _, ok := paths["/v1/secret"]; ok {
			t.Error("InternalService should be hidden (has annotation, no matching selectors)")
		}
	})

	t.Run("all_selectors", func(t *testing.T) {
		t.Parallel()
		// When both INTERNAL and PREVIEW selectors are configured,
		// all annotated elements should be visible.
		doc := runWithSelectors(t, req, []string{"INTERNAL", "PREVIEW"})
		paths := doc["paths"].(map[string]any)

		// All methods visible when all selectors configured.
		for _, p := range []string{"/v1/public", "/v1/internal", "/v1/preview", "/v1/secret"} {
			if _, ok := paths[p]; !ok {
				t.Errorf("%s should be visible when all selectors configured", p)
			}
		}
	})

	t.Run("oneof_filtering", func(t *testing.T) {
		t.Parallel()
		// OneofRequest has three members: visible (no annotation), internal,
		// preview. Under PREVIEW, visible + preview remain → the component
		// schema's oneOf constraint still fires with two members.
		doc := runWithSelectors(t, req, []string{"PREVIEW"})
		components := doc["components"].(map[string]any)
		schemas := components["schemas"].(map[string]any)
		schema, ok := schemas["vis.v1.OneofRequest"].(map[string]any)
		if !ok {
			t.Fatalf("OneofRequest component schema missing; schemas=%v", schemas)
		}
		props := schema["properties"].(map[string]any)

		if _, ok := props["visibleChoice"]; !ok {
			t.Error("visibleChoice should be present (no annotation)")
		}
		if _, ok := props["internalChoice"]; ok {
			t.Error("internalChoice should be filtered out (INTERNAL under PREVIEW)")
		}
		if _, ok := props["previewChoice"]; !ok {
			t.Error("previewChoice should be present (PREVIEW matches)")
		}

		// A single non-empty group is hoisted directly onto the schema's
		// oneOf: one "none set" guard plus one option per visible member.
		oneOf, ok := schema["oneOf"].([]any)
		if !ok {
			t.Fatalf("oneOf missing from OneofRequest schema: %v", schema)
		}
		if len(oneOf) != 3 {
			t.Errorf("expected 3 oneOf entries (guard + 2 visible members), got %d", len(oneOf))
		}
		raw, _ := json.Marshal(oneOf)
		if strings.Contains(string(raw), "internalChoice") {
			t.Errorf("oneOf should not reference internalChoice, got %s", raw)
		}
	})

	t.Run("oneof_collapse_to_single_visible", func(t *testing.T) {
		t.Parallel()
		// OneofCollapseRequest has one visible + one INTERNAL member. Under
		// PREVIEW, only one visible member remains, so the oneOf constraint
		// should collapse to nothing on the component schema.
		doc := runWithSelectors(t, req, []string{"PREVIEW"})
		components := doc["components"].(map[string]any)
		schemas := components["schemas"].(map[string]any)
		schema, ok := schemas["vis.v1.OneofCollapseRequest"].(map[string]any)
		if !ok {
			t.Fatalf("OneofCollapseRequest component schema missing; schemas=%v", schemas)
		}
		props := schema["properties"].(map[string]any)

		if _, ok := props["lonelyVisible"]; !ok {
			t.Error("lonelyVisible should be present")
		}
		if _, ok := props["lonelyInternal"]; ok {
			t.Error("lonelyInternal should be filtered out")
		}
		if _, ok := schema["oneOf"]; ok {
			t.Error("oneOf constraint should collapse when only one member is visible")
		}
		if _, ok := schema["allOf"]; ok {
			t.Error("allOf constraint should not appear when the only group collapses")
		}
	})

	t.Run("body_field_component_filtering", func(t *testing.T) {
		t.Parallel()
		// BodyFieldMethod uses body="payload" so the body is a $ref to the
		// BodyPayload component. Hidden fields on BodyPayload must not leak
		// via that component.
		doc := runWithSelectors(t, req, []string{"PREVIEW"})
		components := doc["components"].(map[string]any)
		schemas := components["schemas"].(map[string]any)
		payload, ok := schemas["vis.v1.BodyPayload"].(map[string]any)
		if !ok {
			t.Fatalf("expected BodyPayload component schema; got schemas=%v", schemas)
		}
		props := payload["properties"].(map[string]any)
		if _, ok := props["title"]; !ok {
			t.Error("title should be present on BodyPayload")
		}
		if _, ok := props["hiddenTitle"]; ok {
			t.Error("hiddenTitle should be filtered out of BodyPayload component")
		}

		// Also verify the request body on the operation references the
		// component (not an inline object with the hidden field).
		paths := doc["paths"].(map[string]any)
		p := paths["/v1/body_field/{id}"].(map[string]any)
		op := p["post"].(map[string]any)
		body := op["requestBody"].(map[string]any)
		bodySchema := body["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
		if bodySchema["$ref"] != "#/components/schemas/vis.v1.BodyPayload" {
			t.Errorf("expected body to $ref BodyPayload, got %v", bodySchema)
		}
	})

	t.Run("query_parameter_filtering", func(t *testing.T) {
		t.Parallel()
		// QueryMethod takes a scalar, a hidden scalar, and a nested message
		// (which itself has a hidden field). Both the top-level hidden
		// scalar and the nested hidden field should be absent from the
		// emitted parameters.
		doc := runWithSelectors(t, req, []string{"PREVIEW"})
		paths := doc["paths"].(map[string]any)
		p := paths["/v1/query"].(map[string]any)
		op := p["get"].(map[string]any)
		params, _ := op["parameters"].([]any)

		names := map[string]bool{}
		for _, raw := range params {
			pm := raw.(map[string]any)
			names[pm["name"].(string)] = true
		}

		if !names["visibleQ"] {
			t.Error("visibleQ should be emitted as a query parameter")
		}
		if names["internalQ"] {
			t.Error("internalQ should be filtered out of query parameters")
		}
		if !names["filter.label"] {
			t.Error("filter.label should be emitted (nested visible field)")
		}
		if names["filter.secretLabel"] {
			t.Error("filter.secretLabel should be filtered out (nested hidden field)")
		}
	})

	t.Run("enum_all_values_hidden", func(t *testing.T) {
		t.Parallel()
		// Under PREVIEW, every value of AllHiddenEnum (all marked INTERNAL)
		// is hidden. The component schema must still exist (a visible field
		// references it), but should fall back to unconstrained `type:
		// string` — no empty `enum: []` (which would be invalid), and no
		// dangling $ref.
		doc := runWithSelectors(t, req, []string{"PREVIEW"})
		components := doc["components"].(map[string]any)
		schemas := components["schemas"].(map[string]any)
		enumSchema, ok := schemas["vis.v1.AllHiddenEnum"].(map[string]any)
		if !ok {
			t.Fatalf("AllHiddenEnum component schema should be present as a $ref target; schemas=%v", schemas)
		}
		if enumSchema["type"] != "string" {
			t.Errorf("expected AllHiddenEnum to fall back to type=string, got %v", enumSchema["type"])
		}
		if _, hasEnum := enumSchema["enum"]; hasEnum {
			t.Errorf("AllHiddenEnum should not carry an enum constraint when all values are hidden, got %v", enumSchema)
		}
	})
}

