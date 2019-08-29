package tex

import (
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"golang.org/x/tools/imports"
	"golang.org/x/xerrors"
	"path/filepath"
	"testing"
)

var (
	sourceDir = "/templates"
	destDir   = "/out"
	templates = map[string]string{
		"example.go.tmpl": `package main

func main() {
	fmt.Println("Hello {{.Message}}")
}
`,
		"bad.go.tmpl": `package main

func main() {
	fmt.Println("Hello {{.Message}}"
}
`,
		"hello.txt": `Hello {{.Message}}`,
	}
)

func fs() (afero.Fs, error) {
	fs := afero.NewMemMapFs()

	if err := fs.Mkdir(sourceDir, 0644); nil != err {
		return nil, err
	}

	if err := fs.Mkdir(destDir, 0644); nil != err {
		return nil, err
	}

	for name, content := range templates {
		if err := afero.WriteFile(fs, filepath.Join(sourceDir, name), []byte(content), 0755); nil != err {
			return nil, err
		}
	}

	return fs, nil
}
func mustFs() afero.Fs {
	fs, err := fs()
	if nil != err {
		panic(err)
	}
	return fs
}

func TestRenderInvalidOptions(t *testing.T) {

	type args struct{
		options Options
	}

	var tests = []struct {
		name    string
		args args
		errType error
		init func(args *args) error
	}{
		{
			name: "empty source",
			args: args{
				options: Options{
					Source: "",
					Dest:   destDir,
					Fs:     mustFs(),
				},
			},
			errType: &ErrInvalidSource{},
		},
		{
			name: "empty dest",
			args: args{
				options: Options{
					Source: filepath.Join(sourceDir, "*.tmpl"),
					Dest:   "",
					Fs:     mustFs(),
				},
			},
			errType: &ErrInvalidDest{},
		},
		{
			name: "no source match",
			args: args{
				options: Options{
					Source: filepath.Join(sourceDir, "*.tpl"),
					Dest:   destDir,
					Fs:     mustFs(),
				},
			},
			errType: &ErrNoNoMatch{},
		},
		{
			name: "dest not dir",
			args: args{
				options: Options{
					Source: filepath.Join(sourceDir, "*.tmpl"),
					Dest:   destDir,
					Fs:     mustFs(),
				},
			},
			errType: &ErrInvalidDest{},
			init: func(args *args) error {
				if err := args.options.Fs.Remove(destDir); nil != err {
					return err
				}
				_, err := args.options.Fs.Create(destDir)
				return err
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			g := NewGomegaWithT(tt)

			if nil != tc.init {
				err := tc.init(&tc.args)
				g.Expect(err).ToNot(HaveOccurred())
			}

			err := Render(tc.args.options)
			g.Expect(err).To(HaveOccurred())
			g.Expect(xerrors.Is(err, tc.errType)).To(BeTrue())
		})
	}
}

func TestRender(t *testing.T) {

	g := NewGomegaWithT(t)

	var mustProcess = func(str string) string {
		bs, err := imports.Process(filepath.Join(destDir, sourceDir, "example.go"), []byte(str), nil)
		g.Expect(err).ToNot(HaveOccurred())
		return string(bs)
	}

	type data struct {
		Message string
	}

	var tests = []struct {
		name         string
		options      Options
		expectedPath string
		expected     string
	}{
		{
			name: "example.go.tmpl",
			options: Options{
				Source:  filepath.Join(sourceDir, "example.go.tmpl"),
				Dest:    destDir,
				WorkDir: "/",
				Fs:      mustFs(),
				Data:    data{Message: "tex"},
				Funcs:   nil,
			},
			expectedPath: filepath.Join(destDir, sourceDir, "example.go"),
			expected: mustProcess(`package main

func main() {
	fmt.Println("Hello tex")
}`),
		},
		{
			name: "bad.go.tmpl",
			options: Options{
				Source:  filepath.Join(sourceDir, "bad.go.tmpl"),
				Dest:    destDir,
				WorkDir: "/",
				Fs:      mustFs(),
				Data:    data{Message: "tex"},
				Funcs:   nil,
			},
			expectedPath: filepath.Join(destDir, sourceDir, "bad.go"),
			expected: `package main

func main() {
	fmt.Println("Hello tex"
}
`,
		},
		{
			name: "hello.txt",
			options: Options{
				Source:  filepath.Join(sourceDir, "hello.txt"),
				Dest:    destDir,
				WorkDir: "/",
				Fs:      mustFs(),
				Data:    data{Message: "tex"},
				Funcs:   nil,
			},
			expectedPath: filepath.Join(destDir, sourceDir, "hello.txt"),
			expected:     `Hello tex`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			g := NewGomegaWithT(tt)

			err := Render(tc.options)
			g.Expect(err).ToNot(HaveOccurred())

			bs, err := afero.ReadFile(tc.options.Fs, tc.expectedPath)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(string(bs)).To(Equal(tc.expected))
		})
	}
}
