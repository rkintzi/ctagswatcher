package globfilter

import "path/filepath"

type Action bool

const (
	Deny  Action = false
	Allow        = true
)

type rule struct {
	pat string
	act Action
}

type Filter struct {
	rules []rule
	act   Action
}

func NewFilter(act Action) *Filter {
	return &Filter{act: act}
}

func (f *Filter) Append(act Action, patterns ...string) error {
	rules := make([]rule, len(patterns))
	for i, pat := range patterns {
		_, err := filepath.Match(pat, "")
		if err != nil {
			return err
		}
		rules[i].pat = pat
		rules[i].act = act
	}
	f.rules = append(f.rules, rules...)
	return nil
}

func (f *Filter) Filter(path string) bool {
	for _, rule := range f.rules {
		if r, _ := filepath.Match(rule.pat, path); r {
			return bool(rule.act)
		}
	}
	return bool(f.act)
}
