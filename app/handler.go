package main

import (
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"GET":     get,
	"SET":     set,
	"HGET":    hget,
	"HSET":    hset,
	"HGETALL": hgetall,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{} // use sync.RWMutex{} because server is supposed to handle requests concurrently.
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Err wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Err wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk
	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "Err wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	_, ok := HSETs[hash]
	if !ok {
		HSETs[hash] = map[string]string{}
	}

	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Err wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Err wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	innerValue, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "array", array: []Value{}}
	}

	elements := []Value{}

	for field, val := range innerValue {
		fieldValue := Value{typ: "bulk", bulk: field}
		elements = append(elements, fieldValue)

		valValue := Value{typ: "bulk", bulk: val}
		elements = append(elements, valValue)
	}

	resultValue := Value{typ: "array", array: elements}

	return resultValue
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}
