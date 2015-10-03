package util
import "time"

// 2014-10-11T15:00:00.000Z のような形式を処理する
// エラーの時は、空のtimeとerrを返す
func ParseDate(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}