package tex

import (
	"bytes"
	"github.com/spf13/afero"
	"golang.org/x/tools/imports"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"text/template"
)

var extRegex = regexp.MustCompile(`(\.[^.]+)\.tm?pl$`)

type Options struct {
	Source  string
	Dest    string
	WorkDir string

	Fs afero.Fs

	Data  interface{}
	Funcs template.FuncMap

	Logger *log.Logger
}

func (o *Options) ApplyDefaults() {
	if "" == o.WorkDir {
		_, o.WorkDir, _, _ = runtime.Caller(1)
	}

	if nil == o.Fs {
		o.Fs = afero.NewOsFs()
	}

	funcs := DefaultFuncMap()
	for n, f := range o.Funcs {
		funcs[n] = f
	}
	o.Funcs = funcs
}

func (o *Options) Validate() error {
	if "" == o.Source {
		return &ErrInvalidSource{o.Source}
	}
	if "" == o.Dest {
		return &ErrInvalidDest{o.Dest}
	}
	return nil
}

func Render(opts Options) error {
	opts.ApplyDefaults()
	if err := opts.Validate(); nil != err {
		return err
	}
	sources, err := afero.Glob(opts.Fs, opts.Source)
	if nil != err {
		return err
	}
	if 0 == len(sources) {
		return &ErrNoNoMatch{Source: opts.Source}
	}

	destStat, err := opts.Fs.Stat(opts.Dest)
	if nil != err {
		if os.IsNotExist(err) {
			return wrapErr(&ErrInvalidDest{Dest: opts.Dest}, "not exists")
		}
		return err
	}
	if !destStat.IsDir() {
		return wrapErr(&ErrInvalidDest{Dest: opts.Dest}, "not a directory")
	}

	for _, source := range sources {
		bs, err := afero.ReadFile(opts.Fs, source)
		if nil != err {
			return err
		}

		relPath, err := filepath.Rel(opts.WorkDir, source)
		if nil != err {
			return err
		}

		tmpl, err := template.New(relPath).Funcs(opts.Funcs).Parse(string(bs))
		if nil != err {
			return err
		}

		extMatch := extRegex.FindStringSubmatch(relPath)
		var outFile string
		if 2 == len(extMatch) {
			outFile = filepath.Join(opts.Dest, extRegex.ReplaceAllString(relPath, extMatch[1]))
		} else {
			outFile = filepath.Join(opts.Dest, relPath)
		}

		outDir := filepath.Dir(outFile)
		outDirStat, err := opts.Fs.Stat(outDir)
		if nil != err {
			if !os.IsNotExist(err) {
				return err
			}
			if err := opts.Fs.MkdirAll(outDir, 0644); nil != err {
				return err
			}
		} else {
			if !outDirStat.IsDir() {
				return wrapErr(&ErrInvalidDest{Dest: opts.Dest}, "not a directory")
			}
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, opts.Data); nil != err {
			return err
		}

		bs, err = imports.Process(outFile, buf.Bytes(), nil)
		if nil != err {
			if nil != opts.Logger {
				opts.Logger.Printf("failed to process file: %s", err.Error())
			}
			if err := afero.WriteFile(opts.Fs, outFile, buf.Bytes(), 0755); nil != err {
				return err
			}
		} else {
			if err := afero.WriteFile(opts.Fs, outFile, bs, 0755); nil != err {
				return err
			}
		}

	}
	return nil
}
