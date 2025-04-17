package pharmacy

import (
	"reflect"
	"testing"
)

func TestTransformCases(t *testing.T) {
	cases := []struct {
		input    string
		expected openDayTime
	}{
		{
			input:    "",
			expected: openDayTime{},
		},
		{
			input: "Mon 08:00 - 17:00",
			expected: openDayTime{
				"Mon": {"08:00", "17:00"},
			},
		},
		{
			input: "Mon - Fri 09:00 - 18:00",
			expected: openDayTime{
				"Mon":  {"09:00", "18:00"},
				"Tue":  {"09:00", "18:00"},
				"Wed":  {"09:00", "18:00"},
				"Thur": {"09:00", "18:00"},
				"Fri":  {"09:00", "18:00"},
			},
		},
	}

	for _, c := range cases {
		result := FormateOpeningHours(c.input)
		if !reflect.DeepEqual(result, c.expected) {
			t.Errorf("For input '%s', expected %v, got %v", c.input, c.expected, result)
		}
	}
}
