package converter

import (
	"encoding/csv"
	"fmt"
	"github.com/leporel/wixtoyandex/config"
	"github.com/microcosm-cc/bluemonday"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"pkg.re/essentialkaos/translit.v2"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

var checkUrlLock = make(chan interface{}, 10)

func ConvertFile(inFile string) error {
	params := config.Cfg.ConvertParams
	processDefaultRow(params)

	if filepath.Ext(inFile) != ".csv" {
		return fmt.Errorf("файл неправильного формата %s, ожидается файл формата .csv", filepath.Ext(inFile))
	}

	return convertFile(inFile)
}

func convertFile(inFile string) error {
	fmt.Println("Парсинг файла", inFile)

	wixFile, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла %s: %s", inFile, err)
	}
	defer wixFile.Close()

	dir := filepath.Dir(inFile)
	//dir := path.Dir(inFile)
	if !filepath.IsAbs(dir) {
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("ошибка при создании файла, не удалось получить текущию директорию: %s", err)
		}
	}

	outFile := filepath.Clean(filepath.Join(dir, strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile))+"_yandex.csv"))

	fmt.Println("Новый файл", outFile)
	yandexFile, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла %s: %s", outFile, err)
	}
	defer yandexFile.Close()

	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	if _, err := yandexFile.Write(bomUtf8); err != nil {
		return fmt.Errorf("ошибка при записи в файл кодировки: %s", err)
	}

	err = convert(wixFile, yandexFile)
	if err != nil {
		return err
	}
	fmt.Println("ГОТОВО!")

	return nil
}

func ConvertFiles(inFolder string) error {
	params := config.Cfg.ConvertParams
	processDefaultRow(params)

	if !filepath.IsAbs(inFolder) {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("не удалось получить текущию директорию: %s", err)
		}
		inFolder = filepath.Join(dir, inFolder)
	}

	if _, err := os.Stat(inFolder); os.IsNotExist(err) {
		return fmt.Errorf("директория не найдена: %s", inFolder)
	}

	err := filepath.Walk(inFolder, func(pathF string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("ошибка доступа к папке %q: %v\n", pathF, err)
			return err
		}
		if f.IsDir() {
			return nil
		}
		if filepath.Ext(f.Name()) == ".csv" {
			err = convertFile(pathF)
		}
		return err
	})
	if err != nil {
		return fmt.Errorf("ошибка при переборе файлов: %s", err)
	}

	fmt.Println("ГОТОВО!")

	return nil
}

func processDefaultRow(params config.ConvertParams) {
	if params.Delivery {
		defaultYandexRow[yandexColNumber["Доставка"]] = "Есть"
	}

	if params.DeliveryCost >= 0 {
		defaultYandexRow[yandexColNumber["Стоимость доставки"]] = fmt.Sprintf("%d", params.DeliveryCost)
	}

	defaultYandexRow[yandexColNumber["Срок доставки"]] = params.DeliveryTime

	if params.DeliveryOnly {
		defaultYandexRow[yandexColNumber["Самовывоз"]] = "Нет"
	}

	if params.NeedOrder {
		defaultYandexRow[yandexColNumber["Купить в магазине без заказа"]] = "Нельзя"
	}

	defaultYandexRow[yandexColNumber["Валюта"]] = params.Currency

	defaultYandexRow[yandexColNumber["Ссылка на товар на сайте магазина"]] = strings.TrimSuffix(params.Url, "/")

	defaultYandexRow[yandexColNumber["Ссылка на картинку"]] = params.WixUrl

	if params.Warranty {
		defaultYandexRow[yandexColNumber["Гарантия производителя"]] = "Есть"
	}
}

func convert(wixFile io.Reader, yandexFile io.Writer) error {
	params := config.Cfg.ConvertParams

	reader := csv.NewReader(wixFile)
	writer := csv.NewWriter(yandexFile)
	//writer.UseCRLF = true
	writer.Comma, _ = utf8.DecodeRune([]byte(params.Delimiter))

	write := func(record []string) error {
		err := writer.Write(record)
		if err != nil {
			return fmt.Errorf("ошибка при записи в файл: %s", err)
		}
		writer.Flush()
		err = writer.Error()
		if err != nil {
			return fmt.Errorf("ошибка при записи в файл: %s", err)
		}
		return nil
	}

	var existingUrl = make(map[string]string)

	line := 0

	var wg sync.WaitGroup

	for {
		line++

		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				wg.Wait()
				return nil
			}
			return fmt.Errorf("ошибка при парсинге файла: %s", err)
		}

		if line == 1 { // first line
			for i, s := range record {
				wixProductFields[i] = strings.TrimSpace(filterCharacters(s, true, unicode.Cyrillic))
			}
			continue
		}

		product := make(map[string]string)

		for i, s := range record {
			nameField, ok := wixProductFields[i]
			if !ok {
				return fmt.Errorf("fuck this stupid wix format file: %v, %v \n", i, s)
			}
			product[nameField] = filterHtml(filterCharacters(s, true, unicode.Cyrillic))
		}

		if product["visible"] != "true" {
			continue
		}

		newRow := makeRow(product)

		productUrl := newRow[yandexColNumber["Ссылка на товар на сайте магазина"]]
		sku := newRow[yandexColNumber["id"]]

		if v, ok := existingUrl[productUrl]; ok {
			fmt.Printf("	[WARNING] дубль ссылок [SKU:%s и %s]: %s\n", sku, v, productUrl)
		} else if params.Url != "" {
			wg.Add(1)

			go func() {
				defer wg.Done()
				ok, err = checkUrl(productUrl)
				if err != nil {
					fmt.Printf("[ERROR] не удалось проверить ссылку %s\n", productUrl)
					return
				}
				if !ok {
					fmt.Printf("[WARNING] нерабочая ссылка [SKU:%s]: %s\n", sku, productUrl)
				}

			}()
		}

		existingUrl[newRow[yandexColNumber["Ссылка на товар на сайте магазина"]]] = sku

		if err = write(newRow); err != nil {
			return err
		}
	}
}

func makeRow(product map[string]string) []string {
	rs := make([]string, len(defaultYandexRow))
	copy(rs, defaultYandexRow)

	rs[yandexColNumber["id"]] = filterSKU(translit.EncodeToISO9B(product["sku"]))
	if inventory, err := strconv.Atoi(product["inventory"]); err == nil && inventory == 0 {
		rs[yandexColNumber["Статус товара"]] = "На заказ"
	}

	// https://www.wix-site.ru/product-page/moy-product
	nr := strings.NewReplacer("`", "", "!", "", ".", "-", " ", "-", "/", "-", ",", "-")
	url := nr.Replace(translit.EncodeToISO9B(strings.ToLower(product["name"])))
	url = strings.TrimRight(url, "-")
	rs[yandexColNumber["Ссылка на товар на сайте магазина"]] = rs[yandexColNumber["Ссылка на товар на сайте магазина"]] + "/product-page/" + url

	// https://static.wixstatic.com/media/5e777_f2ded1e08ea4699d041d7cd9c17e8~mv2.png
	if rs[yandexColNumber["Ссылка на картинку"]] != "" {
		rs[yandexColNumber["Ссылка на картинку"]] = rs[yandexColNumber["Ссылка на картинку"]] + strings.Split(product["productImageUrl"], ";")[0]
	}

	rs[yandexColNumber["Название"]] = product["name"]

	rs[yandexColNumber["Категория"]] = product["collection"]

	if product["discountValue"] != "0.0" {
		price, _ := strconv.Atoi(product["price"])
		discount, _ := strconv.Atoi(product["discountValue"])
		price = price - (price * discount / 100)

		rs[yandexColNumber["Цена"]] = fmt.Sprint(price)
		rs[yandexColNumber["Цена без скидки"]] = strings.ReplaceAll(product["price"], ".", ",")
	}

	rs[yandexColNumber["Цена"]] = strings.ReplaceAll(product["price"], ".", ",")

	rs[yandexColNumber["Описание"]] = product["description"]

	listOptions := ""
	paramExist := true
	for i := 1; paramExist; i++ {
		if name, ok := product[fmt.Sprintf("productOptionName%d", i)]; ok && name != "" {
			dsc := product[fmt.Sprintf("productOptionDescription%d", i)]
			listOptions = fmt.Sprintf("%s%s|%s;", listOptions, name, dsc)
		} else {
			paramExist = false
		}
	}

	rs[yandexColNumber["Характеристики товара"]] = listOptions

	return rs
}

func filterSKU(s string) string {
	r := ""
	for _, char := range s {
		if strings.Contains(skuFilter, strings.ToLower(string(char))) {
			r = r + string(char)
		}
	}
	return r
}

func filterCharacters(s string, ASCIIFilter bool, ranges ...*unicode.RangeTable) string {
	var rs = make([]byte, 0, len([]byte(s)))
	for _, runeValue := range s {
		if ASCIIFilter && runeValue < unicode.MaxASCII {
			rs = append(rs, byte(runeValue))
		} else if unicode.In(runeValue, ranges...) {
			r := fmt.Sprintf("%c", runeValue)
			rs = append(rs, []byte(r)...)
		}
	}
	return string(rs)
}

func filterHtml(s string) string {

	p := bluemonday.NewPolicy()
	p.AllowElements("br")
	rs := p.Sanitize(s)

	rs = html.UnescapeString(rs)
	rs = strings.Replace(rs, "<br>", "\n", -1)

	return rs
}

func checkUrl(url string) (bool, error) {
	checkUrlLock <- struct{}{}
	defer func() { <-checkUrlLock }()

	cli := http.Client{Timeout: 5 * time.Second}

	resp, err := cli.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if strings.Contains(string(body), "Нет такого товара") {
		return false, nil
	}

	return true, nil
}
