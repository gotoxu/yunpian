package assert

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Fataler 定义了当达到给定条件时触发致命错误的最小接口
// testing.T 或者 testing.B 都满足此接口
type Fataler interface {
	Fatal(a ...interface{})
}

type cond struct {
	Fataler           Fataler
	Skip              int
	Format            string
	FormatArgs        []interface{}
	Extra             []interface{}
	DisableDeleteSelf bool
}

var deleteSelf = strings.Repeat("\b", 15)

func (c cond) String() string {
	var b bytes.Buffer
	if c.DisableDeleteSelf {
		fmt.Fprint(&b, "\n")
	} else {
		fmt.Fprint(&b, deleteSelf)
	}

	fmt.Fprintf(&b, pstack(callers(c.Skip+1), c.DisableDeleteSelf))
	if c.Format != "" {
		fmt.Fprintf(&b, c.Format, c.FormatArgs...)
	}
	if len(c.Extra) != 0 {
		fmt.Fprintf(&b, "\n")
		fmt.Fprintf(&b, tsdump(c.Extra...))
	}

	return b.String()
}

// IsType 判断给定对象的类型是否和期望值相同
func IsType(t Fataler, object interface{}, expected interface{}, a ...interface{}) {
	if !reflect.DeepEqual(reflect.TypeOf(object), reflect.TypeOf(expected)) {
		fatal(cond{
			Fataler:    t,
			Format:     "Object expected to be of type %v, but was %v",
			FormatArgs: []interface{}{reflect.TypeOf(expected), reflect.TypeOf(object)},
			Extra:      a,
		})
	}
}

// Len 判断指定的对象是否具有指定的长度
func Len(t Fataler, object interface{}, length int, a ...interface{}) {
	ok, l := getLen(object)
	if !ok {
		fatal(cond{
			Fataler:    t,
			Format:     `"%s" could not be applied builtin len()`,
			FormatArgs: []interface{}{object},
			Extra:      a,
		})
	}

	if l != length {
		fatal(cond{
			Fataler:    t,
			Format:     `"%s" should have %d item(s), but has %d`,
			FormatArgs: []interface{}{object, length, l},
			Extra:      a,
		})
	}
}

func getLen(x interface{}) (ok bool, length int) {
	v := reflect.ValueOf(x)
	defer func() {
		if e := recover(); e != nil {
			ok = false
		}
	}()

	return true, v.Len()
}

// Empty 判断给定的对象v是否为空，v必须可作用于内建的len函数
func Empty(t Fataler, v interface{}, a ...interface{}) {
	ok, l := getLen(v)
	if !ok {
		fatal(cond{
			Fataler:    t,
			Format:     `"%s" could not be applied builtin len()`,
			FormatArgs: []interface{}{v},
			Extra:      a,
		})
		return
	}

	if l != 0 {
		fatal(cond{
			Fataler:    t,
			Format:     "expected an empty object, but got the: %s",
			FormatArgs: []interface{}{tsdump(v)},
			Extra:      a,
		})
	}
}

// NotEmpty 判断给定的集合是否非空, v必须可以作用于内建的len函数
func NotEmpty(t Fataler, v interface{}, a ...interface{}) {
	ok, l := getLen(v)
	if !ok {
		fatal(cond{
			Fataler:    t,
			Format:     `"%s" could not be applied builtin len()`,
			FormatArgs: []interface{}{v},
			Extra:      a,
		})
	}

	if l == 0 {
		fatal(cond{
			Fataler: t,
			Format:  "expected an not empty object, but got an empty one.",
			Extra:   a,
		})
	}
}

// DeepEqual 判断给定的2个对象是否相等，使用reflect.DeepEqual进行判断
func DeepEqual(t Fataler, actual, expected interface{}, a ...interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		fatal(cond{
			Fataler:    t,
			Format:     "expected these to be equal:\nACTUAL:\n%s\nEXPECTED:\n%s",
			FormatArgs: []interface{}{spew.Sdump(actual), tsdump(expected)},
			Extra:      a,
		})
	}
}

// NotDeepEqual 判断给定的2个对象是否不同，使用reflect.DeepEqual进行判断
func NotDeepEqual(t Fataler, actual, expected interface{}, a ...interface{}) {
	if reflect.DeepEqual(actual, expected) {
		fatal(cond{
			Fataler:    t,
			Format:     "expected two different values, but got the same:\n%s",
			FormatArgs: []interface{}{tsdump(actual)},
			Extra:      a,
		})
	}
}

// Error 判断给定的错误信息是否与正则表达式匹配
func Error(t Fataler, err error, re *regexp.Regexp, a ...interface{}) {
	if err == nil && re == nil {
		return
	}

	if err == nil && re != nil {
		fatal(cond{
			Fataler:    t,
			Format:     `expected error: "%s" but got a nil error`,
			FormatArgs: []interface{}{re},
			Extra:      a,
		})
		return
	}

	if err != nil && re == nil {
		fatal(cond{
			Fataler:    t,
			Format:     "unexpected error: %s",
			FormatArgs: []interface{}{err},
			Extra:      a,
		})
		return
	}

	if !re.MatchString(err.Error()) {
		fatal(cond{
			Fataler:    t,
			Format:     `expected error: "%s" but got "%s"`,
			FormatArgs: []interface{}{re, err},
			Extra:      a,
		})
	}
}

// StringContains 判断给定的字符串是否包含字串
func StringContains(t Fataler, s, substr string, a ...interface{}) {
	if !strings.Contains(s, substr) {
		format := `expected substring "%s" was not found in "%s"`

		if strings.Contains(s, "\n") || strings.Contains(substr, "\n") {
			format = `expected substring was not found:\nEXPECTED SUBSTRING:\n%s\nACTUAL:\n%s`
		}

		fatal(cond{
			Fataler:    t,
			Format:     format,
			FormatArgs: []interface{}{substr, s},
			Extra:      a,
		})
	}
}

// StringDoesNotContain 判断给定的字符串是否不包含字串
func StringDoesNotContain(t Fataler, s, substr string, a ...interface{}) {
	if strings.Contains(s, substr) {
		fatal(cond{
			Fataler:    t,
			Format:     `substring "%s" was not supposed to be found in "%s"`,
			FormatArgs: []interface{}{substr, s},
			Extra:      a,
		})
	}
}

// Nil 方法判断 v 是否等于 nil
func Nil(t Fataler, v interface{}, a ...interface{}) {
	vs := tsdump(v)
	sp := " "
	if strings.Contains(vs[:len(vs)-1], "\n") {
		sp = "\n"
	}

	if v != nil {
		if _, ok := v.(error); ok {
			fatal(cond{
				Fataler:    t,
				Format:     `unexpected error: %s`,
				FormatArgs: []interface{}{v},
				Extra:      a,
			})
		} else {
			fatal(cond{
				Fataler:    t,
				Format:     "expected nil value but got:%s%s",
				FormatArgs: []interface{}{sp, vs},
				Extra:      a,
			})
		}
	}
}

// NotNil 判断 v 是否不等于 nil
func NotNil(t Fataler, v interface{}, a ...interface{}) {
	if v == nil {
		fatal(cond{
			Fataler: t,
			Format:  "expected a value but got nil",
			Extra:   a,
		})
	}
}

// True 判断 v 是否为 true
func True(t Fataler, v bool, a ...interface{}) {
	if !v {
		fatal(cond{
			Fataler: t,
			Format:  "expected true but got false",
			Extra:   a,
		})
	}
}

// False 判断 v 是否为 false
func False(t Fataler, v bool, a ...interface{}) {
	if v {
		fatal(cond{
			Fataler: t,
			Format:  "expected false but got true",
			Extra:   a,
		})
	}
}

func fatal(c cond) {
	c.Skip = c.Skip + 2
	c.Fataler.Fatal(c.String())
}

func tsdump(a ...interface{}) string {
	return strings.TrimSpace(spew.Sdump(a...))
}

func pstack(s stack, skipPrefix bool) string {
	first := s[0]
	if isTestFrame(first) {
		return fmt.Sprintf("%s:%d: ", filepath.Base(first.File), first.Line)
	}

	prefix := "        "
	if skipPrefix {
		prefix = ""
	}

	var snew stack
	for _, f := range s {
		snew = append(snew, f)
		if isTestFrame(f) {
			return prefix + snew.String() + "\n"
		}
	}

	return prefix + s.String() + "\n"
}

func isTestFrame(f frame) bool {
	return strings.HasPrefix(f.Name, "Test")
}
