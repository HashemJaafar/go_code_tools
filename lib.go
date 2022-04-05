package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	_test        = "_test.go"
	package_     = "package "
	func_        = "func "
	len_func_    = 5
	len_package_ = 8
)

var (
	golang_files      []string
	golang_test_files []string
	error_not_found   = errors.New("func_name not found")
)

type func_block struct {
	func_name  string
	func_block []string
}

func POPUP[t any](s []t, a int) []t { return append(s[:a], s[a+1:]...) }

func SWAP[t any](s []t, a, b int) { s[a], s[b] = s[b], s[a] }
func err_panic(err error) {
	if err != nil {
		log.Panic(err)
	}
}
func file_lines(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(content), "\n"), nil
}
func is_func_line(line string) bool {
	return len(line) > len_func_ && line[:len_func_] == func_
}
func is_one_line_func(line string) bool {
	if is_func_line(line) {
		if string(line[len(line)-1]) == "}" {
			return true
		}
	}
	return false
}
func main() {
	write_test_tools(".")
	sort_the_test_func_by_the_original_func()
}
func read_folders(folder_name string) {
	files, err := ioutil.ReadDir(folder_name)
	err_panic(err)
	for _, f := range files {
		dir := folder_name + "/" + f.Name()
		switch {
		case f.IsDir():
			read_folders(dir)
		case filepath.Ext(f.Name()) == ".go":
			if strings.Contains(f.Name(), _test) {
				golang_test_files = append(golang_test_files, dir)
			} else {
				golang_files = append(golang_files, dir)
			}
		}
	}
}
func return_all_func_code_block(lines []string) ([]func_block, []string) {
	var funcs_code_block []func_block
	var func_code_block func_block
	var we_are_in_func_code_block bool
	var new_lines []string
	for _, a := range lines {
		if !we_are_in_func_code_block {
			if is_func_line(a) {
				we_are_in_func_code_block = true
				func_code_block.func_name = return_func_name(a)
			} else if a != "" {
				new_lines = append(new_lines, a)
			}
		}
		if we_are_in_func_code_block {
			func_code_block.func_block = append(func_code_block.func_block, a)
		}
		if we_are_in_func_code_block && (a == "}" || is_one_line_func(a)) {
			we_are_in_func_code_block = false
			funcs_code_block = append(funcs_code_block, func_code_block)
			func_code_block.func_block = []string{}
		}
	}
	return funcs_code_block, new_lines
}
func return_func_name(line string) string {
	var func_name string
	if is_func_line(line) {
		for _, a := range line[len_func_:] {
			if a == '(' || a == '[' {
				return func_name
			}
			func_name += string(a)
		}
	}
	return func_name
}
func return_func_test_line(func_name string) string {
	return "func Test_" + func_name + "(t *testing.T) {"
}
func return_func_test_name(func_name string) string {
	return "Test_" + func_name
}
func return_package_name(lines []string) string {
	for _, a := range lines {
		if len(a) >= len_package_ && a[:len_package_] == package_ {
			return a[len_package_:]
		}
	}
	return ""
}
func return_test_file_name(path string) string {
	return path[:len(path)-3] + _test
}
func return_the_func_code_block_from_name(func_name string, slice_of_func_block []func_block) (func_block, int, error) {
	for indexa, a := range slice_of_func_block {
		if a.func_name == func_name {
			return a, indexa, nil
		}
	}
	return func_block{}, 0, error_not_found
}
func sort_the_func_block(slice_of_func_block []func_block) {
	for i := 0; i < len(slice_of_func_block); i++ {
		for j := i + 1; j < len(slice_of_func_block); j++ {
			if slice_of_func_block[i].func_name > slice_of_func_block[j].func_name {
				SWAP(slice_of_func_block, i, j)
			}
		}
	}
}
func sort_the_test_func_by_the_original_func() {
	golang_files = []string{}
	golang_test_files = []string{}
	read_folders(".")

	all_test_func_code_block := []func_block{}
	for _, a := range golang_test_files {
		test_file_name := return_test_file_name(a)
		lines_test, _ := file_lines(test_file_name)
		funcs_code_block_test, lines_test := return_all_func_code_block(lines_test)
		all_test_func_code_block = append(all_test_func_code_block, funcs_code_block_test...)
	}
	for _, a := range golang_files {

		lines, _ := file_lines(a)
		package_name := return_package_name(lines)
		funcs_code_block, lines := return_all_func_code_block(lines)

		if len(funcs_code_block) == 0 {
			continue
		}

		test_file_name := return_test_file_name(a)
		lines_test, err := file_lines(test_file_name)
		_, lines_test = return_all_func_code_block(lines_test)

		if err != nil {
			lines_test = write_the_first_lines_of_test_file(package_name)
		}

		sort_the_func_block(funcs_code_block)

		f, _ := os.Create(a)
		write_lines(f, lines)
		f_test, _ := os.Create(test_file_name)
		write_lines(f_test, lines_test)
		for _, b := range funcs_code_block {
			write_lines(f, b.func_block)
			test_func, index, err := return_the_func_code_block_from_name(return_func_test_name(b.func_name), all_test_func_code_block)
			if err != nil {
				test_func.func_block = write_test_func(b.func_name)
			}
			write_lines(f_test, test_func.func_block)
			POPUP(all_test_func_code_block, index)
		}
	}
}
func write_lines(f *os.File, lines []string) {
	for _, a := range lines {
		f.WriteString("\n" + a)
	}
}
func write_test_func(func_name string) []string {
	return []string{
		return_func_test_line(func_name),
		"\t//TODO",
		"\t//a1:=" + func_name + "()",
		"\t//e1:=",
		"\t//TEST(true,a1,e1)",
		"}",
	}
}
func write_test_tools(path string) {
	f, _ := os.Create(path + "/test_tools.txt")
	f.WriteString(
		`//this file created automatically
package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"text/tabwriter"
)

 func TEST[t any](should_equal bool, actual, expected t) {
	if !(reflect.DeepEqual(actual, expected) == should_equal) {
		p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(p, "\033[35m", "should_equal\t:", should_equal) //purple
		fmt.Fprintln(p, "\033[34m", "actual\t:", actual)             //blue
		fmt.Fprintln(p, "\033[33m", "expected\t:", expected)         //yellow
		p.Flush()
		fmt.Println("\033[31m") //red
		log.Panic()
	}
}`)
}
func write_the_first_lines_of_test_file(package_name string) []string {
	return []string{
		"//this file created automatically",
		"package " + package_name,
		"",
		"import \"testing\"",
	}
}
