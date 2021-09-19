package dispatcher

import "sort"

type (
	Arg interface{}

	Listener func(args ...Arg) (interface{}, error)

	Dispatcher interface {
		Listen(event string, priority int, listener Listener)
		Firing() string
		Until(event string, args ...Arg) (interface{}, error)
		Fire(event string, halt bool, args ...Arg) ([]interface{}, error)
		Listeners(event string) []Listener
	}

	dispatcher struct {
		listeners map[string]map[int][]Listener
		sorted    map[string][]Listener
		firing    []string
	}
)

func New() *dispatcher {
	return &dispatcher{
		listeners: make(map[string]map[int][]Listener),
		sorted:    make(map[string][]Listener),
		firing:    make([]string, 0),
	}
}

func (d *dispatcher) Listen(event string, priority int, listener Listener) {
	d.listeners[event][priority] = append(d.listeners[event][priority], listener)
	delete(d.sorted, event)
}

func (d *dispatcher) Firing() string {
	if len(d.firing) > 0 {
		return d.firing[len(d.firing)-1]
	}
	return ""
}

func (d *dispatcher) Until(event string, args ...Arg) (interface{}, error) {
	return d.Fire(event, true, args...)
}

func (d *dispatcher) Fire(event string, halt bool, args ...Arg) ([]interface{}, error) {
	var responses []interface{}

	d.firingPush(event)

	for _, listener := range d.Listeners(event) {
		response, err := listener(args...)
		if err != nil {
			d.firingPop()
			return responses, err
		}

		if response != nil && halt {
			d.firingPop()
			return []interface{}{response}, nil
		}

		if r, ok := response.(bool); ok && !r {
			break
		}

		responses = append(responses, response)
	}

	d.firingPop()

	if halt {
		return nil, nil
	}

	return responses, nil
}

func (d *dispatcher) Listeners(event string) []Listener {
	if sorted, ok := d.sorted[event]; ok {
		return sorted
	}

	d.sorted[event] = make([]Listener, 0)

	if notSorted, ok := d.listeners[event]; ok {
		keys := make([]int, 0, len(notSorted))
		for k := range notSorted {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
		for _, k := range keys {
			d.sorted[event] = append(d.sorted[event], notSorted[k]...)
		}
	}

	return d.sorted[event]
}

func (d *dispatcher) firingPush(event string) {
	d.firing = append(d.firing, event)
}

func (d *dispatcher) firingPop() {
	d.firing = d.firing[:len(d.firing)-1]
}
