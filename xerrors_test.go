package xerrors

import (
	"fmt"
	"testing"
)

func testSubcall() error {
	err := Wrap(fmt.Errorf("%s", "abc")).WithInt(1001)
	return err
}

func testSubcall2() error {
	err := fmt.Errorf("abc")
	return err
}

func testSubcall3() error {
	err := Wrap(fmt.Errorf("%s", "efg")).WithInt(1002).WithMessage("a good guy")
	return err
}

func testSubcall4() error {
	err := testSubcall3()
	return Wrap(err)
}

func TestError(t *testing.T) {
	err := testSubcall()
	if err != nil {
		t.Logf("err: %s", err.Error())
	}

	err = testSubcall2()
	if err != nil {
		t.Logf("err: %s", err.Error())
	}

	err = testSubcall3()
	if err != nil {
		t.Logf("err: %s", err.Error())
	}

	t.Logf("result: %d", Int(nil))

	err = testSubcall4()
	if err != nil {
		t.Logf("err: %s", err.Error())
	}
}

func TestWrap(t *testing.T) {
	SetSysInternalError(2)

	err1 := Wrap(fmt.Errorf("abc")).WithInt(1001).WithMessage("inner")
	err2 := Wrap(err1).WithInt(1002).WithMessage("outer")
	t.Logf("err: %s", err2.Error())
	t.Logf("retcode: %d", Int(err2))
}
