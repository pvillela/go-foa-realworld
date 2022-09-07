package util

type Func[S any, T any] func(S) T

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
