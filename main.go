package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"text/tabwriter"
)

const (
	_test         = "_test.go"
	package_      = "package "
	func_         = "func "
	Benchmark     = "Benchmark"
	Example       = "Example"
	Fuzz          = "Fuzz"
	Test          = "Test"
	Function      = "Function"
	Method        = "Method"
	len_func_     = 5
	len_package_  = 8
	len_Benchmark = 9
	len_Example   = 7
	len_Fuzz      = 4
	len_Test      = 4
)

var (
	golang_files      []string
	golang_test_files []string
	error_not_found   = errors.New("func_name not found")
	p                 = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fail_test_number  int
)

type function struct {
	comment_line                                   []string
	function_line                                  string
	receiver_type                                  string
	function_name                                  string
	input_type                                     []string
	output_type                                    []string
	function_type                                  string //function or test or banchmark or example or fuzz or method
	func_block                                     []string
	functions_that_used_this_function              []string
	functions_that_used_in_this_function           []string
	number_of_functions_that_used_this_function    int
	number_of_functions_that_used_in_this_function int
	stage                                          int
}

func REMOVE[t any](s []t, a int) []t { return append(s[:a], s[a+1:]...) }

func SWAP[t any](s []t, a, b int) { s[a], s[b] = s[b], s[a] }

func TEST[t any](should_equal bool, actual, expected t) {
	if reflect.DeepEqual(actual, expected) != should_equal {
		fail_test_number++
		p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(p, "\033[32m", "fail_test_number\t:", fail_test_number) //green
		fmt.Fprintln(p, "\033[35m", "should_equal\t:", should_equal)         //purple
		fmt.Fprintln(p, "\033[34m", "actual\t:", actual)                     //blue
		fmt.Fprintln(p, "\033[33m", "expected\t:", expected)                 //yellow
		p.Flush()
		fmt.Println("\033[31m") //red
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(string(debug.Stack()), "\033[0m") //reset
			}
		}()
		log.Panic()
	}
}

func analyse_the_function_line(function_line string) (string, string, []string, []string, []string) {
	var (
		param_receiver string
		function_name  string
		param_type     []string
		param_input    []string
		param_output   []string
	)
	var where_the_curser_is string
	const (
		in_param_receiver = "in_param_receiver"
		in_function_name  = "in_function_name"
		in_param_type     = "in_param_type"
		in_param_input    = "in_param_input"
		in_param_output   = "in_param_output"
	)
	for k1, v1 := range function_line[len_func_:] {
		if k1 == 0 {
			if v1 == '(' {
				where_the_curser_is = in_param_receiver
			} else {
				where_the_curser_is = in_function_name
			}
		}

		switch where_the_curser_is {
		case in_param_receiver:
			if v1 == '(' {
				continue
			}
			if v1 == ')' {
				where_the_curser_is = in_function_name
				continue
			}
			param_receiver += string(v1)
		case in_function_name:
			if v1 == ')' {
				continue
			}
			if v1 == '(' {
				where_the_curser_is = in_param_input
				continue
			}
			if v1 == '[' {
				where_the_curser_is = in_param_type
				continue
			}
			function_name += string(v1)
		case in_param_type:
			if v1 == '[' {
				continue
			}
			if v1 == ']' {
				where_the_curser_is = in_param_input
				continue
			}
			// param_type += string(v1)
		case in_param_input:
			if v1 == '(' {
				continue
			}
			if v1 == ')' {
				where_the_curser_is = in_param_output
			}
			// param_input += string(v1)
		case in_param_output:
			if v1 == '(' {
				continue
			}
			if v1 == ')' {
				where_the_curser_is = ""
				continue
			}
			if v1 == '{' {
				where_the_curser_is = ""
				continue
			}
			// param_output+= string(v1)
		default:
			return param_receiver, function_name, param_type, param_input, param_output
		}
	}
	return param_receiver, function_name, param_type, param_input, param_output
}

func counter() {
	golang_files = []string{}
	golang_test_files = []string{}
	read_folders(".")

	var original_file int
	var original_functions int
	var original_line int
	var test_file int
	var test_functions int
	var test_line int
	for _, a := range golang_files {
		lines, _ := file_lines(a)
		functions, _ := return_all_func_code_block(lines)
		original_file++
		original_functions += len(functions)
		original_line += len(lines)
	}
	for _, a := range golang_test_files {
		lines, _ := file_lines(a)
		functions, _ := return_all_func_code_block(lines)
		test_file++
		test_functions += len(functions)
		test_line += len(lines)
	}
	fmt.Fprintln(p, "original_file\t", original_file)
	fmt.Fprintln(p, "original_functions\t", original_functions)
	fmt.Fprintln(p, "original_line\t", original_line)
	fmt.Fprintln(p, "test_file\t", test_file)
	fmt.Fprintln(p, "test_functions\t", test_functions)
	fmt.Fprintln(p, "test_line\t", test_line)
	fmt.Fprintln(p, "all_file\t", original_file+test_file)
	fmt.Fprintln(p, "all_functions\t", original_functions+test_functions)
	fmt.Fprintln(p, "all_line\t", original_line+test_line)
	p.Flush()
}

func drew_gragh(all_func_code_block []function) {
}

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

func find[t any](elements []t, condition func(int) bool) (t, error) {
	for k1, v1 := range elements {
		if condition(k1) {
			return v1, nil
		}
	}
	var a t
	return a, errors.New("not found")
}

func is_func_line(line string) bool {
	return len(line) > len_func_ && line[:len_func_] == func_
}

func is_one_line_func(line string) bool {
	if is_func_line(line) && string(line[len(line)-1]) == "}" {
		return true
	}
	return false
}

func is_the_func_called_in_this_line(line, func_name string) bool {
	return strings.Contains(line, func_name+"(") || strings.Contains(line, func_name+"[")
}

func is_there_a_common_element[t any](slice1, slice2 []t) bool {
	for _, a := range slice1 {
		for _, b := range slice2 {
			if reflect.DeepEqual(a, b) {
				return true
			}
		}
	}
	return false
}

func main() {
	write_test_tools(".")
	all_Function_Method := sort_the_test_func_by_the_original_func()
	set_functions_dependency(all_Function_Method)
	set_the_stages(all_Function_Method)
	drew_gragh(all_Function_Method)
	print_unused_functions(all_Function_Method)
	counter()
}

func print_unused_functions(all_functions []function) {
	for _, v1 := range all_functions {
		if v1.number_of_functions_that_used_this_function == 0 {
			fmt.Println("unused:", v1.function_name)
		}
	}
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

func return_all_func_code_block(lines []string) ([]function, []string) {
	var function_all []function
	var function_one function
	var we_are_in_function bool
	var new_lines []string

	for _, a := range lines {
		if !we_are_in_function {
			receiver_type, function_type, function_name, is_func_line := type_of_the_function(a)
			if is_func_line {
				we_are_in_function = true
				function_one.receiver_type = receiver_type
				function_one.function_type = function_type
				function_one.function_name = function_name
				function_one.function_line = a
			} else if a != "" {
				new_lines = append(new_lines, a)
			}
		}
		if we_are_in_function {
			function_one.func_block = append(function_one.func_block, a)
		}
		if we_are_in_function && (a == "}" || is_one_line_func(a)) {
			function_all = append(function_all, function_one)
			we_are_in_function = false
			function_one = function{}
		}
	}
	return function_all, new_lines
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

func return_the_func_code_block_from_name(func_name string, slice_of_func_block []function) (function, int, error) {
	for indexa, a := range slice_of_func_block {
		if a.function_name == func_name {
			return a, indexa, nil
		}
	}
	return function{}, 0, error_not_found
}

func search_and_write_func(b function, all_test_func_code_block []function, f_test *os.File, function_type string, function_to_write func(string) []string) []function {
	test_func, index, err := return_the_func_code_block_from_name(function_type+b.function_name, all_test_func_code_block)
	if err != nil {
		test_func.func_block = function_to_write(b.function_name)
	} else {
		all_test_func_code_block = REMOVE(all_test_func_code_block, index)
	}
	write_lines(f_test, test_func.func_block)
	return all_test_func_code_block
}

func set_functions_dependency(all_func_code_block []function) {
	for indexa, a := range all_func_code_block {
		for indexb, b := range all_func_code_block {
			for _, c := range b.func_block {
				if !is_func_line(c) && is_the_func_called_in_this_line(c, a.function_name) {
					all_func_code_block[indexb].functions_that_used_in_this_function = append(all_func_code_block[indexb].functions_that_used_in_this_function, a.function_name)
					all_func_code_block[indexa].functions_that_used_this_function = append(all_func_code_block[indexa].functions_that_used_this_function, b.function_name)
					all_func_code_block[indexb].number_of_functions_that_used_in_this_function++
					all_func_code_block[indexa].number_of_functions_that_used_this_function++
				}
			}
		}
	}
}

func set_the_stages(all_func_code_block []function) []function {
	// step 1: find all the functions_that_used_this_function and set the stage for them to bigger by 1
	// if there is no loop betwen the father and daughter
	for i := 0; i < len(all_func_code_block); i++ {
		for indexa := range all_func_code_block {
			is_there_loop := is_there_a_common_element(all_func_code_block[indexa].functions_that_used_in_this_function, all_func_code_block[indexa].functions_that_used_this_function)
			for _, b := range all_func_code_block[indexa].functions_that_used_this_function {
				_, index, _ := return_the_func_code_block_from_name(b, all_func_code_block)
				if all_func_code_block[index].stage <= all_func_code_block[indexa].stage {
					if is_there_loop {
						all_func_code_block[index].stage = all_func_code_block[indexa].stage
					} else {
						all_func_code_block[index].stage = all_func_code_block[indexa].stage + 1
					}
				}
			}
		}
	}

	for _, a := range all_func_code_block {
		fmt.Fprintln(p, a.function_name, "\t", a.stage)
	}
	p.Flush()
	return all_func_code_block
}

func sort_the_func_block(slice_of_func_block []function) {
	for i := 0; i < len(slice_of_func_block); i++ {
		for j := i + 1; j < len(slice_of_func_block); j++ {
			if slice_of_func_block[i].function_name > slice_of_func_block[j].function_name {
				SWAP(slice_of_func_block, i, j)
			}
		}
	}
}

func sort_the_test_func_by_the_original_func() []function {
	golang_files = []string{}
	golang_test_files = []string{}
	read_folders(".")

	type file_content struct {
		file_name      string
		package_name   string
		file_lines     []string
		file_functions []function
	}

	var files_test []file_content
	var all_Benchmark_Example_Fuzz_Test []function
	var all_Function_Method []function
	var all_Function_Method_garbage []function
	for _, a := range golang_test_files {
		lines, _ := file_lines(a)
		functions, lines := return_all_func_code_block(lines)
		Benchmark_Example_Fuzz_Test, Function_Method_garbage := sperate_Benchmark_Example_Fuzz_Test_from_Function_Method(functions)
		all_Benchmark_Example_Fuzz_Test = append(all_Benchmark_Example_Fuzz_Test, Benchmark_Example_Fuzz_Test...)
		all_Function_Method_garbage = append(all_Function_Method_garbage, Function_Method_garbage...)
		all_Function_Method = append(all_Function_Method, Function_Method_garbage...)

		var file file_content
		file.file_name = a
		file.file_lines = lines
		files_test = append(files_test, file)
	}

	var package_name string
	var file_origenal []file_content
	for _, a := range golang_files {
		lines, _ := file_lines(a)
		package_name = return_package_name(lines)
		functions, lines := return_all_func_code_block(lines)
		Benchmark_Example_Fuzz_Test, Function_Method := sperate_Benchmark_Example_Fuzz_Test_from_Function_Method(functions)
		all_Benchmark_Example_Fuzz_Test = append(all_Benchmark_Example_Fuzz_Test, Benchmark_Example_Fuzz_Test...)
		all_Function_Method = append(all_Function_Method, Function_Method...)

		var file file_content
		file.file_name = a
		file.package_name = package_name
		file.file_lines = lines
		file.file_functions = Function_Method
		file_origenal = append(file_origenal, file)
	}

	for _, v1 := range file_origenal {
		test_file_name := return_test_file_name(v1.file_name)
		file_test, err := find(files_test, func(i int) bool {
			return test_file_name == files_test[i].file_name
		})
		if err != nil {
			file_test.file_lines = write_the_first_lines_of_test_file(v1.package_name)
		}
		sort_the_func_block(v1.file_functions)

		f, _ := os.Create(v1.file_name)
		write_lines(f, v1.file_lines)
		f_test, _ := os.Create(test_file_name)
		write_lines(f_test, file_test.file_lines)
		for _, v2 := range v1.file_functions {
			write_lines(f, v2.func_block)
			all_Benchmark_Example_Fuzz_Test = search_and_write_func(v2, all_Benchmark_Example_Fuzz_Test, f_test, "Benchmark_", write_func_benchmark)
			all_Benchmark_Example_Fuzz_Test = search_and_write_func(v2, all_Benchmark_Example_Fuzz_Test, f_test, "Example_", write_func_example)
			all_Benchmark_Example_Fuzz_Test = search_and_write_func(v2, all_Benchmark_Example_Fuzz_Test, f_test, "Fuzz_", write_func_fuzz)
			all_Benchmark_Example_Fuzz_Test = search_and_write_func(v2, all_Benchmark_Example_Fuzz_Test, f_test, "Test_", write_func_test)
		}
	}
	all_Benchmark_Example_Fuzz_Test = append(all_Benchmark_Example_Fuzz_Test, all_Function_Method_garbage...)
	sort_the_func_block(all_Benchmark_Example_Fuzz_Test)
	write_file("./garbage_test.go", package_name, all_Benchmark_Example_Fuzz_Test)

	return all_Function_Method
}

func sperate_Benchmark_Example_Fuzz_Test_from_Function_Method(functions []function) ([]function, []function) {
	var all_Function_Method_garbage []function
	var all_Benchmark_Example_Fuzz_Test []function
	for _, v1 := range functions {
		if v1.function_type == Function || v1.function_name == Method {
			all_Function_Method_garbage = append(all_Function_Method_garbage, v1)
		} else {
			all_Benchmark_Example_Fuzz_Test = append(all_Benchmark_Example_Fuzz_Test, v1)
		}
	}
	return all_Benchmark_Example_Fuzz_Test, all_Function_Method_garbage
}

func type_of_the_function(line string) (string, string, string, bool) {
	len_line := len(line)
	if len_line > len_func_ && line[:len_func_] == func_ {
		_, function_name, _, _, _ := analyse_the_function_line(line)
		switch {
		case len_line > len_func_+len_Benchmark && line[:len_func_+len_Benchmark] == func_+Benchmark:
			return "", Benchmark, function_name, true
		case len_line > len_func_+len_Example && line[:len_func_+len_Example] == func_+Example:
			return "", Example, function_name, true
		case len_line > len_func_+len_Fuzz && line[:len_func_+len_Fuzz] == func_+Fuzz:
			return "", Fuzz, function_name, true
		case len_line > len_func_+len_Test && line[:len_func_+len_Test] == func_+Test:
			return "", Test, function_name, true
		case len_line > len_func_+1 && line[:len_func_+1] == func_+"(":
			return "", Method, function_name, true
		default:
			return "", Function, function_name, true
		}
	}
	return "", "", "", false
}

func write_file(path, package_name string, all_functions []function) {
	lines_test, err := file_lines(path)
	_, lines_test = return_all_func_code_block(lines_test)
	f_test, _ := os.Create(path)
	if err != nil {
		lines_test = write_the_first_lines_of_test_file(package_name)
	}
	write_lines(f_test, lines_test)
	for _, a := range all_functions {
		write_lines(f_test, a.func_block)
	}
}

func write_func_benchmark(func_name string) []string {
	return []string{
		"func Benchmark_" + func_name + "(b *testing.B) {",
		"\t//TODO",
		"\t//for i := 0; i < b.N; i++ {",
		"\t//\t" + func_name + "()",
		"\t//}",
		"}",
	}
}

func write_func_example(func_name string) []string {
	return []string{
		"func Example_" + func_name + "() {",
		"\t//TODO",
		"\t//a1:=" + func_name + "()",
		"\t//fmt.Println(a1)",
		"}",
	}
}

func write_func_fuzz(func_name string) []string {
	return []string{
		"func Fuzz_" + func_name + "(f *testing.F) {",
		"\t//TODO",
		"\t//f.Add()",
		"\t//f.Fuzz(func(t *testing.T){",
		"\t//\ta1:=" + func_name + "()",
		"\t//\te1:=",
		"\t//\tTEST(true,a1,e1)",
		"\t//})",
		"}",
	}
}

func write_func_test(func_name string) []string {
	return []string{
		"func Test_" + func_name + "(t *testing.T) {",
		"\t//TODO",
		"\t//a1:=" + func_name + "()",
		"\t//e1:=",
		"\t//TEST(true,a1,e1)",
		"}",
	}
}

func write_lines(f *os.File, lines []string) {
	for _, a := range lines {
		f.WriteString("\n" + a)
	}
	f.WriteString("\n")
}

func write_test_tools(path string) {
	f, _ := os.Create(path + "/test_tools.txt")
	f.WriteString(
		`//this file created automatically
//please copy this function to your code to make it work
package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime/debug"
	"text/tabwriter"
)
 var fail_test_number int
 func TEST[t any](should_equal bool, actual, expected t) {
	if reflect.DeepEqual(actual, expected) != should_equal {
		fail_test_number++
		p := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintln(p, "\033[32m", "fail_test_number\t:", fail_test_number) //green
		fmt.Fprintln(p, "\033[35m", "should_equal\t:", should_equal)         //purple
		fmt.Fprintln(p, "\033[34m", "actual\t:", actual)                     //blue
		fmt.Fprintln(p, "\033[33m", "expected\t:", expected)                 //yellow
		p.Flush()
		fmt.Println("\033[31m") //red
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(string(debug.Stack()), "\033[0m") //reset
			}
		}()
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
