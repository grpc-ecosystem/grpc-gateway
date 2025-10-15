//go:build go1.12
// +build go1.12

package genopenapi

// this method will filter the same fields and return the unique one
func getUniqueFields(schemaFieldsRequired []string, fieldsRequired []string) []string {
	var unique []string
	var index *int

	for j, schemaFieldRequired := range schemaFieldsRequired {
		index = nil
		for i, fieldRequired := range fieldsRequired {
			i := i
			if schemaFieldRequired == fieldRequired {
				index = &i
				break
			}
		}
		if index == nil {
			unique = append(unique, schemaFieldsRequired[j])
		}
	}
	return unique
}
