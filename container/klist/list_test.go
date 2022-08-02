package klist

import (
	"fmt"
	"testing"
)

type element struct {
	id int `json:"id"`
}

func (this *element) Equal(ele any) bool {
	return this.id == ele.(*element).id
}

func TestList_Add(t *testing.T) {
	list := New[*element]()
	list.Add(&element{id: 1})
	list.Add(&element{id: 1})
	fmt.Println(list.String())
}

func TestList_Sort(t *testing.T) {
	list := New[*element]()
	list.Add(&element{id: 1})
	list.Add(&element{id: 3})
	list.Add(&element{id: 2})
	list.Each(func(idx int, elem *element) {
		fmt.Println(idx, elem)
	})
	fmt.Println("---------")
	list.Sort(func(a, b *element) bool {
		return a.id > b.id
	})
	list.Each(func(idx int, elem *element) {
		fmt.Println(idx, elem)
	})
	fmt.Println("---------")

	list.Swap(0, 1)
	list.Each(func(idx int, elem *element) {
		fmt.Println(idx, elem)
	})

}

func TestList_Union(t *testing.T) {
	list1 := New[string]("a", "b")
	list2 := New[string]("c", "d")
	list3 := New[string]("b", "d", "e")

	list := New[string]()
	list = list.Union(list1, list2, list3)
	fmt.Println(list.List())
}

func TestList_Chunk(t *testing.T) {
	list := New[string]("b", "d", "e", "f")
	chunks := list.Chunk(2)
	for _, l := range chunks {
		fmt.Println(l.List())
	}
}
