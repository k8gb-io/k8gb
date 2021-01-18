// Package utils provides common functionality to gslb controller
package utils

import (
	"encoding/json"
	"fmt"
)

// ToString converts type to formatted string. If value is struct, function returns formatted JSON. Function retrieves
// null for nil pointer references. Function doesn't return error. In case of marshal error it converts with %v formatter
// Only two possible errors can occur e.g.:
//	UnsupportedTypeError ToString(make(chan int));
//	UnsupportedValueError ToString(math.Inf(1));
//	In both cases function retrieves expected result. The pointer address in the first while "+Inf" in second
func ToString(v interface{}) string {
	value, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(value)
}
