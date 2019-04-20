package main

type LiveInterval struct {
	symbol     Symbol
	location   int
	startPoint int
	endPoint   int
}

type ActiveLiveIntervals []*LiveInterval

type ByStartPoint []LiveInterval
type ByEndPoint []*LiveInterval

func (a ByStartPoint) Len() int           { return len(a) }
func (a ByStartPoint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartPoint) Less(i, j int) bool { return a[i].startPoint < a[j].startPoint }

func (a ByEndPoint) Len() int           { return len(a) }
func (a ByEndPoint) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByEndPoint) Less(i, j int) bool { return a[i].endPoint < a[j].endPoint }

func DeleteInterval(intervals []*LiveInterval, index int) []*LiveInterval {
	return intervals[:index+copy(intervals[index:], intervals[index+1:])]
}
