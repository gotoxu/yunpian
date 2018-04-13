package assert

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const maxStackSize = 32

var (
	gopath  string
	gopaths []string
)

func init() {
	if gopath == "" {
		gopath = os.Getenv("GOPATH")
	}
	setGOPATH(gopath)
}

func setGOPATH(gp string) {
	gopath = gp
	gopaths = nil

	if runtime.GOOS == "windows" {
		for _, p := range strings.Split(gopath, ";") {
			gopaths = append(gopaths, filepath.Join(p, "src")+"/")
		}
	} else {
		for _, p := range strings.Split(gopath, ":") {
			if p != "" {
				gopaths = append(gopaths, filepath.Join(p, "src")+"/")
			}
		}
	}

	gopaths = append(gopaths, filepath.Join(runtime.GOROOT(), "src", "pkg")+"/")
}

func stripGOPATH(f string) string {
	for _, p := range gopaths {
		if strings.HasPrefix(f, p) {
			return f[len(p):]
		}
	}
	return f
}

func stripPackage(n string) string {
	slashI := strings.LastIndex(n, "/")
	if slashI == -1 {
		slashI = 0
	}
	dotI := strings.Index(n[slashI:], ".")
	if dotI == -1 {
		return n
	}
	return n[slashI+dotI+1:]
}

// frame 包含了调用栈的文件名，行号和函数名
type frame struct {
	File string
	Line int
	Name string
}

func (f frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Name)
}

type stack []frame

func (s stack) String() string {
	var b bytes.Buffer
	writeStack(&b, s)
	return b.String()
}

type multi struct {
	stacks []stack
}

func (m *multi) Stacks() []stack {
	return m.stacks
}

func (m *multi) Add(s stack) {
	m.stacks = append(m.stacks, s)
}

func (m *multi) AddCallers(skip int) {
	m.Add(callers(skip + 1))
}

func (m *multi) String() string {
	var b bytes.Buffer
	for i, s := range m.stacks {
		if i != 0 {
			fmt.Fprintf(&b, "\n(Stack %d)\n", i+1)
		}
		writeStack(&b, s)
	}
	return b.String()
}

func (m *multi) Copy() *multi {
	m2 := &multi{
		stacks: make([]stack, len(m.stacks)),
	}
	copy(m2.stacks, m.stacks)
	return m2
}

func caller(skip int) frame {
	pc, file, line, _ := runtime.Caller(skip + 1)
	fun := runtime.FuncForPC(pc)
	return frame{
		File: stripGOPATH(file),
		Line: line,
		Name: stripPackage(fun.Name()),
	}
}

func callers(skip int) stack {
	pcs := make([]uintptr, maxStackSize)
	num := runtime.Callers(skip+2, pcs)
	stack := make(stack, num)
	for i, pc := range pcs[:num] {
		fun := runtime.FuncForPC(pc)
		file, line := fun.FileLine(pc - 1)
		stack[i].File = stripGOPATH(file)
		stack[i].Line = line
		stack[i].Name = stripPackage(fun.Name())
	}
	return stack
}

func callersMulti(skip int) *multi {
	m := new(multi)
	m.AddCallers(skip + 1)
	return m
}

func writeStack(b *bytes.Buffer, s stack) {
	var width int
	for _, f := range s {
		if l := len(f.File) + numDigits(f.Line) + 1; l > width {
			width = l
		}
	}
	last := len(s) - 1
	for i, f := range s {
		b.WriteString(f.File)
		b.WriteRune(rune(':'))
		n, _ := fmt.Fprintf(b, "%d", f.Line)
		for i := width - len(f.File) - n; i != 0; i-- {
			b.WriteRune(rune(' '))
		}
		b.WriteString(f.Name)
		if i != last {
			b.WriteRune(rune('\n'))
		}
	}
}

func numDigits(i int) int {
	var n int
	for {
		n++
		i = i / 10
		if i == 0 {
			return n
		}
	}
}
