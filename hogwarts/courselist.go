//go:build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) (ans []string) {
	ans = make([]string, 0)

	visited := make(map[string]bool)
	lived := make(map[string]bool)

	var topSort func(string)
	topSort = func(s string) {
		visited[s] = true
		for _, next := range prereqs[s] {
			if visited[next] && !lived[next] {
				panic("have a cycle")
			}
			if !visited[next] {
				topSort(next)
			}
		}
		ans = append(ans, s)
		lived[s] = true
	}

	for course := range prereqs {
		if !visited[course] {
			topSort(course)
		}
	}
	return
}
