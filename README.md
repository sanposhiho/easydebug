# easydebug
tool: add debug statements after statements to store values

- You can add debug statements after all statements to store values.
- You can remove all debug statements with this tool, of course.
- You can edit debug function `dmp` to arrange debug as you like


## Installation

```
go get -u github.com/sanposhiho/easydebug
```

## Usage

### add debug statements

```
easydebug -f target.go -mode 0
```

### remove debug statements

```
easydebug -f target.go -mode 1
```

## sample

### before

```
package main

func test() int {
	hoge := 1

	fuga := 3

	if hoge == 2 {
		fuga = abusoluteTwo()

		hoge = 12

	}

	return hoge + fuga
}

func abusoluteTwo() int {
	return 2
}
```

### after

```
package main

func test() int {
	hoge := 1
	dmp("hoge", hoge)

	fuga := 3
	dmp("fuga", fuga)

	if hoge == 2 {
		fuga = abusoluteTwo()
		dmp("fuga", fuga)

		hoge = 12
		dmp("hoge", hoge)

	}

	return hoge + fuga
}

func abusoluteTwo() int {
	return 2
}

// generated from goeasydebug
// function for data dump
func dmp(valueName string, v ...interface{}) {
  for _, vv := range(v) {
      // arrange debug as you like
      fmt.Printf("%s: %#v\n",valueName, vv)
  }
}
```
