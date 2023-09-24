package main

import "fmt"

type Reader interface {
	ReadBook()
}

type Writer interface {
	WriteBook()
}

// 具体类型
type Book struct {
}

func (this *Book) ReadBook() {
	fmt.Println("Read a Book")
}

func (this *Book) WriteBook() {
	fmt.Println("Write a Book")
}

func main() {
	//b: pair<type:Book, value:book{}地址>
	b := &Book{}

	//r: pair<type:, value:>
	var r Reader
	//r: pair<type:Book, value:book{}地址>
	r = b

	r.ReadBook()

	var w Writer

	/*
		使用断言 r.(Writer) 将 r 转换为 Writer 接口类型并将其赋值给 w。
		这里的断言成功是因为 r 实际上指向了一个 Book 类型的实例，
		而 Book 类型同时实现了 Reader 和 Writer 接口，所以这个断言是有效的。
	*/
	//r: pair<type:Book, value:book{}地址>
	w = r.(Writer) //此处的断言为什么会成功？ 因为w r 具体的type是一致

	w.WriteBook()
}
