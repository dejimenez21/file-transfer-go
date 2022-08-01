package main

type file struct {
	Name    string
	Ext     string
	Size    int64
	Content []byte
}

func (f *file) FullName() string {
	return f.Name + "." + f.Ext
}
