package bot

import "strconv"

func ToArray(slice []int) string {
	var s string

	for i := 0; i < len(slice); i++ {
		if i > 0 {
			s += ", "
		}
		s += strconv.Itoa(slice[i])
	}
	return s

}
