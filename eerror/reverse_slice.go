package eerror

type reverseSlice struct {
	sl  Stack
	len int
}

func (r *reverseSlice) Len() int {
	return r.len
}

func (r *reverseSlice) Cap() int {
	return cap(r.sl)
}

func (r *reverseSlice) Push(v interface{}) {
	if r.len == cap(r.sl) {
		newSl := make(Stack, cap(r.sl)*2)
		for i := range r.sl {
			newSl[len(newSl)-i-1] = r.sl[len(r.sl)-1-i]
		}
		r.sl = newSl
	}
	r.sl[len(r.sl)-1-r.len] = v
	r.len++
}

func (r *reverseSlice) Pop() interface{} {
	if r.len == 0 {
		return nil
	}
	res := r.sl[len(r.sl)-1]
	r.sl = r.sl[:len(r.sl)-1]
	r.len--
	return res
}

func (r *reverseSlice) Get() Stack {
	if r.len == 0 {
		return nil
	}
	return r.sl[len(r.sl)-r.len:]
}

func newReverseSlice(size int) *reverseSlice {
	return &reverseSlice{
		sl:  make(Stack, size),
		len: 0,
	}
}
