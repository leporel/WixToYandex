package cmd

import (
	"fmt"
	config2 "github.com/leporel/wixtoyandex/config"
	"github.com/leporel/wixtoyandex/converter"
	"github.com/spf13/cobra"
	"os"
)

var version = "unknown"

var (
	rootCmd = &cobra.Command{}

	inFolder   string
	inFile     string
	configFile string
)

func init() {
	rootCmd = &cobra.Command{
		Short: "Конвертация csv выгрузки товаров Wix в формат для шаблона Yandex Market'a",
		Long: `                Конвертация csv выгрузки товаров Wix в формат для шаблона Yandex Market'a
	Создает аналогичный файл csv, от куда можно скопировать товары в шаблон для яндекс маркета.
	Подробнее по ссылке https://github.com/leporel/wixtoyandex
	Для справки запустите wty --help`,
		Example: "wty.exe -f C:\\catalog_products_2.csv",
		Run: func(c *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Printf("Программа запущена без аргументов, все файлы в папке %s будут преобразованы и сохранены в туже папку\n\n", inFolder)
			}

			config2.InitConfig(configFile)

			if inFile != "" {
				if err := converter.ConvertFile(inFile); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				if err := converter.ConvertFiles(inFolder); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		},
		Version: version,
	}

	rootCmd.PersistentFlags().StringVarP(&inFile, "file", "f", "", "конвертировать один файл, преобразованный результат будет в той же директории")
	rootCmd.PersistentFlags().StringVarP(&inFolder, "inputFolder", "i", "wix", `папка с csv файлами wix`)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./config.toml", "путь к файлу конфигурации")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
