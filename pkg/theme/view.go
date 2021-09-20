package theme

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/iagapie/go-spring/pkg/helper"
	"gopkg.in/ini.v1"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type (
	View interface {
		ID() string
		Name() string
		CfgProps() CfgProps
		CfgComps() CfgComps
		Cfg() *ini.File
		SetCfg(cfg *ini.File)
		Content() string
		SetContent(content string)
		File() string
		SetFile(file string)
		AddFuncs(funcs FuncMap)
		Parsed() bool
		Load() error
		Save() error
		Execute(w io.Writer, data interface{}) error
		Render(data interface{}) (string, error)
	}

	FuncMap  map[string]interface{}
	CfgProps map[string]string
	CfgComp  struct {
		Name  string
		Alias string
		Props CfgProps
	}
	CfgComps map[string]CfgComp

	ViewOption interface {
		apply(v *view)
	}
)

type (
	view struct {
		id         string
		file       string
		content    string
		delimLeft  string
		delimRight string
		start      string
		end        string
		parsed     bool
		cfg        *ini.File
		funcs      FuncMap
		t          *template.Template
	}

	viewOption func(v *view)
)

const (
	dDelimLeft  = "{{"
	dDelimRight = "}}"
	dStart      = "[cfg]"
	dEnd        = "[/cfg]"
)

func WithViewFile(file string) ViewOption {
	return viewOption(func(v *view) {
		v.SetFile(file)
	})
}

func WithViewContent(content string) ViewOption {
	return viewOption(func(v *view) {
		v.SetContent(content)
	})
}

func WithViewFuncs(funcs FuncMap) ViewOption {
	return viewOption(func(v *view) {
		v.funcs = funcs
	})
}

func WithViewData(data *ini.File) ViewOption {
	return viewOption(func(v *view) {
		v.SetCfg(data)
	})
}

func WithViewDataSep(start, end string) ViewOption {
	return viewOption(func(v *view) {
		v.start = start
		v.end = end
	})
}

func WithViewDelims(left, right string) ViewOption {
	return viewOption(func(v *view) {
		v.delimLeft = left
		v.delimRight = right
	})
}

func NewView(opts ...ViewOption) View {
	id := uuid.NewString()

	v := &view{
		id:         id,
		delimLeft:  dDelimLeft,
		delimRight: dDelimRight,
		start:      dStart,
		end:        dEnd,
		parsed:     false,
		cfg:        ini.Empty(),
		funcs:      make(FuncMap),
	}

	for _, opt := range opts {
		opt.apply(v)
	}

	return v
}

func (v *view) ID() string {
	return v.id
}

func (v *view) Name() string {
	return filepath.Base(v.File())
}

func (v *view) CfgProps() CfgProps {
	return keysToCfgProps(v.cfg.Section("").Keys())
}

func (v *view) CfgComps() CfgComps {
	comps := make(CfgComps)
	for _, s := range v.cfg.Sections() {
		if name := s.Name(); name != ini.DefaultSection {
			alias := name
			if index := strings.IndexRune(name, ' '); index != -1 {
				alias = name[index+1:]
				name = name[:index]
			}
			comps[alias] = CfgComp{
				Name:  name,
				Alias: alias,
				Props: keysToCfgProps(s.Keys()),
			}
		}
	}
	return comps
}

func (v *view) Cfg() *ini.File {
	return v.cfg
}

func (v *view) SetCfg(cfg *ini.File) {
	v.cfg = cfg
}

func (v *view) Content() string {
	return v.content
}

func (v *view) SetContent(content string) {
	v.parsed = false
	v.content = content
}

func (v *view) File() string {
	return v.file
}

func (v *view) SetFile(file string) {
	v.id = file
	v.file = file
}

func (v *view) AddFuncs(funcs FuncMap) {
	for name, fn := range funcs {
		v.funcs[name] = fn
	}
}

func (v *view) Parsed() bool {
	return v.parsed
}

func (v *view) Load() error {
	if err := v.read(); err != nil {
		return err
	}
	if err := v.findCfg(); err != nil {
		return err
	}
	return v.parse()
}

func (v *view) Save() error {
	if v.file == "" {
		return fmt.Errorf("file path is empty, use SetFile")
	}

	data := new(bytes.Buffer)
	_, err := v.cfg.WriteTo(data)
	if err != nil {
		return err
	}

	f, err := os.Create(v.file)
	if err != nil {
		return err
	}
	defer f.Close()

	c := v.content
	if data.Len() > 0 {
		c = fmt.Sprintf("%s\n%s\n%s\n%s", v.start, data.String(), v.end, c)
	}

	_, err = f.WriteString(c)
	return err
}

func (v *view) Execute(w io.Writer, data interface{}) error {
	if err := v.parse(); err != nil {
		return err
	}
	return v.t.Execute(w, data)
}

func (v *view) Render(data interface{}) (string, error) {
	w := new(bytes.Buffer)
	if err := v.Execute(w, data); err != nil {
		return "", err
	}
	return w.String(), nil
}

func (v *view) read() error {
	if helper.FileExists(v.file) {
		b, err := os.ReadFile(v.file)
		v.content = string(b)
		return err
	}
	return nil
}

func (v *view) parse() error {
	if !v.parsed {
		v.parsed = true
		v.t = template.New(v.id).Funcs(template.FuncMap(v.funcs))
		_, err := v.t.Parse(v.content)
		return err
	}
	return nil
}

func (v *view) findCfg() error {
	cfg := ""
	for {
		start := strings.Index(v.content, v.start)
		if start == -1 {
			return v.parseCfg(cfg)
		}

		end := strings.Index(v.content, v.end)
		if end == -1 {
			return fmt.Errorf("%s not found in %s", v.end, v.content)
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
	v.cfg = f
	return nil
}

func (fn viewOption) apply(v *view) {
	fn(v)
}

func keysToCfgProps(keys []*ini.Key) CfgProps {
	props := make(CfgProps)
	for _, key := range keys {
		props[key.Name()] = key.Value()
	}
	return props
}
