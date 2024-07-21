package elogger

func Any(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
