package klist

// Transform  使用映射函数将 []T1 转换成 []T2。
// This function has two type parameters, T1 and T2.
// 映射函数 f 接受两个类型类型 T1 和 T2。
// 本函数可以处理所有类型的切片数据。
func Transform[T1, T2 any](s *List[T1], f func(T1) T2) *List[T2] {
	r := New[T2]()
	s.Each(func(idx int, elem T1) {
		r.Add(f(elem))
	})
	return r
}
