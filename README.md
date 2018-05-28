# ptail
Read the file from the end with golang and execute the hook

# usage
```golang
l := newPtail("/path/to/example.log")
l.parse(100, func(line string) {
  fmt.Println(line)
})
```

# author
@pyama86
