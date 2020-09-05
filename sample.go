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
