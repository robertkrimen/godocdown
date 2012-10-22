package main

import (
	_dbg "github.com/robertkrimen/dbg"
)

func dbg(input ...interface{}) {
	_dbg.Dbg(input...)
}
