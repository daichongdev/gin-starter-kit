package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONTime 自定义时间类型，专门用于JSON序列化
type JSONTime time.Time

// MarshalJSON 实现json.Marshaler接口
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	t := time.Time(jt)
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + t.Format("2006-01-02 15:04:05") + `"`), nil
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (jt *JSONTime) UnmarshalJSON(data []byte) error {
	if string(data) == `""` || string(data) == "null" {
		*jt = JSONTime(time.Time{})
		return nil
	}

	str := string(data[1 : len(data)-1]) // 去掉引号
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		// 尝试RFC3339格式
		t, err = time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
	}
	*jt = JSONTime(t)
	return nil
}

// Value 实现driver.Valuer接口，用于数据库存储
func (jt JSONTime) Value() (driver.Value, error) {
	t := time.Time(jt)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

// Scan 实现sql.Scanner接口，用于数据库读取
func (jt *JSONTime) Scan(value interface{}) error {
	if value == nil {
		*jt = JSONTime(time.Time{})
		return nil
	}
	if t, ok := value.(time.Time); ok {
		*jt = JSONTime(t)
		return nil
	}
	return fmt.Errorf("cannot scan %T into JSONTime", value)
}

// ToTime 转换为标准time.Time类型
func (jt JSONTime) ToTime() time.Time {
	return time.Time(jt)
}

// Now 创建当前时间的JSONTime
func Now() JSONTime {
	return JSONTime(time.Now())
}

// FromTime 从time.Time创建JSONTime
func FromTime(t time.Time) JSONTime {
	return JSONTime(t)
}
