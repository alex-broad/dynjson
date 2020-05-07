# dynjson
User-customizable JSON formats for dynamic APIs.

## Introduction

dynjson allow APIs to return only fields selected by the API client :

```
GET https://api.example.com/v1/foos
[{"id":1,foo":1,"bar":2,"baz":3}]

GET https://api.example.com/v1/foos?select=foo
[{"foo":1}]

GET https://api.example.com/v1/foos/1?select=foo
{"foo":1}
```

## Installation

go get github.com/cocoon-space/dynjson

## Usage

Standard use :

```go
type APIResult struct {
    Foo int     `json:"foo"`
    Bar string  `json:"bar"`
}

f := dynjson.NewFormatter()

res := &APIResult{Foo:1, Bar:"bar"}
o, err := f.FormatObject(o, []string{"foo"}, nil)
if err != nil {
    // handle error
}
err := json.Marshal(w, o) // {"foo": 1}
```