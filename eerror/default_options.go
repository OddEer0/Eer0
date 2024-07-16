package eerror

func DefaultCode(code ErrorCode) {
	mu.Lock()
	defer mu.Unlock()
	defaultErrorCode = code
}

func DefaultStackDepth(size int) {
	mu.Lock()
	defer mu.Unlock()
	stackDepth = size
}

func DefaultErrorMethod(fn ErrorMethod) {
	mu.Lock()
	defer mu.Unlock()
	errorFunc = fn
}

func DefaultCauseWrapSeparator(separator string) {
	mu.Lock()
	defer mu.Unlock()
	wrapSeparator = separator
}

func DefaultWrapCap(capacity int) {
	mu.Lock()
	defer mu.Unlock()
	wrapDepth = capacity
}

func DefaultInfoSize(size int) {
	mu.Lock()
	defer mu.Unlock()
	infoCap = size
}
