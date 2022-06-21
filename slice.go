package pgxquery

const (
	notFound = -1
)

func stringSliceFind(ss []string, str string) int {
	for i := range ss {
		if ss[i] != str {
			continue
		}

		return i
	}

	return notFound
}

func stringSliceRemove(ss []string, i int) []string {
	ss[i], ss[len(ss)-1] = ss[len(ss)-1], ss[i]
	ss = ss[:len(ss)-1]

	return ss
}
