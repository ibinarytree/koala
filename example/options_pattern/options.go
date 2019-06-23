package main

import "fmt"

type Options struct {
	StrOption1 string
	StrOption2 string
	StrOption3 string
	IntOption1 int
	IntOption2 int
	IntOption3 int
}

func InitOptions1(strOption1 string, strOption2 string, strOption3 string,
	IntOption1 int, IntOption2 int, IntOption3 int) {
	options := &Options{}
	options.StrOption1 = strOption1
	options.StrOption2 = strOption2
	options.StrOption3 = strOption3
	options.IntOption1 = IntOption1
	options.IntOption2 = IntOption2
	options.IntOption3 = IntOption3

	fmt.Printf("init options1:%#v\n", options)
	return
}

func InitOptions2(opts ...interface{}) {
	options := &Options{}
	for index, opt := range opts {
		switch index {
		case 0:
			str, ok := opt.(string)
			if !ok {
				return
			}
			options.StrOption1 = str
		case 1:
			str, ok := opt.(string)
			if !ok {
				return
			}
			options.StrOption2 = str
		case 2:
			str, ok := opt.(string)
			if !ok {
				return
			}
			options.StrOption3 = str
		case 3:
			val, ok := opt.(int)
			if !ok {
				return
			}
			options.IntOption1 = val
		case 4:
			val, ok := opt.(int)
			if !ok {
				return
			}
			options.IntOption2 = val
		case 5:
			val, ok := opt.(int)
			if !ok {
				return
			}
			options.IntOption3 = val
		}
	}
	fmt.Printf("init options2:%#v\n", options)
}

type Option func(opts *Options)

func InitOption3(opts ...Option) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	fmt.Printf("init options3:%#v\n", options)
}

func WithStringOption1(str string) Option {
	return func(opts *Options) {
		opts.StrOption1 = str
	}
}

func WithStringOption2(str string) Option {
	return func(opts *Options) {
		opts.StrOption2 = str
	}
}

func WithStringOption3(str string) Option {
	return func(opts *Options) {
		opts.StrOption3 = str
	}
}

func WithIntOption1(val int) Option {
	return func(opts *Options) {
		opts.IntOption1 = val
	}
}

func WithIntOption2(val int) Option {
	return func(opts *Options) {
		opts.IntOption2 = val
	}
}

func WithIntOption3(val int) Option {
	return func(opts *Options) {
		opts.IntOption3 = val
	}
}

func main() {
	InitOptions1("str1", "str2", "str3", 1, 2, 3)
	InitOptions2("str1", "str2", "str3", 1, 2, 3)
	InitOption3(
		WithStringOption1("str1"),
		WithStringOption3("str3"),
		WithIntOption3(3),
		WithIntOption2(2),
		WithIntOption1(1),
		WithStringOption2("str2"),
	)
}
