# Assert [![Build Status](https://travis-ci.org/gotoxu/assert.svg?branch=master)](https://travis-ci.org/gotoxu/assert)
Assert 包为Go的单元测试提供了一些工具函数，让你可以写出更加简洁的测试代码。该项目的开发基于facebook的开源项目[ensure](https://github.com/facebookgo/ensure)，并在其之上新增了更多的断言函数。

## Installation
To install assert, use `go get`

```go
go get -u github.com/gotoxu/assert
```

Import the `assert` package into your code using this template

```go
package yours

import (
	"testing"
	"github.com/gotoxu/assert"
)

func TestSomething(t *testing.T) {
	assert.True(t, true)
}
```

## Example usages

**Empty**: 该函数判断一个给定的切片长度是否等于0

```go
func TestSliceEmpty(t *testing.T) {
	arr := []int{}
	assert.Empty(t, arr)
}
```

**NotEmpty**: 该函数判断一个给定的切片长度是否不等于0

```go
func TestSliceNotEmpty(t *testing.T) {
	arr := []int{1, 2, 3}
	assert.NotEmpty(t, arr)
}
```

**Nil**: 该函数判断一个给定的对象是否等于nil

```go
func TestNil(t *testing.T) {
	var e error
	assert.Nil(t, e)
}
```

**NotNil**: 该函数判断一个给定的对象是否不等于nil

```go
func TestNil(t *testing.T) {
	e := errors.New("Some error")
	assert.NotNil(t, e)
}
```

**Len**: 该函数判断一个给定切片的长度是否与期望值相等

```go
func TestLen(t *testing.T) {
	arr := []int{1, 2, 3}
	assert.Len(t, arr, 3)
}
```

**IsType**: 该函数判断一个给定对象的类型是否与期望值对象的类型相同

```go
func TestIsType(t *testing.T) {
	arr1 := []int{1, 2, 3}
	arr2 := []string{"1", "2", "3"}

	// 这里测试会无法通过
	assert.IsType(t, arr1, arr2)
}
```

更多的断言函数请查阅源代码，其中有详细的注释。