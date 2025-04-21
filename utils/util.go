package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

var dayToInt = map[string]int{
	"Monday":    1,
	"Mon":       1,
	"Tuesday":   2,
	"Tue":       2,
	"Wednesday": 3,
	"Wed":       3,
	"Thursday":  4,
	"Thur":      4,
	"Friday":    5,
	"Fri":       5,
	"Saturday":  6,
	"Sat":       6,
	"Sunday":    7,
	"Sun":       7,
}

type OpenDayTime map[string][]string

func (o OpenDayTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]string(o))
}

func (o *OpenDayTime) UnmarshalJSON(data []byte) error {
	var temp map[string][]string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	*o = OpenDayTime(temp)
	return nil
}

func (o OpenDayTime) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *OpenDayTime) Scan(value interface{}) error {
	if value == nil {
		*o = OpenDayTime{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan OpenDayTime: %v", value)
	}
	return o.UnmarshalJSON(bytes)
}

// "Mon - Wed 08:00 - 17:00 / Thur, Sat 20:00 - 02:00", to openDayTime formate
func FormateOpeningHours(s string) OpenDayTime {
	// Parse the string and populate the openDayTime map
	o := make(OpenDayTime)
	// Split the input string by '/' to handle multiple day-time ranges
	parts := strings.Split(s, "/")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		startDay := 0
		endDay := 0
		startTime := ""
		endTime := ""
		dayAndTime := strings.Split(part, " ")

		//if dayAndTime[1] == "-" range mode
		//if layAndTime[1]  != "-" single mode
		if dayAndTime[1] == "-" {
			startDay = dayToInt[strings.TrimSpace(dayAndTime[0])]
			endDay = dayToInt[strings.TrimSpace(dayAndTime[2])]
			startTime = strings.TrimSpace(dayAndTime[3])
			endTime = strings.TrimSpace(dayAndTime[5])
			addDaysInRange(o, startDay, endDay, startTime, endTime)
		} else {
			days := []int{}
			for _, v := range dayAndTime {
				if day, ok := dayToInt[strings.Trim(v, ",")]; ok {
					days = append(days, day)
				} else {
					startTime = strings.TrimSpace(dayAndTime[len(dayAndTime)-3])
					endTime = strings.TrimSpace(dayAndTime[len(dayAndTime)-1])
				}
			}
			for _, day := range days {
				addDaysInRange(o, day, day, startTime, endTime)
			}
		}
	}
	return o
}

func addDaysInRange(o OpenDayTime, startDay, endDay int, startTime, endTime string) {
	days := []string{"Mon", "Tue", "Wed", "Thur", "Fri", "Sat", "Sun"}

	for _, day := range days[startDay-1 : endDay] {
		o[day] = []string{startTime, endTime}
	}
}

