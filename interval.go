package main

type LiveInterval struct {
	symbol     Symbol
	location   int
	startPoint int
	endPoint   int
}

type ActiveLiveIntervals []*LiveInterval

type ByStartPoint []LiveInterval

func (a ByStartPoint) Len() int           { return len(a) }
func (a ByStartPoint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartPoint) Less(i, j int) bool { return a[i].startPoint < a[j].startPoint }

func insertByEndPoint(active ActiveLiveIntervals, iv *LiveInterval) ActiveLiveIntervals {
	i := 0
	for i < len(active) && active[i].endPoint <= iv.endPoint {
		i++
	}
	active = append(active, iv)
	copy(active[i+1:], active[i:len(active)-1])
	active[i] = iv
	return active
}
