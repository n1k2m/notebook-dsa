package hashtable

import (
	"notebook-dsa/event"
	"sort"
	"strings"
)

const tableSize = 10007

// node - узел цепочки (метод разрешения коллизий - цепочки)
type node struct {
	data event.Event
	next *node
}

type HashTable struct {
	buckets [tableSize]*node
	count   int
}

func New() *HashTable {
	return &HashTable{}
}

func bucketIndex(e *event.Event) int {
	return int(event.HashKey(e) % tableSize)
}

func (ht *HashTable) Add(e event.Event) {
	ht.Insert(&e)
}

func (ht *HashTable) Insert(e *event.Event) {
	idx := bucketIndex(e)
	ht.buckets[idx] = &node{data: *e, next: ht.buckets[idx]}
	ht.count++
}

type SearchResult struct {
	Events []*event.Event
	Ops    int64
}

func (ht *HashTable) FindByDatetime(y, mo, d, h, mi int) SearchResult {
	key := event.Event{Year: y, Month: mo, Day: d, Hour: h, Minute: mi}
	return ht.Search(&key)
}

func (ht *HashTable) Search(key *event.Event) SearchResult {
	res := SearchResult{}
	idx := bucketIndex(key)
	curr := ht.buckets[idx]
	for curr != nil {
		res.Ops++
		if event.Equal(&curr.data, key) {
			res.Events = append(res.Events, &curr.data)
		}
		curr = curr.next
	}
	return res
}

func (ht *HashTable) Delete(key *event.Event) bool {
	idx := bucketIndex(key)
	curr := ht.buckets[idx]
	var prev *node
	for curr != nil {
		if event.Equal(&curr.data, key) {
			if prev != nil {
				prev.next = curr.next
			} else {
				ht.buckets[idx] = curr.next
			}
			ht.count--
			return true
		}
		prev = curr
		curr = curr.next
	}
	return false
}

func (ht *HashTable) AllSorted() []event.Event {
	events := make([]event.Event, 0, ht.count)
	for i := 0; i < tableSize; i++ {
		curr := ht.buckets[i]
		for curr != nil {
			events = append(events, curr.data)
			curr = curr.next
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return event.Cmp(&events[i], &events[j]) < 0
	})
	return events
}

func (ht *HashTable) FilterPlace(substr string) []event.Event {
	sub := strings.ToLower(substr)
	var result []event.Event
	for i := 0; i < tableSize; i++ {
		curr := ht.buckets[i]
		for curr != nil {
			if strings.Contains(strings.ToLower(curr.data.Place), sub) {
				result = append(result, curr.data)
			}
			curr = curr.next
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return event.Cmp(&result[i], &result[j]) < 0
	})
	return result
}

func (ht *HashTable) Count() int {
	return ht.count
}

func (ht *HashTable) Collect() []event.Event {
	events := make([]event.Event, 0, ht.count)
	for i := 0; i < tableSize; i++ {
		curr := ht.buckets[i]
		for curr != nil {
			events = append(events, curr.data)
			curr = curr.next
		}
	}
	return events
}
