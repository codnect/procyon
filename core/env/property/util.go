package property

import "strconv"

func flatMap(m map[string]interface{}) map[string]interface{} {
	flattenMap := map[string]interface{}{}

	for key, value := range m {
		switch child := value.(type) {
		case map[string]interface{}:
			nm := flatMap(child)

			for nk, nv := range nm {
				flattenMap[key+"."+nk] = nv
			}
		case []interface{}:
			for i := 0; i < len(child); i++ {
				flattenMap[key+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			flattenMap[key] = value
		}
	}

	return flattenMap
}
