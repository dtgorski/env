// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2021

package env

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestUnmarshal_1(t *testing.T) {
	if err := Unmarshal(nil); err == nil {
		t.Error("unexpected")
	}
	if err := Unmarshal("x"); err == nil {
		t.Error("unexpected")
	}
	if err := Unmarshal(&t); err == nil {
		t.Error("unexpected")
	}
}

func TestUnmarshal_2(t *testing.T) {
	type Testee struct {
		ValInt     int     `env:"VAL_INT"`
		ValInt8    int8    `env:"VAL_INT8"`
		ValInt16   int16   `env:"VAL_INT16"`
		ValInt32   int32   `env:"VAL_INT32"`
		ValInt64   int64   `env:"VAL_INT64"`
		ValUint    uint    `env:"VAL_UINT"`
		ValUint8   uint8   `env:"VAL_UINT8"`
		ValUint16  uint16  `env:"VAL_UINT16"`
		ValUint32  uint32  `env:"VAL_UINT32"`
		ValUint64  uint64  `env:"VAL_UINT64"`
		ValFloat32 float32 `env:"VAL_FLOAT32"`
		ValFloat64 float64 `env:"VAL_FLOAT64"`

		Section struct {
			ValBool   bool   `env:"VAL_BOOL"`
			ValString string `env:"VAL_STR"`
		}

		ArrStr  [1]string `env:"ARR_STR"`
		ListStr []string  `env:"LIST_STR"`

		ListInt   []int   `env:"LIST_INT"`
		ListInt8  []int8  `env:"LIST_INT8"`
		ListInt16 []int16 `env:"LIST_INT16"`
		ListInt32 []int32 `env:"LIST_INT32"`
		ListInt64 []int64 `env:"LIST_INT64"`

		ListUint   []uint   `env:"LIST_UINT"`
		ListUint8  []uint8  `env:"LIST_UINT8"`
		ListUint16 []uint16 `env:"LIST_UINT16"`
		ListUint32 []uint32 `env:"LIST_UINT32"`
		ListUint64 []uint64 `env:"LIST_UINT64"`

		ListFloat32 []float32 `env:"LIST_FLT32"`
		ListFloat64 []float64 `env:"LIST_FLT64"`
	}

	env := Testee{}
	envs := []struct {
		key   string
		value string
		test  func() bool
	}{
		{"VAL_INT" /*    */, "1" /*   */, func() bool { return env.ValInt == 1 }},
		{"VAL_INT8" /*   */, "2" /*   */, func() bool { return env.ValInt8 == 2 }},
		{"VAL_INT16" /*  */, "3" /*   */, func() bool { return env.ValInt16 == 3 }},
		{"VAL_INT32" /*  */, "4" /*   */, func() bool { return env.ValInt32 == 4 }},
		{"VAL_INT64" /*  */, "-5" /*  */, func() bool { return env.ValInt64 == -5 }},
		{"VAL_UINT" /*   */, "6" /*   */, func() bool { return env.ValUint == 6 }},
		{"VAL_UINT8" /*  */, "7" /*   */, func() bool { return env.ValUint8 == 7 }},
		{"VAL_UINT16" /* */, "8" /*   */, func() bool { return env.ValUint16 == 8 }},
		{"VAL_UINT32" /* */, "9 " /*  */, func() bool { return env.ValUint32 == 9 }},
		{"VAL_UINT64" /* */, " -5" /* */, func() bool { return env.ValUint64 == 0 }},
		{"VAL_FLOAT32" /**/, "1.2" /* */, func() bool { return env.ValFloat32 == 1.2 }},
		{"VAL_FLOAT64" /**/, "3.4" /* */, func() bool { return env.ValFloat64 == 3.4 }},

		{"VAL_BOOL" /*   */, "true" /**/, func() bool { return env.Section.ValBool }},
		{"VAL_STR" /*    */, "foo " /**/, func() bool { return env.Section.ValString == "foo" }},

		{"ARR_STR" /*    */, "x" /*   */, func() bool { return env.ArrStr[0] == "" }},
		{"LIST_STR" /*   */, "a,b" /* */, func() bool { return env.ListStr[0] == "a" && env.ListStr[1] == "b" }},

		{"LIST_INT" /*   */, "1,-1" /**/, func() bool { return env.ListInt[0] == 1 && env.ListInt[1] == -1 }},
		{"LIST_INT8" /*  */, "2,-2" /**/, func() bool { return env.ListInt8[0] == 2 && env.ListInt8[1] == -2 }},
		{"LIST_INT16" /* */, "3,-3" /**/, func() bool { return env.ListInt16[0] == 3 && env.ListInt16[1] == -3 }},
		{"LIST_INT32" /* */, "4,-4" /**/, func() bool { return env.ListInt32[0] == 4 && env.ListInt32[1] == -4 }},
		{"LIST_INT64" /* */, "5,-5" /**/, func() bool { return env.ListInt64[0] == 5 && env.ListInt64[1] == -5 }},

		{"LIST_UINT" /*  */, "5,-5" /**/, func() bool { return env.ListUint[0] == 5 && env.ListUint[1] == 0 }},
		{"LIST_UINT8" /* */, "6,-6" /**/, func() bool { return env.ListUint8[0] == 6 && env.ListUint8[1] == 0 }},
		{"LIST_UINT16" /**/, "7,-7" /**/, func() bool { return env.ListUint16[0] == 7 && env.ListUint16[1] == 0 }},
		{"LIST_UINT32" /**/, "8,-8" /**/, func() bool { return env.ListUint32[0] == 8 && env.ListUint32[1] == 0 }},
		{"LIST_UINT64" /**/, "9,0" /* */, func() bool { return env.ListUint64[0] == 9 && env.ListUint64[1] == 0 }},

		{"LIST_FLT32" /**/, "0.1,1" /**/, func() bool { return env.ListFloat32[0] == 0.1 && env.ListFloat32[1] == 1 }},
		{"LIST_FLT64" /**/, ",1e10" /**/, func() bool { return env.ListFloat64[0] == 1e10 }},
	}

	os.Clearenv()
	for _, e := range envs {
		_ = os.Setenv(e.key, e.value)
	}
	if err := Unmarshal(&env); err != nil {
		t.Fatal(err)
	}
	for _, e := range envs {
		if e.test != nil && !e.test() {
			t.Errorf("unexpected at %s", e.key)
		}
	}
}

func TestUnmarshal_3(t *testing.T) {
	type Testee struct {
		ValInt int `env:"VAL_INT,file"`
	}

	tmpFileName := writeTempOnce([]byte("1234\n"))
	defer func() { _ = os.Remove(tmpFileName) }()

	os.Clearenv()
	_ = os.Setenv("VAL_INT_FILE", tmpFileName)

	env := Testee{}
	if err := Unmarshal(&env); err != nil {
		t.Fatal(err)
	}
	if env.ValInt != 1234 {
		t.Error("unexpected")
	}

	_ = os.Setenv("VAL_INT", "5678")
	if err := Unmarshal(&env); err != nil {
		t.Fatal(err)
	}
	if env.ValInt != 5678 {
		t.Error("unexpected")
	}
}

func writeTempOnce(b []byte) string {
	var file *os.File
	var err error

	if file, err = ioutil.TempFile("", ""); err != nil {
		return ""
	}
	if _, err = file.Write(b); err != nil {
		return ""
	}
	if err = file.Close(); err != nil {
		return ""
	}
	return file.Name()
}

func TestCoerce(t *testing.T) {
	for _, v := range []string{"true", "True", "on", "oN", "ON", "yes", "Yes", "1"} {
		if !coerceBool(v) {
			t.Errorf("unexpected at %s", v)
		}
	}
}
