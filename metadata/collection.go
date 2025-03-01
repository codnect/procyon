package metadata

type Collection map[Key]Metadata

func From[M Metadata](collection Collection, key Key) (M, bool) {
	var metadata M
	if result, ok := collection[key]; ok {
		return metadata, false
	} else {
		var cast bool
		metadata, cast = result.(M)
		if cast {
			return metadata, true
		}
	}

	return metadata, false
}
