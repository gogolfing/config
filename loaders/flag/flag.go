package flag

import (
	flaglib "flag"
	"os"

	"github.com/gogolfing/config"
)

const DashSeparatorKeyParser = config.SeparatorKeyParser("-")

type Loader struct {
	Args []string

	Aliases map[string]string

	KeyParser config.KeyParser

	LoadDefaults bool

	FlagSet *flaglib.FlagSet
}

func NewLoader(name string) *Loader {
	return &Loader{
		Args:         os.Args[1:],
		Aliases:      map[string]string{},
		KeyParser:    DashSeparatorKeyParser,
		LoadDefaults: false,
		FlagSet:      flaglib.NewFlagSet(name, flaglib.ContinueOnError),
	}
}

func (l *Loader) AddAlias(name, alias string) *Loader {
	l.Aliases[name] = alias
	return l
}

func (l *Loader) Load() (*config.Values, error) {
	if !l.FlagSet.Parsed() {
		args := l.Args
		if args == nil {
			args = os.Args[1:]
		}
		err := l.FlagSet.Parse(args)
		if err != nil {
			if err != flaglib.ErrHelp {
				return nil, err
			}
		}
	}
	v := config.NewValues()
	visit := l.FlagSet.Visit
	if l.LoadDefaults {
		visit = l.FlagSet.VisitAll
	}
	visit(func(f *flaglib.Flag) {
		l.putFlagIntoValues(v, f)
	})
	return v, nil
}

func (l *Loader) putFlagIntoValues(v *config.Values, f *flaglib.Flag) {
	//all flaglib.Values are supposed to implement flaglib.Getters
	getter, ok := f.Value.(flaglib.Getter)
	if !ok {
		return
	}
	name := f.Name
	if alias, ok := l.Aliases[f.Name]; ok {
		name = alias
	}
	key := l.KeyParser.Parse(name)
	v.Put(key, getter.Get())
}
