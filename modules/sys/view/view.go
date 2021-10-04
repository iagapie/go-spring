package view

import (
	"bytes"
	"fmt"
	"github.com/iagapie/go-spring/modules/sys/helper"
	"gopkg.in/ini.v1"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type (
	Props   map[string]string
	Comps   map[string]*Comp
	Comp    struct {
		Name  string
		Alias string
		Props Props
	}

	View interface {
		File() string
		Name() string
		Exists() bool
		Props() Props
		Prop(name string) string
		CfgComps() Comps
		Content() string
		Funcs(funcMap template.FuncMap)
		Load() error
		Execute(w io.Writer, vars interface{}) error
		Render(vars interface{}) (string, error)
	}

	Option interface {
		apply(v *view)
	}

	option func(v *view)

	view struct {
		mu         sync.Mutex
		file       string
		content    string
		props      Props
		comps      Comps
		delimLeft  string
		delimRight string
		start      string
		end        string
		funcs      template.FuncMap
		t          *template.Template
	}
)

const (
	cDelimLeft  = "{{"
	cDelimRight = "}}"
	cStart      = "[cfg]"
	cEnd        = "[/cfg]"
)

func WithFuncs(funcMap template.FuncMap) Option {
	return option(func(v *view) {
		v.Funcs(funcMap)
	})
}

func WithCfgSep(start, end string) Option {
	return option(func(v *view) {
		v.start = start
		v.end = end
	})
}

func WithDelims(left, right string) Option {
	return option(func(v *view) {
		v.delimLeft = left
		v.delimRight = right
	})
}

func New(file string, opts ...Option) View {
	v := &view{
		file:       file,
		delimLeft:  cDelimLeft,
		delimRight: cDelimRight,
		start:      cStart,
		end:        cEnd,
		props:      make(Props),
		comps:      make(Comps),
		funcs:      globalFuncs(),
	}

	for _, opt := range opts {
		opt.apply(v)
	}

	return v
}

func (v *view) File() string {
	return v.file
}

func (v *view) Name() string {
	return filepath.Base(v.File())
}

func (v *view) Exists() bool {
	return helper.FileExists(v.File())
}

func (v *view) Props() Props {
	v.mu.Lock()
	defer v.mu.Unlock()
	return cpProps(v.props)
}

func (v *view) Prop(name string) string {
	v.mu.Lock()
	defer v.mu.Unlock()
	if value, ok := v.props[name]; ok {
		return value
	}
	return ""
}

func (v *view) CfgComps() Comps {
	v.mu.Lock()
	defer v.mu.Unlock()
	return cpComps(v.comps)
}

func (v *view) Content() string {
	return v.content
}

func (v *view) Funcs(funcMap template.FuncMap) {
	for name, fn := range funcMap {
		v.funcs[name] = fn
	}
}

func (v *view) Load() error {
	if v.t != nil {
		return nil
	}
	if err := v.read(); err != nil {
		return err
	}
	if err := v.cfg(); err != nil {
		return err
	}
	return v.parse()
}

func (v *view) Execute(w io.Writer, vars interface{}) error {
	if err := v.Load(); err != nil {
		return err
	}
	return v.t.Execute(w, vars)
}

func (v *view) Render(vars interface{}) (string, error) {
	w := new(bytes.Buffer)
	if err := v.Execute(w, vars); err != nil {
		return "", err
	}
	return w.String(), nil
}

func (v *view) read() error {
	if !v.Exists() {
		return fmt.Errorf("view %s not found", v.File())
	}
	b, err := os.ReadFile(v.File())
	v.content = string(b)
	return err
}

func (v *view) cfg() error {
	cfg := ""
	for {
		start := strings.Index(v.content, v.start)
		if start == -1 {
			return v.parseCfg(cfg)
		}

		end := strings.Index(v.content, v.end)
		if end == -1 {
			return fmt.Errorf("%s not found in %s", v.end, v.File())
		}

		cfg += strings.TrimSpace(v.content[start+len(v.start) : end])
		v.content = strings.TrimSpace(v.content[:start] + v.content[end+len(v.end):])
	}
}

func (v *view) parseCfg(cfg string) error {
	f, err := ini.Load([]byte(cfg))
	if err != nil {
		return err
	}
	for _, s := range f.Sections() {
		if name := s.Name(); name != ini.DefaultSection {
			alias := name
			if index := strings.IndexRune(name, ' '); index != -1 {
				alias = name[index+1:]
				name = name[:index]
			}
			v.comps[alias] = &Comp{
				Name:  name,
				Alias: alias,
				Props: keysToCfgProps(s.Keys()),
			}
		} else {
			v.props = keysToCfgProps(s.Keys())
		}
	}
	return nil
}

func (v *view) parse() (err error) {
	v.t, err = template.New(v.File()).Funcs(v.funcs).Parse(v.Content())
	return
}

func (fn option) apply(v *view) {
	fn(v)
}

func keysToCfgProps(keys []*ini.Key) Props {
	props := make(Props)
	for _, key := range keys {
		props[key.Name()] = key.Value()
	}
	return props
}

func cpProps(props Props) Props {
	cp := make(Props)
	for k, v := range props {
		cp[k] = v
	}
	return cp
}

func cpComp(comp *Comp) *Comp {
	return &Comp{
		Name:  comp.Name,
		Alias: comp.Alias,
		Props: cpProps(comp.Props),
	}
}

func cpComps(comps Comps) Comps {
	cp := make(Comps)
	for k, v := range comps {
		cp[k] = cpComp(v)
	}
	return cp
}
