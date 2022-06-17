package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitors = 3
const delay = 3

func main() {
	for {
		showMenu()
		command := readCommand()

		switch command {
		case 1:
			fmt.Println("monitorando ...")
			startMonitor()
		case 2:
			fmt.Println("Exibindo logs")
			readLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Escolha não identificada")
			os.Exit(-1)
		}
	}
}

func showMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do programa")
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	return command
}

func startMonitor() {
	sites := readSites()

	for i := 0; i < monitors; i++ {
		for _, site := range sites {
			checkSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
}

func checkSite(site string) {
	rsp, err := http.Get(site)

	if err != nil {
		fmt.Println("Error:", err)
	}
	if rsp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		writeLogs(site, true)
	} else {
		fmt.Println("Site:", site, "está com problema!")
		writeLogs(site, false)
	}
}

func readSites() []string {
	sites := []string{}
	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Error:", err)
	}

	reader := bufio.NewReader(file)

	for {
		row, err := reader.ReadString('\n')
		row = strings.TrimSpace(row)

		sites = append(sites, row)

		if err == io.EOF {
			break
		}
	}

	file.Close()
	return sites
}

func writeLogs(site string, status bool) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}
	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	file.Close()
}

func readLogs() {
	file, err := ioutil.ReadFile("logs.txt")

	if err != nil {
		fmt.Println("sem logs no momento")
	}

	fmt.Println(string(file))
}
