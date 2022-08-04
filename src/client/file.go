package main

type file struct {
	Name string
	Ext  string
	Size int64
}

func (f *file) FullName() string {
	return f.Name + "." + f.Ext
}
