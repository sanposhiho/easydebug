package main

func test() int {
	hoge := 1
	dmp("hoge", hoge)

	println(hoge)

	hoge = 2
	dmp("hoge", hoge)

	fuga := 3
	dmp("fuga", fuga)

	if hoge == 2 {
		fuga = abusoluteTwo()
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
      fmt.Printf("%s: %#v\n",valueName, vv)
  }
}
