package main

import (
	"bufio"
	"fmt"
	"notebook-dsa/avl"
	"notebook-dsa/bst"
	"notebook-dsa/event"
	"notebook-dsa/hashtable"
	"os"
	"strconv"
	"strings"
	"time"
)

// --- Ввод ---

var reader = bufio.NewReader(os.Stdin)

func readLine() string {
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func readInt(prompt string, lo, hi int) int {
	for {
		fmt.Printf("%s [%d-%d]: ", prompt, lo, hi)
		s := readLine()
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if err == nil && v >= lo && v <= hi {
			return v
		}
		fmt.Printf("Ошибка: введите целое число от %d до %d\n", lo, hi)
	}
}

func readString(prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		s := readLine()
		s = strings.TrimSpace(s)
		if s != "" {
			return s
		}
		fmt.Println("Ошибка: поле не может быть пустым")
	}
}

func readEvent() event.Event {
	fmt.Println("\nНовое событие")
	var e event.Event
	for {
		e.Year = readInt("Год", 1900, 2100)
		e.Month = readInt("Месяц", 1, 12)
		e.Day = readInt("День", 1, event.DaysInMonth(e.Year, e.Month))
		if event.IsValidDate(e.Year, e.Month, e.Day) {
			break
		}
		fmt.Println("Некорректная дата, попробуйте снова")
	}
	e.Hour = readInt("Час", 0, 23)
	e.Minute = readInt("Минута", 0, 59)
	e.Weekday = event.CalcWeekday(e.Year, e.Month, e.Day)
	fmt.Printf("День недели: %s\n", event.WeekdayName(e.Weekday))
	e.Place = readString("Место")
	e.Desc = readString("Описание")
	return e
}

func readDateTime() (y, mo, d, h, mi int) {
	y = readInt("Год", 1900, 2100)
	mo = readInt("Месяц", 1, 12)
	d = readInt("День", 1, event.DaysInMonth(y, mo))
	h = readInt("Час", 0, 23)
	mi = readInt("Минута", 0, 59)
	return
}

// --- Сохранение / загрузка ---

func saveEvents(events []event.Event, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, e := range events {
		fmt.Fprintf(w, "%d %d %d %d %d\n", e.Year, e.Month, e.Day, e.Hour, e.Minute)
		fmt.Fprintln(w, e.Place)
		fmt.Fprintln(w, e.Desc)
	}
	return w.Flush()
}

func loadEvents(filename string) ([]event.Event, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var events []event.Event
	scanner := bufio.NewScanner(f)
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		var y, mo, d, h, mi int
		_, err := fmt.Sscanf(line, "%d %d %d %d %d", &y, &mo, &d, &h, &mi)
		if err != nil {
			continue
		}
		if !scanner.Scan() {
			break
		}
		place := scanner.Text()
		if !scanner.Scan() {
			break
		}
		desc := scanner.Text()

		if !event.IsValidDate(y, mo, d) || !event.IsValidTime(h, mi) {
			fmt.Println("Пропущена некорректная запись в файле")
			continue
		}
		events = append(events, event.NewEvent(y, mo, d, h, mi, desc, place))
	}
	return events, nil
}

// --- Меню BST ---

func menuBST() {
	tree := bst.New()
	for {
		fmt.Println("\nБинарное дерево поиска")
		fmt.Println("1 - Показать все события (по дате)")
		fmt.Println("2 - Добавить событие")
		fmt.Println("3 - Удалить событие")
		fmt.Println("4 - Фильтр по месту")
		fmt.Println("5 - Поиск по дате и времени")
		fmt.Println("6 - Сохранить в файл")
		fmt.Println("7 - Загрузить из файла")
		fmt.Println("0 - Назад")
		choice := readInt("Выбор", 0, 7)

		switch choice {
		case 1:
			cnt := tree.Count()
			if cnt == 0 {
				fmt.Println("Список пуст")
			} else {
				fmt.Printf("\n--- Все события (%d шт.), отсортированы по дате ---\n", cnt)
				tree.Inorder(event.Print)
			}
		case 2:
			e := readEvent()
			tree.Add(e)
			fmt.Println("Событие добавлено")
		case 3:
			if tree.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			fmt.Println("\nВведите дату и время удаляемого события:")
			y, mo, d, h, mi := readDateTime()
			before := tree.Count()
			tree.DeleteByDatetime(y, mo, d, h, mi)
			if tree.Count() < before {
				fmt.Println("Событие удалено")
			} else {
				fmt.Println("Событие не найдено")
			}
		case 4:
			if tree.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			sub := readString("Подстрока места")
			fmt.Printf("\n--- Фильтр по месту: \"%s\" ---\n", sub)
			found := false
			tree.FilterPlace(sub, func(e *event.Event) {
				event.Print(e)
				found = true
			})
			if !found {
				fmt.Println("Ничего не найдено")
			}
		case 5:
			if tree.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			fmt.Println("\nПоиск по дате и времени:")
			y, mo, d, h, mi := readDateTime()
			start := time.Now()
			res := tree.FindByDatetime(y, mo, d, h, mi)
			elapsed := time.Since(start).Nanoseconds()
			if len(res.Events) == 0 {
				fmt.Println("Событий не найдено")
			} else {
				for _, e := range res.Events {
					event.Print(e)
				}
			}
			fmt.Printf("Сравнений: %d | Время: %d нс | Высота дерева: %d\n",
				res.Ops, elapsed, tree.Height())
		case 6:
			filename := readString("Имя файла")
			if err := saveEvents(tree.Collect(), filename); err != nil {
				fmt.Println("Ошибка сохранения:", err)
			} else {
				fmt.Printf("Сохранено (%d событий)\n", tree.Count())
			}
		case 7:
			filename := readString("Имя файла")
			events, err := loadEvents(filename)
			if err != nil {
				fmt.Println("Ошибка загрузки:", err)
				continue
			}
			tree = bst.New()
			for i := range events {
				tree.Add(events[i])
			}
			fmt.Printf("Загружено (%d событий)\n", tree.Count())
		case 0:
			return
		}
	}
}

// --- Меню AVL ---

func menuAVL() {
	tree := avl.New()
	for {
		fmt.Println("\nАВЛ-дерево")
		fmt.Println("1 - Показать все события (по дате)")
		fmt.Println("2 - Добавить событие")
		fmt.Println("3 - Удалить событие")
		fmt.Println("4 - Фильтр по месту")
		fmt.Println("5 - Поиск по дате и времени")
		fmt.Println("6 - Сохранить в файл")
		fmt.Println("7 - Загрузить из файла")
		fmt.Println("0 - Назад")
		choice := readInt("Выбор", 0, 7)

		switch choice {
		case 1:
			cnt := tree.Count()
			if cnt == 0 {
				fmt.Println("Список пуст")
			} else {
				fmt.Printf("\n--- Все события (%d шт.), отсортированы по дате ---\n", cnt)
				tree.Inorder(event.Print)
			}
		case 2:
			e := readEvent()
			tree.Add(e)
			fmt.Println("Событие добавлено")
		case 3:
			if tree.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			fmt.Println("\nВведите дату и время удаляемого события:")
			y, mo, d, h, mi := readDateTime()
			before := tree.Count()
			tree.DeleteByDatetime(y, mo, d, h, mi)
			if tree.Count() < before {
				fmt.Println("Событие удалено")
			} else {
				fmt.Println("Событие не найдено")
			}
		case 4:
			if tree.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			sub := readString("Подстрока места")
			fmt.Printf("\n--- Фильтр по месту: \"%s\" ---\n", sub)
			found := false
			tree.FilterPlace(sub, func(e *event.Event) {
				event.Print(e)
				found = true
			})
			if !found {
				fmt.Println("Ничего не найдено")
			}
		case 5:
			if tree.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			fmt.Println("\nПоиск по дате и времени:")
			y, mo, d, h, mi := readDateTime()
			start := time.Now()
			res := tree.FindByDatetime(y, mo, d, h, mi)
			elapsed := time.Since(start).Nanoseconds()
			if len(res.Events) == 0 {
				fmt.Println("Событий не найдено")
			} else {
				for _, e := range res.Events {
					event.Print(e)
				}
			}
			fmt.Printf("Сравнений: %d | Время: %d нс | Высота дерева: %d\n",
				res.Ops, elapsed, tree.Height())
		case 6:
			filename := readString("Имя файла")
			if err := saveEvents(tree.Collect(), filename); err != nil {
				fmt.Println("Ошибка сохранения:", err)
			} else {
				fmt.Printf("Сохранено (%d событий)\n", tree.Count())
			}
		case 7:
			filename := readString("Имя файла")
			events, err := loadEvents(filename)
			if err != nil {
				fmt.Println("Ошибка загрузки:", err)
				continue
			}
			tree = avl.New()
			for i := range events {
				tree.Add(events[i])
			}
			fmt.Printf("Загружено (%d событий)\n", tree.Count())
		case 0:
			return
		}
	}
}

// --- Меню HashTable ---

func menuHash() {
	ht := hashtable.New()
	for {
		fmt.Println("\nХеш-таблица")
		fmt.Println("1 - Показать все события (по дате)")
		fmt.Println("2 - Добавить событие")
		fmt.Println("3 - Удалить событие")
		fmt.Println("4 - Фильтр по месту")
		fmt.Println("5 - Поиск по дате и времени")
		fmt.Println("6 - Сохранить в файл")
		fmt.Println("7 - Загрузить из файла")
		fmt.Println("0 - Назад")
		choice := readInt("Выбор", 0, 7)

		switch choice {
		case 1:
			cnt := ht.Count()
			if cnt == 0 {
				fmt.Println("Список пуст")
			} else {
				events := ht.AllSorted()
				fmt.Printf("\n--- Все события (%d шт.), отсортированы по дате ---\n", cnt)
				for i := range events {
					event.Print(&events[i])
				}
			}
		case 2:
			e := readEvent()
			ht.Add(e)
			fmt.Println("Событие добавлено")
		case 3:
			if ht.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			fmt.Println("\nВведите дату и время удаляемого события:")
			y, mo, d, h, mi := readDateTime()
			key := event.Event{Year: y, Month: mo, Day: d, Hour: h, Minute: mi}
			if ht.Delete(&key) {
				fmt.Println("Событие удалено")
			} else {
				fmt.Println("Событие не найдено")
			}
		case 4:
			if ht.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			sub := readString("Подстрока места")
			fmt.Printf("\n--- Фильтр по месту: \"%s\" ---\n", sub)
			events := ht.FilterPlace(sub)
			if len(events) == 0 {
				fmt.Println("Ничего не найдено")
			} else {
				for i := range events {
					event.Print(&events[i])
				}
			}
		case 5:
			if ht.Count() == 0 {
				fmt.Println("Список пуст")
				continue
			}
			fmt.Println("\nПоиск по дате и времени:")
			y, mo, d, h, mi := readDateTime()
			start := time.Now()
			res := ht.FindByDatetime(y, mo, d, h, mi)
			elapsed := time.Since(start).Nanoseconds()
			if len(res.Events) == 0 {
				fmt.Println("Событий не найдено")
			} else {
				for _, e := range res.Events {
					event.Print(e)
				}
			}
			fmt.Printf("Сравнений: %d | Время: %d нс\n", res.Ops, elapsed)
		case 6:
			filename := readString("Имя файла")
			if err := saveEvents(ht.AllSorted(), filename); err != nil {
				fmt.Println("Ошибка сохранения:", err)
			} else {
				fmt.Printf("Сохранено (%d событий)\n", ht.Count())
			}
		case 7:
			filename := readString("Имя файла")
			events, err := loadEvents(filename)
			if err != nil {
				fmt.Println("Ошибка загрузки:", err)
				continue
			}
			ht = hashtable.New()
			for i := range events {
				ht.Add(events[i])
			}
			fmt.Printf("Загружено (%d событий)\n", ht.Count())
		case 0:
			return
		}
	}
}

// --- Главное меню ---

func main() {
	for {
		fmt.Println("\nЗаписная книжка - DSA")
		fmt.Println("1 - BST (бинарное дерево поиска)")
		fmt.Println("2 - AVL (сбалансированное дерево)")
		fmt.Println("3 - HashTable (хеш-таблица)")
		fmt.Println("0 - Выход")
		choice := readInt("Выбор", 0, 3)

		switch choice {
		case 1:
			menuBST()
		case 2:
			menuAVL()
		case 3:
			menuHash()
		case 0:
			fmt.Println("Программа завершена")
			return
		}
	}
}
