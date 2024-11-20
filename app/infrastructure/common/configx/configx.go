package configx

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"dario.cat/mergo"
	"github.com/go-playground/validator"
	"github.com/ilyakaznacheev/cleanenv"
)

// Loader represents a config loader
type Loader struct {
	defaultFileName string
	fileFlag        string
	fileName        string
}

// WithDefaultFileName allows you to specify a custom default file name
// Default file name will be check in case no other file name is specified
func (l *Loader) WithDefaultFileName(flnm string) *Loader {
	l.defaultFileName = flnm

	return l
}

// WithFileFlag allows you to specify a custom file flag from where
// file name will be read. This file name will be used in case it is spcified the flag
// before to backup into the default one
// Take into consideration that going for flag will disable filename option
func (l *Loader) WithFileFlag(f string) *Loader {
	l.fileFlag = f
	l.fileName = ""

	return l
}

// WithFileName allows you to specify an specific file name
// Take into consideration that going for filename will disable flag option
func (l *Loader) WithFileName(flnm string) *Loader {
	l.fileName = flnm
	l.fileFlag = ""

	return l
}

// OnlyEnvironment makes loader only read from environment variables
func (l *Loader) OnlyEnvironment() *Loader {
	l.fileFlag = ""
	l.fileName = ""

	return l
}

// Load will fill the specified configuration structure from the source
func (l Loader) Load(cfg interface{}) error {
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return fmt.Errorf("configuration has to be a pointer to a struct but got %T", cfg)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return fmt.Errorf("error reading environment variables: %v", err)
	}

	// Copying structure to parse file configuration without affecting
	// environment variables in case there are some defined
	fileCfg := reflect.New(reflect.ValueOf(cfg).Elem().Type()).Interface()
	fileName := &l.fileName

	// Puts priority in case file name was specified instead of going for the flag
	// If both are empty we do nothing...
	if *fileName == "" && l.fileFlag != "" {
		f := flag.Lookup(l.fileFlag)

		if f == nil {
			fileName = flag.String(l.fileFlag, "", "Specify configuration file")
			flag.Parse()
		} else {
			*fileName = f.Value.String()
		}

		flag.Parse()

		// Only will set to read default config file if loader was configured to expect file flag
		// but the application was executed without specifying it and backup config file really exists
		if _, err := os.Stat(l.defaultFileName); *fileName == "" && err == nil {
			*fileName = l.defaultFileName
		}
	}

	// If it ends up in no file name defined neither from config flag, specific file name
	// nor backup file we don't try to read anything
	if *fileName != "" {
		if err := cleanenv.ReadConfig(*fileName, fileCfg); err != nil {
			return fmt.Errorf("error reading configuration from file: %v", err)
		}
	}

	// This will always put priority into the configuration comming from env variables
	if err := mergo.MergeWithOverwrite(cfg, fileCfg); err != nil {
		return fmt.Errorf("unexpected error merging configuration (%v)", err)
	}

	// Validating in case it was defined within the structure
	if err := validator.New().Struct(cfg); err != nil {
		return fmt.Errorf("configuration %v", err)
	}

	return nil
}

// NewLoader generates a default config loader
func NewLoader() *Loader {
	return &Loader{
		defaultFileName: ".env",
		fileFlag:        "config",
		fileName:        "",
	}
}
