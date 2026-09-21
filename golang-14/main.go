package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ListBooks struct {
	Author   string
	NameBook string
}

type Books [5]ListBooks

func (b Books) createFile(n int, name string) error {
	filename := fmt.Sprintf("student_%d.txt", n)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o664)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer f.Close()

	t := time.Now().Format(time.TimeOnly)
	_, err = f.WriteString(fmt.Sprintf("%s\n 478\n %d\n %s\n", name, n, t))
	if err != nil {
		return err
	}

	for i, v := range b {
		_, err = f.WriteString(fmt.Sprintf("%d. %s - %s\n", i+1, v.Author, v.NameBook))
		if err != nil {
			return err
		}
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	var count int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	_, err = f.WriteString(fmt.Sprintf("Количество строк: %d", count))
	if err != nil {
		return err
	}

	return nil
}

func createDirs(n int, paths map[string][]string) error {
	base := filepath.Join(fmt.Sprintf("project_%d", n))
	for k, v := range paths {
		for _, v2 := range v {
			path := filepath.Join(base, k, v2)
			if err := os.MkdirAll(path, 0o775); err != nil {
				return err
			}
		}
	}
	return nil
}

func visit(p string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fmt.Println(p)
	return nil
}

func createFilesInfo(n int, desc map[string]string) error {
	err := filepath.WalkDir(fmt.Sprintf("project_%d", n), func(p string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			dirname := filepath.Base(p)
			if descr, ok := desc[dirname]; ok {
				d := []byte(descr)
				infop := filepath.Join(p, "info.txt")
				err := os.WriteFile(infop, d, 0o775)
				if err != nil {
					return err
				}
				t := []byte(time.Now().Format("02-01-2006 15:04:05"))
				timemd := filepath.Join(p, "README.md")
				err = os.WriteFile(timemd, t, 0o775)
				if err != nil {
					return err
				}
			}
		}
		fmt.Println(p)
		return nil
	})
	return err
}

type statis struct {
	CountFiles  int
	CountDirs   int
	CountSizeKB float64
	CountSizeMB float64
	ListExt     map[string]int
	TopFiveBig  []string
}

type fileInfo struct {
	Path string
	Size int64
}

func selectDir(path *string) (string, error) {
	str, err := filepath.Abs(*path)
	if err != nil {
		return "", err
	}
	return str, nil
}

func collectFiles(root string) ([]fileInfo, error) {
	var files []fileInfo

	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintln(os.Stderr, "пропуск:", err)
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		files = append(files, fileInfo{Path: p, Size: info.Size()})
		return nil
	})
	return files, err
}

func (s statis) statisFunc(path string) (statis, error) {
	ss := statis{ListExt: make(map[string]int)}
	wanted := []string{".js", ".json", ".txt", ".md"}
	var totalBytes int64

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			ss.CountDirs++
			return nil
		}

		ss.CountFiles++
		if ext := filepath.Ext(p); slices.Contains(wanted, ext) {
			ss.ListExt[ext]++
		}

		fi, err := d.Info()
		if err != nil {
			return err
		}
		totalBytes += fi.Size()
		return nil
	})
	if err != nil {
		return statis{}, err
	}

	ss.CountSizeKB = float64(totalBytes) / 1024
	ss.CountSizeMB = float64(totalBytes) / 1024 / 1024

	files, err := collectFiles(path)
	if err != nil {
		return statis{}, err
	}
	slices.SortFunc(files, func(a, b fileInfo) int {
		return cmp.Compare(b.Size, a.Size)
	})
	for i, f := range files[:min(5, len(files))] {
		mb := float64(f.Size) / 1024 / 1024
		ss.TopFiveBig = append(ss.TopFiveBig, fmt.Sprintf("%d: %s - %.2f мб", i+1, f.Path, mb))
	}

	return ss, nil
}

func createFile100k() error {
	file, err := os.OpenFile("data_14.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o775)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < 100000; i++ {
		randn := rand.Intn(1001-1) + 1
		s := fmt.Sprintf("%d, %d, Вариант 14\n", i, randn)
		if _, err := writer.WriteString(s); err != nil {
			return err
		}
	}
	return nil
}

/*
*Сумму всех чисел
Среднее арифметическое
Максимальное и минимальное число
Количество строк
*
*/
func readFile100k() error {
	file, err := os.Open("data_14.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	fileresult, err := os.OpenFile("processed_14.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o642)
	if err != nil {
		return err
	}
	defer fileresult.Close()
	filedop, err := os.OpenFile("dop_14.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o642)
	if err != nil {
		return err
	}
	defer filedop.Close()
	scanner_dop := bufio.NewWriter(filedop)
	defer scanner_dop.Flush()

	var sum int64 = 0
	var avg float64 = 0
	var count int
	var max_1 int64 = 0
	var min_1 int64 = math.MaxInt64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
		if int64(count)%10000 == 0 {
			fmt.Printf("Прогресс: %d\n", int64(count)/1000)
		}

		str := strings.Split(scanner.Text(), ",") // [1],[10],[Вариант 14]

		firstn, _ := strconv.Atoi(strings.TrimSpace(str[0]))
		secondn, _ := strconv.Atoi(strings.TrimSpace(str[1]))
		if int64(secondn) > max_1 {
			max_1 = int64(secondn)
		}
		if int64(secondn) < min_1 {
			min_1 = int64(secondn)
		}
		if firstn > 500 && secondn > 500 {
			scanner_dop.Write([]byte(scanner.Text() + "\n"))
		}
		sum += int64(firstn) + int64(secondn)
	}
	avg = float64(sum) / float64(count*2)
	sss := fmt.Sprintf("sum = %d avg = %.2f count = %d max = %d min = %d\n", sum, avg, count, max_1, min_1)
	fmt.Print(sss)
	fileresult.Write([]byte(sss))
	return nil
}

func main() {
	path := flag.String("p", ".", "принимает путь к директории через аргументы командной строки")
	del := flag.Bool("d", false, "удаляет project_14")
	flag.Parse()
	create := Books{
		{Author: "Петя", NameBook: "Основы Go"},
		{Author: "Миша", NameBook: "GOD OF AI"},
		{Author: "Джордж Оруэлл", NameBook: "1984"},
		{Author: "Роберт Мартин", NameBook: "Чистый код"},
		{Author: "Адитья Бхаргава", NameBook: "Грокаем алгоритмы"},
	}
	create.createFile(14, "Мишкевич Максим")

	dirs := map[string][]string{
		"src":  {"modules", "components", "utils"},
		"data": {"input", "output"},
		"":     {"temp"},
	}

	desc := map[string]string{
		"src":        "Исходный код",
		"modules":    "Крупные модули",
		"components": "Компоненты интерфейса",
		"utils":      "Вспомогательные функции",
		"data":       "Данные проекта",
		"input":      "Входные данные",
		"output":     "Выходные результаты",
		"temp":       "Временные файлы",
	}

	if err := createDirs(14, dirs); err != nil {
		fmt.Println(err)
		return
	}

	if *del {
		err := os.RemoveAll("project_14/")
		fmt.Println("Папка project_14 удалена")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := createFilesInfo(14, desc); err != nil {
		fmt.Println(err)
		return
	}

	err := os.Rename("project_14/temp", "project_14/data/temp")
	if err != nil {
		fmt.Println(err)
		return
	}

	delete(desc, "")
	dirs["data"] = []string{"input", "output", "temp"}

	err = os.Rename("project_14/data/output", "project_14/data/result")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 45; i++ {
		fmt.Print("-")
	}
	fmt.Println()

	err = filepath.Walk("project_14", visit)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	var str string
	str, err = selectDir(path)
	if err != nil {
		return
	}

	err = os.RemoveAll("project_14/data/temp")
	fmt.Println("Папка temp удалена")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = filepath.Walk(str, visit)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	for i := 0; i < 45; i++ {
		fmt.Print("-")
	}
	fmt.Println()

	f1, err := os.OpenFile("report_14.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o764)
	if err != nil {
		fmt.Println(err)
		return
	}

	var s statis
	s, err = s.statisFunc(str)
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("Путь: %s\nКоличество файлов: %d\nКоличество папок: %d\nРазмер: %.2f mb %.2f kb\n", str, s.CountFiles, s.CountDirs, s.CountSizeMB, s.CountSizeKB))
	f1.WriteString(fmt.Sprintf("Путь: %s\nКоличество файлов: %d\nКоличество папок: %d\nРазмер: %.2f mb %.2f kb\n", str, s.CountFiles, s.CountDirs, s.CountSizeMB, s.CountSizeKB))
	for k, v := range s.ListExt {
		fmt.Printf("%s - %d раз(а)\n", k, v)
		f1.WriteString(fmt.Sprintf("%s - %d раз(а)\n", k, v))
	}
	fmt.Println()
	for _, v := range s.TopFiveBig {
		fmt.Println(v)
		f1.WriteString(v + "\n")
	}

	err = createFile100k()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = readFile100k()
	if err != nil {
		fmt.Println(err)
		return
	}
}
