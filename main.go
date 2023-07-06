package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const delay = 3 * time.Second

func main() {

	exibeIntroducao()

	for {
		exibeMenu()
		comando := leComando()
		switch comando {
		case 1:
			iniciarMonitoramento()
			time.Sleep(delay)
		case 2:
			leArquivoLogs()
			time.Sleep(delay)
		case 0:
			fmt.Println("Saindo do programa...")
			os.Remove("log.txt")
			fmt.Println("Obrigado!")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando, tente novamente.")
		}
	}
}

func exibeIntroducao() {
	nome := "Vinicius"
	versao := 0.9

	fmt.Println("Olá, Srx.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("--------------------------------------------------------------")
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
	fmt.Println("--------------------------------------------------------------")
}

func leComando() int {
	var comandoLido int

	fmt.Print("Digite uma opção: ")
	fmt.Scan(&comandoLido)
	return comandoLido
}

func iniciarMonitoramento() {
	sites := criarSlice()

	fmt.Println("Iniciando Monitoramento...")

	for i := 0; i < len(sites); i++ {
		resp, err := http.Get(sites[i])
		retornaErro(err)
		if resp.StatusCode == 200 {
			fmt.Println("O site:", sites[i], "foi carregado com sucesso!")
			registraLog(sites[i], true)
		} else {
			fmt.Println("O site:", sites[i], "está com problemas. ", "Status Code: ", resp.StatusCode)
			registraLog(sites[i], false)
		}
	}
	fmt.Println("Monitoramento realizado com sucesso!")
}

func criarSlice() []string {
	var comando string
	var sites []string

	fmt.Print("Você deseja escolher os sites para monitorar? (S/n): ")
	fmt.Scan(&comando)
	comando = strings.ToLower(comando)

	if comando == "s" {
		var quantidade int

		fmt.Print("Quantos sites você deseja escolher?: ")
		fmt.Scan(&quantidade)

		for i := 0; i < quantidade; i++ {
			var site string

			fmt.Printf("Digite o nome do site[%d]: ", i+1)
			fmt.Scan(&site)
			if !strings.Contains(site, "https://") {
				site = "https://" + site
			}

			sites = append(sites, site)
		}
	} else {
		fmt.Println("Será escolhido sites padrões...")
		sites = append(sites, "https://g1.globo.com/")
		sites = append(sites, "https://www.bbc.com/")
		sites = append(sites, "https://edition.cnn.com/")
		sites = append(sites, "https://www.investing.com/")
		sites = append(sites, "https://www.infomoney.com.br/")
	}
	return sites
}

func leArquivoLogs() {
	fmt.Println("Lendo logs...")
	arquivo, err := os.Open("log.txt")
	retornaErro(err)
	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		fmt.Println(linha)
		if err == io.EOF {
			break
		}
	}
	arquivo.Close()
}

func retornaErro(err error) {
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
}

func registraLog(site string, status bool) {
	_, exist := os.Stat("log.txt")
	arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if exist == nil {
		arquivo.WriteString("\n")
	}
	retornaErro(err)
	arquivo.WriteString("Data: " + time.Now().Format("02/01/06 15:04:05") +
		" - Site: " + site + " - Online: " + strconv.FormatBool(status))
	arquivo.Close()
}
