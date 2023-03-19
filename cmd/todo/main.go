package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"todo"
)

const (
	todoFile = ".todos.json"
)

func main() {

	add := flag.Bool("add", false, "add a new todo")

	complete := flag.Int("complete", 0, "mark a todo as completed")

	delete := flag.Int("delete", 0, "deletes a todo from the file")

	list := flag.Bool("list", false, "prints list of todos")

	flag.Parse()

	todos := &todo.Todos{}

	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		checkError(err)
		todos.Add(task)
		err = todos.Store(todoFile)
		checkError(err)
	case *complete > 0:
		err := todos.Complete(*complete)
		checkError(err)

		err = todos.Store(todoFile)
		checkError(err)
	case *delete > 0:
		err := todos.Delete(*delete)
		checkError(err)

		err = todos.Store(todoFile)
		checkError(err)
	case *list:
		todos.Print()
	default:
		fmt.Fprintln(os.Stdout, "Invalid command")
		os.Exit(0)
	}

}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := strings.TrimSpace(scanner.Text())

	if len(text) < 1 {
		return "", errors.New("todo cannot be empty")
	}

	return text, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
