package util

func ThunkOf[T any](t T) func() T {
	return func() T {
		return t
	}
}

func ConstOf[S any, T any](t T) func(S) T {
	return func(S) T {
		return t
	}
}

func LiftToNullary[S, T any](f func(S) T) func(func() S) func() T {
	return func(sSrc func() S) func() T {
		return func() T {
			return f(sSrc())
		}
	}
}
