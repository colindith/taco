package newsfeed

import "fmt"

type Getter interface {
	GetAll() []Item
}

type Adder interface {
	Add(item Item)
}

type Item struct {
	Title string `json:"title"`
	Post string `json:"Post"`
}

type Repo struct {
	Items []Item
}

func New() *Repo {
	return &Repo {
		Items: []Item{},
	}
}

func (r *Repo) Add(item Item) {
	fmt.Println("----------------------")
  fmt.Println(item)
	r.Items = append(r.Items, item)
}

func (r *Repo) GetAll() []Item {
	return r.Items
}