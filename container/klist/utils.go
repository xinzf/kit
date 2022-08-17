package klist

import "github.com/xinzf/kit/container/kvar"

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

func Chunk[T any](list *List[T], splitNum int) []*List[T] {
	if list.Size() == 0 {
		return []*List[T]{}
	}

	if splitNum < 1 {
		return []*List[T]{list}
	}

	chunks := make([]*List[T], 0)
	{
		chunks = append(chunks, New[T]())
	}

	i := 0
	list.Each(func(idx int, elem T) {
		chunks[len(chunks)-1].Add(elem)
		i = i + 1
		if i == splitNum {
			i = 0
			if idx < list.Size()-1 {
				chunks = append(chunks, New[T]())
			}
		}
	})
	return chunks
}

// Intersect 交集，属于set且属于others的元素为元素的集合，others 只要有一个不相符即抛弃
func Intersect[T any](lists ...*List[T]) (newList *List[T]) {
	newList = New[T]()
	if lists == nil || len(lists) == 0 {
		return
	}

	this := lists[0]
	others := lists[0:]
	this.Each(func(_ int, a T) {
		same := true
		for _, other := range others {
			idx, _ := other.Find(func(b T) bool {
				return kvar.New(a).Equal(kvar.New(b))
			})
			if idx == -1 {
				same = false
				break
			}
		}
		if same {
			newList.Add(a)
		}
	})
	return
}

// Diff 差集，属于set且不属于others的元素为元素的集合，others只要有一个相符，即抛弃
func Diff[T any](lists ...*List[T]) (newList *List[T]) {
	newList = New[T]()
	if lists == nil || len(lists) == 0 {
		return
	}

	this := lists[0]
	others := lists[0:]
	this.Each(func(_ int, a T) {
		same := true
		for _, other := range others {
			if idx, _ := other.Find(func(b T) bool {
				return kvar.New(a).Equal(kvar.New(b))
			}); idx != -1 {
				same = false
				break
			}
		}
		if !same {
			newList.Add(a)
		}
	})

	return
}

// Union 并集，属于set或属于others的元素为元素的集合。
func Union[T any](lists ...*List[T]) (newList *List[T]) {
	newList = New[T]()
	if lists == nil || len(lists) == 0 {
		return
	}

	for _, other := range lists {
		other.Each(func(_ int, a T) {
			idx, _ := newList.Find(func(b T) bool {
				return kvar.New(a).Equal(kvar.New(b))
			})
			if idx == -1 {
				newList.Add(a)
			}
		})
	}
	return
}
