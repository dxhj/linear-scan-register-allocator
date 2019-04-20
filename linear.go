package main

import (
	"fmt"
	"sort"
)

type Symbol struct {
	name     string
	register Register
	location int
}

var pool *RegisterPool = &RegisterPool{
	registers: []Register{EAX, EBX},
}

var currentStackLocation = 0

func expireOldIntervals(activeIntervals ActiveLiveIntervals, interval *LiveInterval) ActiveLiveIntervals {
	sort.Sort(ByEndPoint(activeIntervals))
	for i, activeInterval := range activeIntervals {
		if activeInterval.endPoint >= interval.startPoint {
			return activeIntervals
		}
		pool.freeRegister(activeInterval.symbol.register)
		return DeleteInterval(activeIntervals, i)
	}
	return activeIntervals
}

func spillAtInterval(activeIntervals ActiveLiveIntervals, interval *LiveInterval) {
	sort.Sort(ByEndPoint(activeIntervals))
	spill := activeIntervals[len(activeIntervals)-1]
	if spill.endPoint > interval.endPoint {
		fmt.Printf("ACTION: SPILL INTERVAL (%p)\n", spill)
		fmt.Printf("ACTION: ALLOCATE REGISTER %s(%d) TO INTERVAL(%p)\n", spill.symbol.register.getRegisterName(), spill.symbol.register, interval)
		interval.symbol.register = spill.symbol.register
		spill.symbol.location = currentStackLocation
		activeIntervals = DeleteInterval(activeIntervals, len(activeIntervals)-1)
	} else {
		fmt.Printf("ACTION: SPILL INTERVAL(%p)\n", interval)
		interval.symbol.location = currentStackLocation
	}
	currentStackLocation += 4
}

func main() {
	liveIntervals := []LiveInterval{
		LiveInterval{symbol: Symbol{name: "a", location: -1}, startPoint: 1, endPoint: 4},
		LiveInterval{symbol: Symbol{name: "b", location: -1}, startPoint: 2, endPoint: 6},
		LiveInterval{symbol: Symbol{name: "c", location: -1}, startPoint: 3, endPoint: 10},
		LiveInterval{symbol: Symbol{name: "d", location: -1}, startPoint: 5, endPoint: 9},
		LiveInterval{symbol: Symbol{name: "e", location: -1}, startPoint: 7, endPoint: 8},
	}
	var activeIntervals ActiveLiveIntervals

	sort.Sort(ByStartPoint(liveIntervals))

	for i, interval := range liveIntervals {
		fmt.Printf("(%p). SYMBOL: %s | LOCATION: %d | STARTPOINT: %d | ENDPOINT: %d\n", &liveIntervals[i], interval.symbol.name, interval.symbol.location, interval.startPoint, interval.endPoint)

		activeIntervals = expireOldIntervals(activeIntervals, &liveIntervals[i])

		if pool.isEmpty() {
			spillAtInterval(activeIntervals, &liveIntervals[i])
		} else {
			register, err := pool.getRegister()
			fmt.Printf("ACTION: ALLOCATE REGISTER %s(%d) TO INTERVAL(%p)\n", register.getRegisterName(), register, &liveIntervals[i])
			if err == nil {
				liveIntervals[i].symbol.register = register
			}
			activeIntervals = append(activeIntervals, &liveIntervals[i])
		}
	}

	for _, interval := range liveIntervals {
		fmt.Println(interval.symbol)
	}
}
