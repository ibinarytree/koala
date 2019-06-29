package main

type Generator interface {
	Run(opt *Option) error
}
