package component

type (
	psItem map[string][]map[string]interface{}

	PartialStack struct {
		activePartial psItem
		partialStack  []psItem
	}
)

func NewPartialStack() *PartialStack {
	return &PartialStack{
		partialStack: make([]psItem, 0),
	}
}

func (ps *PartialStack) StackPartial() {
	if ps.activePartial != nil {
		ps.partialStack = append([]psItem{ps.activePartial}, ps.partialStack...)
	}
	ps.activePartial = psItem{
		"components": make([]map[string]interface{}, 0),
	}
}

func (ps *PartialStack) UnstackPartial() {
	if len(ps.partialStack) > 0 {
		ps.activePartial = ps.partialStack[0]
		ps.partialStack = ps.partialStack[1:]
	} else {
		ps.activePartial = nil
	}
}

func (ps *PartialStack) AddComponent(alias string, comp Component) {
	ps.activePartial["components"] = append(ps.activePartial["components"], map[string]interface{}{
		"name": alias,
		"obj":  comp,
	})
}

func (ps *PartialStack) Component(name string) Component {
	if ps.activePartial == nil {
		return nil
	}
	if comp := ps.findComponentFromStack(name, ps.activePartial); comp != nil {
		return comp
	}
	for _, stack := range ps.partialStack {
		if comp := ps.findComponentFromStack(name, stack); comp != nil {
			return comp
		}
	}
	return nil
}

func (ps *PartialStack) findComponentFromStack(name string, stack psItem) Component {
	for _, info := range stack["components"] {
		if info["name"].(string) == name {
			return info["obj"].(Component)
		}
	}
	return nil
}
