package jsonhelper

import (
    "encoding/json"
    "time"
)

type JSONObject map[string]interface{}

func NewJSONObject() JSONObject {
    return make(JSONObject)
}

func NewJSONObjectFromMap(m map[string]interface{}) JSONObject {
    return JSONObject(m)
}

func NewJSONObjectFromBuf(data []byte) (JSONObject, error) {
    j := make(JSONObject)
	err := json.Unmarshal(data, &j)
    return j, err
}

func (p JSONObject) String() string {
    b, _ := json.Marshal(&p)
    return string(b)
}

func (p JSONObject) Del(key string) {
    delete(p, key)
}

func (p JSONObject) Set(key string, value interface{}) {
    p[key] = value
}

func (p JSONObject) Get(key string) interface{} {
    value, _ := p[key]
    return value
}

func (p JSONObject) Len() int {
    return len(p)
}

func (p JSONObject) GetAsString(key string) string {
    value, _ := p[key]
    return JSONValueToString(value)
}

func (p JSONObject) GetAsInt(key string) int {
    value, _ := p[key]
    return JSONValueToInt(value)
}

func (p JSONObject) GetAsInt32(key string) int32 {
    value, _ := p[key]
    return JSONValueToInt32(value)
}

func (p JSONObject) GetAsUint32(key string) uint32 {
    value, _ := p[key]
    v := JSONValueToInt64(value)
    return uint32(v)
}

func (p JSONObject) GetAsInt64(key string) int64 {
    value, _ := p[key]
    return JSONValueToInt64(value)
}

func (p JSONObject) GetAsFloat64(key string) float64 {
    value, _ := p[key]
    return JSONValueToFloat64(value)
}

func (p JSONObject) GetAsBool(key string) bool {
    value, _ := p[key]
    return JSONValueToBool(value)
}

func (p JSONObject) GetAsObject(key string) JSONObject {
    value, _ := p[key]
    return JSONValueToObject(value)
}

func (p JSONObject) GetAsArray(key string) JSONArray {
    value, _ := p[key]
    return JSONValueToArray(value)
}

func (p JSONObject) GetAsDuration(key string) time.Duration {
    value, _ := p[key]
    return time.Duration(JSONValueToInt64(value))
}

func (p JSONObject) GetAsTime(key string, format string) time.Time {
    value, _ := p[key]
    return JSONValueToTime(value, format)
}

func (p JSONObject) Compact(removeFalse bool, removeEmptyStrings bool, removeZero bool, removeEmptyArrays bool, removeEmptyObjects bool) JSONObject {
    if len(p) == 0 {
        if removeEmptyObjects {
            return nil
        }
        return p
    }
    m := make(map[string]interface{})
    for k, v := range p {
        var value interface{}
        value = v
        switch t := v.(type) {
        case nil:
            continue
        case string:
            if removeEmptyStrings && len(t) == 0 {
                continue
            }
        case JSONObject:
            value = t.Compact(removeFalse, removeEmptyStrings, removeZero, removeEmptyArrays, removeEmptyObjects)
        case JSONArray:
            value = t.Compact(removeFalse, removeEmptyStrings, removeZero, removeEmptyArrays, removeEmptyObjects)
        case map[string]interface{}:
            value = NewJSONObjectFromMap(t).Compact(removeFalse, removeEmptyStrings, removeZero, removeEmptyArrays, removeEmptyObjects)
        case []interface{}:
            value = NewJSONArrayFromArray(t).Compact(removeFalse, removeEmptyStrings, removeZero, removeEmptyArrays, removeEmptyObjects)
        case float64:
            if removeZero && t == 0.0 {
                continue
            }
        case float32:
            if removeZero && t == 0.0 {
                continue
            }
        case int64:
            if removeZero && t == 0 {
                continue
            }
        case int32:
            if removeZero && t == 0 {
                continue
            }
        case int:
            if removeZero && t == 0 {
                continue
            }
        case int16:
            if removeZero && t == 0 {
                continue
            }
        case int8:
            if removeZero && t == 0 {
                continue
            }
        case byte:
            if removeZero && t == 0 {
                continue
            }
        case bool:
            if removeFalse && t == false {
                continue
            }
        }
        if value == nil {
            continue
        }
        m[k] = value
    }
    if removeEmptyObjects && len(m) == 0 {
        return nil
    }
    return NewJSONObjectFromMap(m)
}
