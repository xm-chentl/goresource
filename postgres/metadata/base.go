package metadata

import "fmt"

func FormatField(field string) string {
	return fmt.Sprintf(`"%s"`, field)
}
