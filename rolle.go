package rollee

import "sync"

type ID = int

type List struct {
	ID     ID
	Values []int
}

func Fold(initialValue int, f func(int, int) int, l List) map[ID]int {
	for _, val := range l.Values {
		initialValue = f(initialValue, val)
	}

	return map[ID]int{l.ID: initialValue}
}

func FoldChan(initialValue int, f func(int, int) int, ch chan List) map[ID]int {
	result := make(map[ID]int)

	for c := range ch {
		acc, exists := result[c.ID]
		if !exists {
			acc = initialValue
		}
		for _, val := range c.Values {
			acc = f(acc, val)
		}
		result[c.ID] = acc
	}

	return result
}

func FoldChanX(initialValue int, f func(int, int) int, chs ...chan List) map[ID]int {
	result := make(map[ID]int)
	var wg sync.WaitGroup
	mu := sync.Mutex{}

	processCh := func(ch chan List) {
		defer wg.Done()
		localResult := FoldChan(initialValue, f, ch)

		mu.Lock()
		defer mu.Unlock()

		for id, value := range localResult {
			if existingValue, exists := result[id]; exists {
				result[id] = f(existingValue, value)
			} else {
				result[id] = value
			}
		}
	}

	wg.Add(len(chs))
	for _, ch := range chs {
		go processCh(ch)
	}
	wg.Wait()

	return result
}
