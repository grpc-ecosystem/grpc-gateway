package genswagger

// This patch is introduced for several cases that are applicable in
// atlas-app-toolkit:
//
// - Ability to wrap Responses with correct error codes (200 - for GET, 201 - for POST/PUT/PATCH, 204 - for DELETE)
//
// - Ability to identify and append correct documentation with atlas.app.toolkit
// special types: filtering, sorting, paging, field_selection, atlas.rpc.identifier.
//
// - Ability to break up recursive rules introduced by many-to-many definitions.
//
// - Unused refs removal.

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-openapi/spec"
)

const (
	titleAnnotation   = "@title"
	exampleAnnotation = "@example"
)

var (
	sw       spec.Swagger
	seenRefs = map[string]bool{}
)

// filterPathVars returns new params list with: required "true" and path "in" variables only if they present
// in URL path request
func filterPathVars(path string, params []spec.Parameter) []spec.Parameter {
	var newParams []spec.Parameter
	for _, param := range params {
		if !param.ParamProps.Required || !(param.ParamProps.In == "path") || strings.Contains(path, fmt.Sprintf("{%s}", param.ParamProps.Name)) {
			newParams = append(newParams, param)
		}
	}

	return newParams
}

func atlasSwagger(b []byte, withPrivateMethods, withCustomAnnotations bool) string {
	if err := json.Unmarshal(b, &sw); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// remove params that are not part of path
	for path, item := range sw.Paths.Paths {
		if item.Get != nil {
			newParams := filterPathVars(path, item.Get.Parameters)
			item.Get.Parameters = newParams
		}
		if item.Post != nil {
			newParams := filterPathVars(path, item.Post.Parameters)
			item.Post.Parameters = newParams
		}
		if item.Put != nil {
			newParams := filterPathVars(path, item.Put.Parameters)
			item.Put.Parameters = newParams
		}
		if item.Patch != nil {
			newParams := filterPathVars(path, item.Patch.Parameters)
			item.Patch.Parameters = newParams
		}
		if item.Delete != nil {
			newParams := filterPathVars(path, item.Delete.Parameters)
			item.Delete.Parameters = newParams
		}
		if item.Head != nil {
			newParams := filterPathVars(path, item.Head.Parameters)
			item.Head.Parameters = newParams
		}
		if item.Options != nil {
			newParams := filterPathVars(path, item.Options.Parameters)
			item.Options.Parameters = newParams
		}
	}

	// Fix collection operators and IDs and gather refs along Paths.

	var refs []spec.Ref
	fixedPaths := map[string]spec.PathItem{}
	privateMethodsOperations := make(map[string][]string, 0)
	for pn, pi := range sw.Paths.Paths {
		var pnElements []string
		for _, v := range strings.Split(pn, "/") {
			if strings.HasSuffix(v, "id.resource_id}") || strings.HasSuffix(v, ".id}") {
				pnElements = append(pnElements, "{id}")
			} else {
				pnElements = append(pnElements, v)
			}
		}
		pn := strings.Join(pnElements, "/")
		for on, op := range pathItemAsMap(pi) {
			if op == nil {
				continue
			}
			if !withPrivateMethods {
				if IsStringInSlice(op.OperationProps.Tags, "private") {
					privateMethodsOperations[pn] = append(privateMethodsOperations[pn], on)
				}
			}

			var fixedParams []spec.Parameter
			for _, param := range op.Parameters {

				// Fix Collection Operators
				if strings.HasPrefix(param.Description, "atlas.api.") {
					switch strings.TrimPrefix(param.Description, "atlas.api.") {

					case "filtering":
						fixedParams = append(fixedParams, *(spec.QueryParam("_filter")).WithDescription(`

A collection of response resources can be filtered by a logical expression string that includes JSON tag references to values in each resource, literal values, and logical operators. If a resource does not have the specified tag, its value is assumed to be null.

Literal values include numbers (integer and floating-point), and quoted (both single- or double-quoted) literal strings, and 'null'. The following operators are commonly used in filter expressions:

|  Op   |  Description               | 
|  --   |  -----------               | 
|  ==   |  Equal                     | 
|  !=   |  Not Equal                 | 
|  >    |  Greater Than              | 
|   >=  |  Greater Than or Equal To  | 
|  <    |  Less Than                 | 
|  <=   |  Less Than or Equal To     | 
|  and  |  Logical AND               | 
|  ~    |  Matches Regex             | 
|  !~   |  Does Not Match Regex      | 
|  or   |  Logical OR                | 
|  not  |  Logical NOT               | 
|  ()   |  Groupping Operators       |

						`).Typed("string", ""))

					case "sorting":
						fixedParams = append(fixedParams, *(spec.QueryParam("_order_by")).WithDescription(`

A collection of response resources can be sorted by their JSON tags. For a 'flat' resource, the tag name is straightforward. If sorting is allowed on non-flat hierarchical resources, the service should implement a qualified naming scheme such as dot-qualification to reference data down the hierarchy. If a resource does not have the specified tag, its value is assumed to be null.)

Specify this parameter as a comma-separated list of JSON tag names. The sort direction can be specified by a suffix separated by whitespace before the tag name. The suffix 'asc' sorts the data in ascending order. The suffix 'desc' sorts the data in descending order. If no suffix is specified the data is sorted in ascending order.

						`).Typed("string", ""))

					case "field_selection":
						fixedParams = append(fixedParams, *(spec.QueryParam("_fields")).WithDescription(`

A collection of response resources can be transformed by specifying a set of JSON tags to be returned. For a “flat” resource, the tag name is straightforward. If field selection is allowed on non-flat hierarchical resources, the service should implement a qualified naming scheme such as dot-qualification to reference data down the hierarchy. If a resource does not have the specified tag, the tag does not appear in the output resource.

Specify this parameter as a comma-separated list of JSON tag names.

						`).Typed("string", ""))

					case "paging":
						fixedParams = append(
							fixedParams,
							*(spec.QueryParam("_offset")).WithDescription(`

The integer index (zero-origin) of the offset into a collection of resources. If omitted or null the value is assumed to be '0'.

							`).Typed("integer", ""),
							*(spec.QueryParam("_limit")).WithDescription(`

The integer number of resources to be returned in the response. The service may impose maximum value. If omitted the service may impose a default value.

							`).Typed("integer", ""),
							*(spec.QueryParam("_page_token")).WithDescription(`

The service-defined string used to identify a page of resources. A null value indicates the first page.

							`).Typed("string", ""),
						)
					// Skip ID
					default:
					}
					// Replace resource_id with id
				} else if strings.HasSuffix(param.Name, "id.resource_id") || strings.HasSuffix(param.Name, ".id") {
					param.Name = "id"
					fixedParams = append(fixedParams, param)
				} else if strings.HasPrefix(param.Description, "tagging.api.") {
					switch strings.TrimPrefix(param.Description, "tagging.api.") {
					case "filtering":
						fixedParams = append(fixedParams, *(spec.QueryParam("_tfilter")).WithDescription("This parameter is used for filtering by tags.").Typed("string", ""))
					case "sorting":
						fixedParams = append(fixedParams, *(spec.QueryParam("_torder_by")).WithDescription("This parameter is used for sorting by tags.").Typed("string", ""))
					default:
						fixedParams = append(fixedParams, param)
					}
				} else {
					// Gather ref in body.
					if param.In == "body" && param.Schema != nil {
						refs = append(refs, param.Schema.Ref)
					}
					fixedParams = append(fixedParams, param)
				}
			}
			op.Parameters = fixedParams

			// Wrap responses
			if op.Responses.StatusCodeResponses != nil {
				// check if StatusCodeResponses has 201 >= x < 300 then delete 200 and don't go to isNilRef check
				exists := false
				for code := range op.Responses.StatusCodeResponses {
					if code >= 201 && code < 300 {
						exists = true
					}
					break
				}
				if exists {
					delete(op.Responses.StatusCodeResponses, 200)
				} else {
					rsp := op.Responses.StatusCodeResponses[200]
					if !isNilRef(rsp.Schema.Ref) {
						s, _, err := rsp.Schema.Ref.GetPointer().Get(sw)
						if err != nil {
							panic(err)
						}

						schema := s.(spec.Schema)
						if schema.Properties == nil {
							schema.Properties = map[string]spec.Schema{}
						}

						def := sw.Definitions[trim(rsp.Schema.Ref)]
						if rsp.Description == "" {
							rsp.Description = on + " operation response"
						}

						switch on {
						case "DELETE":
							if len(def.Properties) == 0 {
								rsp.Description = "No Content"
								rsp.Schema = nil
								op.Responses.StatusCodeResponses[opToStatusCode(on)] = rsp
								delete(sw.Definitions, trim(rsp.Ref))
								delete(op.Responses.StatusCodeResponses, 200)
								break
							}
							sw.Definitions[trim(rsp.Schema.Ref)] = schema
							refs = append(refs, rsp.Schema.Ref)
							op.Responses.StatusCodeResponses[200] = rsp
						default:
							sw.Definitions[trim(rsp.Schema.Ref)] = schema
							refs = append(refs, rsp.Schema.Ref)
							delete(op.Responses.StatusCodeResponses, 200)
							op.Responses.StatusCodeResponses[opToStatusCode(on)] = rsp
						}
					}
				}
			}

			op.ID = strings.Join(op.Tags, "") + op.ID

		}

		pitem := fixedPaths[pn]
		for opName, opPtr := range pathItemAsMap(pi) {
			if opPtr == nil {
				continue
			}
			opPtr := opPtr
			switch opName {
			case "GET":
				pitem.Get = opPtr
			case "PUT":
				pitem.Put = opPtr
			case "POST":
				pitem.Post = opPtr
			case "DELETE":
				pitem.Delete = opPtr
			case "PATCH":
				pitem.Patch = opPtr
			}
		}
		fixedPaths[pn] = pitem
	}

	sw.Paths.Paths = fixedPaths

	// Break recursive rules introduced by many-to-many.
	for _, r := range refs {
		seenRefs[trim(r)] = true
		s, _, err := r.GetPointer().Get(sw)
		if err != nil {
			continue
		}

		if _, ok := s.(spec.Schema); ok {
			checkRecursion(s.(spec.Schema), r, []string{})
		}
	}

	// Cleanup unused definitions.
	for dn, v := range sw.Definitions {
		// hidden definitions should become explicit.
		if strings.HasPrefix(dn, "_") {
			sw.Definitions[strings.TrimPrefix(dn, "_")] = v
			delete(sw.Definitions, dn)
			seenRefs[dn] = true
		}

		if seenRefs[dn] == false {
			delete(sw.Definitions, dn)
		}
	}

	for pn, on := range privateMethodsOperations {
		pi := sw.Paths.Paths[pn]
		for _, operation := range on {
			switch operation {
			case "GET":
				pi.Get = nil
			case "POST":
				pi.Post = nil
			case "PUT":
				pi.Put = nil
			case "DELETE":
				pi.Delete = nil
			case "PATCH":
				pi.Patch = nil
			}
		}

		if IsPathEmpty(pi) {
			delete(sw.Paths.Paths, pn)
			continue
		}

		sw.Paths.Paths[pn] = pi
	}

	if !withPrivateMethods {
		sw.Definitions = filterDefinitions()
	}

	if withCustomAnnotations {
		sw.Definitions = applyCustomAnnotations(sw.Definitions)
	}

	bOut, err := json.MarshalIndent(sw, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling result: %v", err)
		os.Exit(1)
	}

	return fmt.Sprintf("%s", bOut)
}

func getPropRef(p spec.Schema) spec.Ref {
	if len(p.Type) == 1 && p.Type[0] == "array" {
		return p.Items.Schema.Ref
	}

	return p.Ref
}

func setPropRef(p *spec.Schema, r spec.Ref) {
	if len(p.Type) == 1 && p.Type[0] == "array" {
		p.Items.Schema.Ref = r
	} else {
		p.Ref = r
	}
}

func checkRecursion(s spec.Schema, r spec.Ref, path []string) spec.Ref {
	var newRefLength int
	var newRefName string

	var newProps = map[string]spec.Schema{}

	npath := path[:]
	npath = append(npath, trim(r))

	newProps = map[string]spec.Schema{}
	for np, p := range s.Properties {
		if p.Description == "atlas.api.identifier" {
			p.Description = "The resource identifier."
			if np == "id" {
				p.ReadOnly = true
			}
		}

		// TBD: common pattern.
		if np == "created_time" || np == "updated_time" || np == "id" {
			p.ReadOnly = true
		}

		// FIXME: copy additionalProperties as-is.
		if addProps := p.AdditionalProperties; addProps != nil {
			if addProps.Schema != nil && !isNilRef(addProps.Schema.Ref) {
				seenRefs[trim(addProps.Schema.Ref)] = true
			}
		}

		newProps[np] = p

		sr := getPropRef(p)

		if isNilRef(sr) {
			continue
		}

		for i, prefs := range npath {
			if trim(sr) == prefs {
				delete(newProps, np)
				if newRefLength < len(npath)-i {
					newRefName = strings.Join(reverse(npath[i:]), "_In_")
					newRefLength = len(npath) - i
				}
			}
		}

		if _, ok := newProps[np]; !ok {
			continue
		}

		ss, _, _ := sr.GetPointer().Get(sw)
		if _, ok := ss.(spec.Schema); !ok {
			continue
		}

		nr := checkRecursion(ss.(spec.Schema), sr, npath)

		if trim(nr) != trim(sr) {
			if newRefName == "" {
				newRefName = strings.TrimPrefix(trim(nr), trim(sr)+"_In_")
			}

			delete(newProps, np)

			if len(p.Type) == 1 && p.Type[0] == "array" {
				newProps[np] = *spec.ArrayProperty(spec.RefProperty(nr.String()))
			} else {
				newProps[np] = *spec.RefProperty(nr.String())
			}
		} else {
			seenRefs[trim(sr)] = true
		}
	}

	if newRefName != "" {
		seenRefs[newRefName] = true
		// underscore hides definitions from following along recursive path.
		sw.Definitions["_"+newRefName] = *(&spec.Schema{}).WithProperties(newProps)
		return spec.MustCreateRef("#/definitions/" + newRefName)
	} else {
		s.Properties = newProps
		sw.Definitions[trim(r)] = s
	}

	return r
}

func trim(r spec.Ref) string {
	return strings.TrimPrefix(r.String(), "#/definitions/")
}

func isNilRef(r spec.Ref) bool {
	return r.String() == ""
}

func reverse(s []string) []string {
	news := make([]string, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		news[i] = s[len(s)-1-i]
	}

	return news
}

func pathItemAsMap(pi spec.PathItem) map[string]*spec.Operation {
	return map[string]*spec.Operation{
		"GET":    pi.Get,
		"POST":   pi.Post,
		"PUT":    pi.Put,
		"DELETE": pi.Delete,
		"PATCH":  pi.Patch,
	}
}

func opToStatusCode(on string) int {
	return map[string]int{
		"GET":    200,
		"POST":   201,
		"PUT":    201,
		"PATCH":  201,
		"DELETE": 204,
	}[on]
}

func opToTextCode(on string) string {
	return map[string]string{
		"GET":    "OK",
		"POST":   "CREATED",
		"PUT":    "UPDATED",
		"PATCH":  "UPDATED",
		"DELETE": "DELETED",
	}[on]
}

func IsStringInSlice(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}

	return false
}

func IsPathEmpty(pi spec.PathItem) bool {
	if pi.Get != nil || pi.Post != nil || pi.Put != nil || pi.Patch != nil || pi.Delete != nil {
		return false
	}
	return true
}

func filterDefinitions() (newDefinitions spec.Definitions) {
	marh, _ := sw.MarshalJSON()
	v := map[string]interface{}{}
	if err := json.Unmarshal(marh, &v); err != nil {
		panic(err.Error())
	}
	defs, _ := v["definitions"].(map[string]interface{})
	newDefinitions = make(spec.Definitions)

	for rk := range gatherRefs(v["paths"]) {
		rName := refToName(rk)
		newDefinitions[rName] = sw.Definitions[rName]
		for rrName := range gatherDefinitionRefs(rk, defs) {
			newDefinitions[rrName] = sw.Definitions[rrName]
		}
	}

	return newDefinitions
}

func gatherDefinitionRefs(ref string, defs map[string]interface{}) map[string]struct{} {
	var refs = make(map[string]struct{})

	gatherDefinitionRefsAux(refToName(ref), defs, refs)
	return refs
}

func gatherDefinitionRefsAux(ref string, defs map[string]interface{}, refs map[string]struct{}) {
	for r := range gatherRefs(defs[ref]) {
		refs[r] = struct{}{}
		gatherDefinitionRefsAux(r, defs, refs)
	}

	return
}

func gatherRefs(v interface{}) map[string]struct{} {
	refs := map[string]struct{}{}
	switch v := v.(type) {
	case map[string]interface{}:
		for k, vv := range v {
			if k == "$ref" {
				refs[refToName(vv.(string))] = struct{}{}
			}

			for rk := range gatherRefs(vv) {
				refs[rk] = struct{}{}
			}
		}
	case []interface{}:
		for _, vv := range v {
			for rk := range gatherRefs(vv) {
				refs[rk] = struct{}{}
			}
		}
	}

	return refs
}

func refToName(ref string) string {
	return strings.TrimPrefix(ref, "#/definitions/")
}

func applyCustomAnnotations(defs spec.Definitions) spec.Definitions {
	for k, v := range defs {
		v.Description = v.Title + v.Description
		v.Title = ""

		if title, ok := getAnnotationValue(v.Description, titleAnnotation); ok {
			v.Title, ok = title.(string)
			if !ok {
				log.Printf("unsupported title type in %s\n", v.Type)
			}

			v.Description = removeSpecials(v.Description, titleAnnotation)

			defs[k] = v
		}

		if example, ok := getAnnotationValue(v.Description, exampleAnnotation); ok {
			v.Example = example
			v.Description = removeSpecials(v.Description, exampleAnnotation)

			defs[k] = v
		}

		for fk, fv := range v.Properties {
			fv.Description = fv.Title + fv.Description
			fv.Title = ""

			if example, ok := getAnnotationValue(fv.Description, exampleAnnotation); ok {
				fv.Example = example
				fv.Description = removeSpecials(fv.Description, exampleAnnotation)

				defs[k].Properties[fk] = fv
			}
		}

		if v.Definitions != nil {
			v.Definitions = applyCustomAnnotations(v.Definitions)
		}
	}

	return defs
}

func removeSpecials(comment, annotation string) (res string) {
	lines := strings.Split(comment, "\n")
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, annotation) && strings.HasSuffix(trimmed, "<<<EOF") {
			pos := stringPosition(lines, "EOF", i)
			if pos != -1 && pos < len(lines) {
				lines = append(lines[:i], lines[pos+1:]...)
			}
			lines = append(lines[:i])

		} else if strings.HasPrefix(trimmed, annotation) {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}

	return strings.Join(lines, "\n")
}

func getAnnotationValue(comment, prefix string) (interface{}, bool) {
	var (
		value string
		lines = strings.Split(comment, "\n")
	)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, prefix) && strings.HasSuffix(line, "<<<EOF") {
			pos := stringPosition(lines, "EOF", i)
			if pos != -1 && pos < len(lines) {
				value = strings.Join(lines[i+1:pos], "\n")
			} else {
				value = strings.Join(lines[i+1:], "\n")
			}

		} else if strings.HasPrefix(line, prefix) {
			value = strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}

		if value != "" {
			break
		}
	}

	if value != "" {
		var (
			raw json.RawMessage
		)

		if err := json.Unmarshal([]byte(value), &raw); err == nil {
			res, err := messageWalk(raw)
			if err != nil {
				log.Println(err.Error())
				return nil, false
			}

			return res, true
		}
	}

	return nil, false
}

func stringPosition(strings []string, str string, start int) (pos int) {
	if start >= len(strings) {
		return -1
	}

	for i := start; i < len(strings); i++ {
		if str == strings[i] {
			return i
		}
	}

	return -1
}

func messageWalk(j json.RawMessage) (interface{}, error) {
	var (
		mapVal   map[string]json.RawMessage
		sliceVal []json.RawMessage
		strVal   string
		numVal   float64
	)

	if err := json.Unmarshal(j, &mapVal); err == nil {
		res := make(map[string]interface{}, len(mapVal))
		for k, v := range mapVal {
			if res[k], err = messageWalk(v); err != nil {
				return nil, err
			}
		}

		return res, nil
	} else if err := json.Unmarshal(j, &sliceVal); err == nil {
		res := make([]interface{}, len(sliceVal))
		for k, v := range sliceVal {
			if res[k], err = messageWalk(v); err != nil {
				return nil, err
			}
		}

		return res, nil
	} else if err := json.Unmarshal(j, &numVal); err == nil {
		return numVal, nil
	} else if err := json.Unmarshal(j, &strVal); err == nil {
		return strVal, nil
	} else {
		return nil, fmt.Errorf("unsuported message type")
	}
}
