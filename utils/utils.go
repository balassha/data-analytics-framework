package utils

import (
	"os"
	"strconv"
	"strings"
)

func ParseStartAndEndTime(delivery string) (int, int, error) {
	items := strings.Split(delivery, " ")
	start, err := strconv.Atoi(string([]rune(items[1])[0]))
	if err != nil {
		return 0, 0, err
	}
	if start == 12 {
		start = 0
	}
	end, err := strconv.Atoi(string([]rune(items[3])[0]))
	if err != nil {
		return 0, 0, err
	}
	if end != 12 {
		end += 12
	}
	return start, end, nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 0, 64)
	return err == nil
}

func IsTimeValid(s string, am string) bool {
	length := len(s)
	if length < 3 || length > 4 {
		return false
	}

	if string([]rune(s)[length-2:]) != am {
		return false
	}
	if length == 3 {
		str := string([]rune(s)[0])
		num, _ := strconv.Atoi((str))
		if !IsNumeric(str) || num < 1 || num > 12 {
			return false
		}
	} else {
		str := string([]rune(s)[:2])
		num, _ := strconv.Atoi((str))
		if !IsNumeric(str) || num < 1 || num > 12 {
			return false
		}
	}
	return true
}
