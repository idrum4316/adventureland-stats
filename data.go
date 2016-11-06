package main

import (
	"sync"
)

var db = map[string]map[string]string{}
var db_lock = sync.RWMutex{}