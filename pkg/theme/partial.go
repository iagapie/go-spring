package theme

type (
	Partial interface {
		View
	}

	partial struct {
		View
		t Theme
	}
)

func PartialInTheme(t Theme, name string) Partial {
	if v := t.Datasource().SelectOne("partials", name, "html"); v != nil {
		return &partial{
			View: v,
			t:    t,
		}
	}
	return nil
}
