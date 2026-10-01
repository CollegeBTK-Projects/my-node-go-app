package main

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const group = "478"

func resources() {
	cores, _ := cpu.Counts(true)
	infos, _ := cpu.Info()

	fmt.Println("Количество логических ядер:", cores)
	fmt.Println("Модель процессора:", infos[0].ModelName)

	sum := 0.0
	for i, c := range infos {
		fmt.Printf("Ядро %d: %.0f МГц\n", i, c.Mhz)
		sum += c.Mhz
	}
	fmt.Printf("Средняя частота: %.0f МГц\n", sum/float64(len(infos)))

	m, _ := mem.VirtualMemory()
	gb := 1024.0 * 1024 * 1024
	total := float64(m.Total) / gb
	free := float64(m.Available) / gb
	used := total - free

	fmt.Printf("Общий объём: %.2f ГБ\n", total)
	fmt.Printf("Свободно: %.2f ГБ\n", free)
	fmt.Printf("Использовано: %.2f ГБ (%.1f%%)\n", used, used/total*100)

	if runtime.GOOS == "windows" {
		fmt.Println("Недоступно в Windows")
	} else {
		l, _ := load.Avg()
		fmt.Println("1, 5, 15 минут:", []float64{l.Load1, l.Load5, l.Load15})
	}

	fmt.Println()
	fmt.Println("Группа:", group)
	if free/total*100 < 20 {
		fmt.Println("<20%!")
	} else {
		fmt.Println("Память в норме")
	}
}

func network() {
	ifaces, _ := net.Interfaces()
	mainIface := ""

	for _, i := range ifaces {
		fmt.Println("Интерфейс:", i.Name)
		addrs, _ := i.Addrs()
		for _, a := range addrs {
			ip, _, _ := net.ParseCIDR(a.String())
			if ip.To4() != nil {
				fmt.Println("IPv4:", ip)
				if mainIface == "" && !ip.IsLoopback() {
					mainIface = i.Name + " (" + ip.String() + ")"
				}
			} else {
				fmt.Println("IPv6:", ip)
			}
		}
		mac := "нет"
		if p := strings.Split(i.HardwareAddr.String(), ":"); len(p) == 6 {
			mac = strings.ToUpper(strings.Join(p[:3], ":")) + ":**:**:**"
		}
		fmt.Println("MAC:", mac)
		fmt.Println("Внутренний:", map[bool]string{true: "да", false: "нет"}[i.Flags&net.FlagLoopback != 0])
		fmt.Println()
	}
	fmt.Println("Всего интерфейсов:", len(ifaces))
	fmt.Println("Основной интерфейс:", mainIface)

	u, _ := user.Current()
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "недоступно"
	}

	fmt.Println("Имя пользователя:", u.Username)
	fmt.Println("UID:", u.Uid)
	fmt.Println("GID:", u.Gid)
	fmt.Println("Домашняя директория:", u.HomeDir)
	fmt.Println("Оболочка:", shell)
	fmt.Println("Группа:", group)
	fmt.Println("Проверка root:", map[bool]string{true: "Да", false: "нет"}[os.Geteuid() == 0])
}

func basicInfo() {
	version, _ := host.Info()
	hostname, _ := os.Hostname()

	names := map[string]string{
		"windows": "Windows_NT",
		"linux":   "Linux",
		"darwin":  "Darwin",
	}

	ost := names[runtime.GOOS]

	fmt.Println("ОС: ", runtime.GOOS)
	fmt.Println("Архитектура: ", runtime.GOARCH)
	fmt.Println("Имя ОС: ", ost)
	fmt.Println("Версия ОС: ", version.KernelVersion)
	fmt.Println("Имя хоста: ", hostname)
	fmt.Println("Аптайм: ", version.Uptime)

	switch runtime.GOOS {
	case "windows":
		fmt.Println("Вы работаете в Windows")
	case "linux":
		fmt.Println("Вы работаете в Linux")
	case "darwin":
		fmt.Println("Вы работаете в macOS")
	default:
		fmt.Println("Неизвестная платформа")
	}
}

func main() {
	resources()
	network()
	basicInfo()
}
