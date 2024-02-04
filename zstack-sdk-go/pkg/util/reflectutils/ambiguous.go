package reflectutils

import (
	"fmt"
)

const (
	TagAmbiguousPrefix = "zstack-ambiguous-prefix"
	TagDeprecatedBy    = "zstack-deprecated-by"
	TagOldDeprecatedBy = "deprecated-by"
)

func expandAmbiguousPrefix(fields SStructFieldValueSet) SStructFieldValueSet {
	keyIndexMap := make(map[string][]int)
	for i := range fields {
		if fields[i].Info.Ignore {
			continue
		}
		key := fields[i].Info.MarshalName()
		values, ok := keyIndexMap[key]
		if !ok {
			values = make([]int, 0, 2)
		}
		keyIndexMap[key] = append(values, i)
	}
	for _, indexes := range keyIndexMap {
		if len(indexes) > 1 {
			// ambiguous found
			for _, idx := range indexes {
				if amPrefix, ok := fields[idx].Info.Tags[TagAmbiguousPrefix]; ok {
					fields[idx].Info.Name = fmt.Sprintf("%s%s", amPrefix, fields[idx].Info.Name)
					if depBy, ok := fields[idx].Info.Tags[TagDeprecatedBy]; ok {
						fields[idx].Info.Tags[TagDeprecatedBy] = fmt.Sprintf("%s%s", amPrefix, depBy)
					}
					if depBy, ok := fields[idx].Info.Tags[TagOldDeprecatedBy]; ok {
						fields[idx].Info.Tags[TagOldDeprecatedBy] = fmt.Sprintf("%s%s", amPrefix, depBy)
					}
				}
			}
		}
	}
	return fields
}
