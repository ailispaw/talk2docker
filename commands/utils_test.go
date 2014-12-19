package commands

import (
	"testing"
	"time"
)

func TestTruncate(t *testing.T) {
	var (
		actual   = Truncate("1234567890", 6)
		expected = "123456"
	)
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestFormatDateTime(t *testing.T) {
	var (
		actual   = FormatDateTime(time.Date(2014, 12, 18, 18, 8, 30, 0, time.Local))
		expected = "2014-12-18 18:08:30"
	)
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestFormatBool(t *testing.T) {
	var (
		actual   = FormatBool(true, "OK", "NG")
		expected = "OK"
	)
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	actual = FormatBool(false, "OK", "NG")
	expected = "NG"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestFormatInt(t *testing.T) {
	var (
		actual   = FormatInt(12345)
		expected = "12,345"
	)
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	actual = FormatInt(-1234567)
	expected = "-1,234,567"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestFormatFloat(t *testing.T) {
	var (
		actual   = FormatFloat(12345.678)
		expected = "12,345.678"
	)
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	actual = FormatFloat(1234567.8914)
	expected = "1,234,567.891"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	actual = FormatFloat(-1234567.8915)
	expected = "-1,234,567.892"
	if actual != expected {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
