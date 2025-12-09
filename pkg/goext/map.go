package goext

type MapFunc[A any, B any] func(A) B

func Map[A any, B any](input []A, m MapFunc[A, B]) []B {
	output := make([]B, len(input))
	for i, element := range input {
		output[i] = m(element)
	}
	return output
}
