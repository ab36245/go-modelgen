package godefs

func doMap[A, B any](as []A, f func(A) B) []B {
	var bs []B
	for _, a := range as {
		bs = append(bs, f(a))
	}
	return bs
}
