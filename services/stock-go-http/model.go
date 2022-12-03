package main

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewCategory struct {
	Name string `json:"name"`
}
