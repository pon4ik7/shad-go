//go:build !solution

package hotelbusiness

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	mpLoad := make(map[int]int)
	maxDate := 0
	for _, g := range guests {
		mpLoad[g.CheckInDate]++
		mpLoad[g.CheckOutDate]--
		maxDate = max(maxDate, max(g.CheckInDate, g.CheckOutDate))
	}
	cur := 0
	ans := make([]Load, 0)
	for i := 0; i <= maxDate; i++ {
		if mpLoad[i] != 0 {
			cur += mpLoad[i]
			ans = append(ans, Load{StartDate: i, GuestCount: cur})
		}
	}
	return ans
}
