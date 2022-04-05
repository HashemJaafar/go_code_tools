//this file created automatically
package main

import (
	"fmt"
	"testing"
)

func Test_IS_IN(t *testing.T) {
	// a := IS_IN(1, []int{1, 2, 3})
	// e := true
	// TEST(true, a, e)
	// a = IS_IN(5, []int{1, 2, 3})
	// e = false
	// TEST(true, a, e)
}
func Test_SWAP(t *testing.T) {
	//TODO
	//a1:=SWAP()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_err_panic(t *testing.T) {
	//TODO
	//a1:=err_panic()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_file_lines(t *testing.T) {
	a, a2 := file_lines("/home/hashem/go_code_tools/a.txt")
	e := []string{"this file for test", "this file for test", ""}
	TEST(true, a, e)
	TEST(true, a2, nil)
}
func Test_is_func_line(t *testing.T) {
	a := is_func_line("func Test_is_func_line(t *testing.T) {")
	e := true
	TEST(true, a, e)
	a = is_func_line("Test_is_func_line")
	e = false
	TEST(true, a, e)
}
func Test_is_one_line_func(t *testing.T) {
	//TODO
	//a1:=is_one_line_func()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_main(t *testing.T) {
	//TODO
	//a1:=main()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_read_folders(t *testing.T) {
	read_folders("/home/hashem/go_code_tools")
	TEST(true, golang_files, []string{
		"/home/hashem/go_code_tools/lib.go",
		"/home/hashem/go_code_tools/line_counter.go",
		"/home/hashem/go_code_tools/make_test_func.go",
		"/home/hashem/go_code_tools/test_tools.go",
		"/home/hashem/go_code_tools/unused_func.go",
	})
	TEST(true, golang_test_files, []string{
		"/home/hashem/go_code_tools/lib_test.go",
		"/home/hashem/go_code_tools/make_test_func_test.go",
		"/home/hashem/go_code_tools/test_tools_test.go",
	})
}
func Test_return_all_func_code_block(t *testing.T) {
	a1, _ := return_all_func_code_block([]string{"func SWAP[t any](s []t, a, b int) { s[a], s[b] = s[b], s[a] }"})
	TEST(true, len(a1), 1)
	TEST(true, a1[0].func_name, "SWAP")
	TEST(true, a1[0].func_block, []string{"func SWAP[t any](s []t, a, b int) { s[a], s[b] = s[b], s[a] }"})
}
func Test_return_func_name(t *testing.T) {
	a := return_func_name("func Test_return_func_name(t *testing.T) {")
	e := "Test_return_func_name"
	TEST(true, a, e)
}
func Test_return_func_test_line(t *testing.T) {
	//TODO
	//a1:=return_func_test_line()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_return_func_test_name(t *testing.T) {
	//TODO
	//a1:=return_func_test_name()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_return_package_name(t *testing.T) {
	x, _ := file_lines("/home/hashem/go_code_tools/lib_test.go")
	a := return_package_name(x)
	e := "main"
	TEST(true, a, e)
}
func Test_return_test_file_name(t *testing.T) {
	//TODO
	//a1:=return_test_file_name()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_return_the_func_code_block_from_name(t *testing.T) {
	//TODO
	//a1:=return_the_func_code_block_from_name()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_sort_the_func_block(t *testing.T) {
	lines, _ := file_lines("/home/hashem/go_code_tools/lib.go")
	funcs_code_block, lines := return_all_func_code_block(lines)
	sort_the_func_block(funcs_code_block)
	for _, b := range funcs_code_block {
		fmt.Println(b.func_name)
	}
	// TEST(true, a, e)
}
func Test_sort_the_test_func_by_the_original_func(t *testing.T) {
	//TODO
	//a1:=sort_the_test_func_by_the_original_func()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_write_all_test_func_in_all_files(t *testing.T) {
	//TODO
	//a1:=write_all_test_func_in_all_files()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_write_lines(t *testing.T) {
	//TODO
	//a1:=write_lines()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_write_test_func(t *testing.T) {
	//TODO
	//a1:=write_test_func()
	//e1:=
	//TEST(true,a1,e1)
}
func Test_write_test_tools(t *testing.T) {
	write_test_tools("/home/hashem/go_code_tools")
	//TEST(true,a,e)
}
func Test_write_the_first_lines_of_test_file(t *testing.T) {
	//TODO
	//a1:=write_the_first_lines_of_test_file()
	//e1:=
	//TEST(true,a1,e1)
}
