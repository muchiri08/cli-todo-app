package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}
type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}

func (t *Todos) Delete(index int) error {
	ls := *t

	if index < 1 || index > len(ls) {
		return errors.New("invalid index")
	}

	*t = append(ls[:index-1], ls[index:]...)

	return nil
}

func (t *Todos) Load(filename string) error {

	_, err := os.Stat(filename)

	//check if file exist. if not create it
	if os.IsNotExist(err) {
		file, err := os.Create(filename)

		if err != nil {
			return err
		}
		defer file.Close()
	}

	file, err := os.ReadFile(filename)

	if err != nil {
		return err
	}
	if len(file) == 0 {
		return err
	}
	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignCenter, Text: "CreatedAt"},
			{Align: simpletable.AlignCenter, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	for index, item := range *t {
		index++

		id := blue(fmt.Sprintf("%d", index))
		task := blue(item.Task)
		done := blue(fmt.Sprintf("%t", item.Done))
		createdAt := blue(item.CreatedAt.Format(time.RFC822))
		completedAt := blue(item.CompletedAt.Format(time.RFC822))

		if item.Done {
			id = green(fmt.Sprintf("%d", index))
			task = green(fmt.Sprintf("\u2713 %s", item.Task))
			done = green(fmt.Sprintf("%t", item.Done))
			createdAt = green(item.CreatedAt.Format(time.RFC822))
			completedAt = green(item.CompletedAt.Format(time.RFC822))
		}

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: id},
			{Text: task},
			{Text: done},
			{Text: createdAt},
			{Text: completedAt},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("You have %d pending todos!", t.countPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) countPending() int {
	totals := 0

	for _, item := range *t {
		if !item.Done {
			totals++
		}
	}
	return totals
}
