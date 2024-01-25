package uts

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

func CommandLine(mapFunc map[string]func(), ordersDesc string) {
forLabel:
	for {
		reader := bufio.NewReader(os.Stdin)
		order, _ := reader.ReadString('\n')
		order = strings.ReplaceAll(order, "\r", "")
		order = strings.ReplaceAll(order, "\n", "")

		switch order {
		case "gc":
			printMemStats("After")
			runtime.GC()
			printMemStats("Befor")
		case "exit":
			os.Exit(0)
			break forLabel
		default:
			if f, ok := mapFunc[order]; ok {
				f()
			} else {
				Log(ordersDesc)
			}
		}
	}
}

func printMemStats(mark string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	Log("%v:Alloc = %vKB,TotalAlloc = %vKB,HeapAlloc = %vKB,HeapObjects = %vKB, GC Times = %vn", mark, m.Alloc/1024, m.TotalAlloc/1024, m.HeapAlloc/1024, m.HeapObjects/1024, m.NumGC)
}
