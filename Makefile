hoge: foo
	echo "hoge"
foo: hoge var
	echo "foo"
var: hoge var2
var2: var3
var3: hoge