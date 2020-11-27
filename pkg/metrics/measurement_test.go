package metrics

import (
	"testing"
)

func Test_Measurement_001(t *testing.T) {
	good := []string{
		"",
		"m1 int32",
		"m1,m2 uint32",
		"s1 bool s2 string",
		"s1 bool s2,s3 float32",
	}
	for i, v := range good {
		if m, err := parseMetrics(v); err != nil {
			t.Error(i, err)
		} else {
			t.Log(v, "=>", m)
		}
	}
}

func Test_Measurement_002(t *testing.T) {
	bad := []string{
		"m1",
		"m1,m2",
		"s1 bool s2",
		"s1 bool s2,s3, float32",
		"s1 bool2 s2,s3 float32",
	}
	for i, v := range bad {
		if _, err := parseMetrics(v); err == nil {
			t.Error(i, "Expected error for", v)
		} else {
			t.Log(i, v, "=>", err)
		}
	}
}

func Test_Measurement_003(t *testing.T) {
	bad := []string{
		"m1",
		"m1,m2",
		"s1 bool s2",
		"s1 bool s2,s3, float32",
		"s1 bool2 s2,s3 float32",
		"s1,s1 bool",
		"s1 bool s1 float32",
	}
	for i, v := range bad {
		if _, err := parseMetrics(v); err == nil {
			t.Error(i, "Expected error for", v)
		} else {
			t.Log(i, v, "=>", err)
		}
	}
}

func Test_Measurement_004(t *testing.T) {
	if _, err := NewMeasurement("", ""); err == nil {
		t.Error("Expected error with empty name")
	}
	if _, err := NewMeasurement("test", ""); err == nil {
		t.Error("Expected error with empty metrics")
	}
	if m, err := NewMeasurement("test", "m1 uint8"); err != nil {
		t.Error(err)
	} else if len(m.Metrics()) != 1 {
		t.Error("Unexpected number of metrics")
	} else if len(m.Tags()) != 0 {
		t.Error("Unexpected number of tags")
	} else {
		t.Log(m)
	}

	if m, err := NewMeasurement("test", "m1 uint8", NewField("tag")); err != nil {
		t.Error(err)
	} else if len(m.Metrics()) != 1 {
		t.Error("Unexpected number of metrics")
	} else if len(m.Tags()) != 1 {
		t.Error("Unexpected number of tags")
	} else {
		t.Log(m)
	}

	if _, err := NewMeasurement("test", "m1 uint8", NewField("tag"), NewField("tag")); err == nil {
		t.Error("Expected error with duplicate tag name")
	} else {
		t.Log("Expected error", err)
	}
}
