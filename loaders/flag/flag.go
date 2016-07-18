//Package flag defines a config.Loader type that allows laoding values from
//command line arguments using the language's flag package.
//
//We import the flag package but give it the name flaglib to avoid confusion.
//References to flaglib throughout the documentation refer to the standard "flag" package.
package flag

import (
	flaglib "flag"
	"os"

	"github.com/gogolfing/config"
)

//DashSeparatorKeyParaser is the default KeyParser for a Loader (set in New()).
const DashSeparatorKeyParser = config.SeparatorKeyParser("-")

//Loader provides settings to load values from command line arguments (or any
//slice of strings).
//Loader implements config.Loader.
type Loader struct {
	//Args is the slice of strings sent to flaglib.FlagSet.Parse().
	//It is set to os.Args[1:] by New().
	Args []string

	//Aliases provides a way to change the name of a command line flag before it
	//gets parsed by KeyParser and inserted into the resulting Values instance.
	//If a flaglib.Flag's Name is a key in this map, then the associated value
	//will be used to create the key instead.
	Aliases map[string]string

	//KeyParser is used to turn a flaglib.Flag's Name (or alias from Aliases)
	//into a Key for insertion into the resulting Values.
	KeyParser config.KeyParser

	//LoadDefaults tells Loader to include in the resulting Values all flags that
	//are defined even if they are not set.
	//The default operation is the to only insert flags that have been set explicitly
	//in Args.
	LoadDefaults bool

	//FlagSet is the provided type for managing arguments from the flaglib package.
	//Work with this instance directly to add flags you want parsed and inserted
	//into the resulting Values.
	FlagSet *flaglib.FlagSet
}

//New creates a *Loader with Args set to os.Args[1:],
//Aliases set to a new map,
//KeyParser set to DashSeparatorKeyParser,
//LoadDefaults set to false,
//and FlagSet set to flaglib.NewFlagSet(name, flaglib.ContinueOnError).
func New(name string) *Loader {
	return &Loader{
		Args:         os.Args[1:],
		Aliases:      map[string]string{},
		KeyParser:    DashSeparatorKeyParser,
		LoadDefaults: false,
		FlagSet:      flaglib.NewFlagSet(name, flaglib.ContinueOnError),
	}
}

//AddAlias inserts name, alias into l.Aliases.
func (l *Loader) AddAlias(name, alias string) *Loader {
	l.Aliases[name] = alias
	return l
}

//Load is the config.Loader required method.
//It calls l.FlagSet.Parse(l.Args) if l.FlagSet.Parsed() is false.
//It then calls one of the flaglib.FlagSet.Visit*() methods depending on the value
//of l.LoadDefaults, and parses each flag's Name or alias and inserts it into the
//returned Values.
func (l *Loader) Load() (*config.Values, error) {
	if !l.FlagSet.Parsed() {
		args := l.Args
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
