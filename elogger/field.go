package elogger

func Any(key string, value interface{}) Field {
	return Field{
		Key: key,
		Value: Value{
			Val: value,
		},
	}
}

func Group(key string, field ...Field) Field {
	value := make(map[string]interface{}, len(field))
	for _, f := range field {
		value[f.Key] = f.Value
	}
	return Field{
		Key: key,
		Value: Value{
			Val: value,
		},
	}
}

func GroupValue(field ...Field) Value {
	value := make(map[string]interface{}, len(field))
	for _, f := range field {
		value[f.Key] = f.Value
	}
	return Value{
		Val: value,
	}
}
