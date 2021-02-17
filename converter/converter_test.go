package converter

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/leporel/wixtoyandex/config"
	"testing"
	"unicode"
)

func TestRangeTable(t *testing.T) {
	for _, r := range unicode.Cyrillic.R16 {
		fmt.Println("\nLo:", r.Lo, "Hi:", r.Hi, "Stride:", r.Stride)
		for c := r.Lo; c <= r.Hi; c += r.Stride {
			fmt.Print(string(rune(c)) + " ")
		}
	}
}

func TestConverter(t *testing.T) {
	var in = []byte(`handleId,fieldType,name,description,productImageUrl,collection,sku,ribbon,price,surcharge,visible,discountMode,discountValue,inventory,weight,productOptionName1,productOptionType1,productOptionDescription1,productOptionName2,productOptionType2,productOptionDescription2,productOptionName3,productOptionType3,productOptionDescription3,productOptionName4,productOptionType4,productOptionDescription4,productOptionName5,productOptionType5,productOptionDescription5,productOptionName6,productOptionType6,productOptionDescription6,additionalInfoTitle1,additionalInfoDescription1,additionalInfoTitle2,additionalInfoDescription2,additionalInfoTitle3,additionalInfoDescription3,additionalInfoTitle4,additionalInfoDescription4,additionalInfoTitle5,additionalInfoDescription5,additionalInfoTitle6,additionalInfoDescription6,customTextField1,customTextCharLimit1,customTextMandatory1,customTextField2,customTextCharLimit2,customTextMandatory2
product_aa6471bb-72d8-874e-f885-8f50d6f92812,Product,Салют Высший пилотаж,"<p>⭕️Калибр - 0.8""<br>⠀<br>💢Кол-во выстрелов - 20<br>⠀<br>⏰Время работы - 30&nbsp;сек<br>⠀<br>🔺Высота разрыва - 20&nbsp;м<br>⠀<br>💥Эффект:&nbsp;кокосовое дерево, ива, пион, мерцание, парчовая корона</p>",5e7d77_5b250734eaaf45fe9d6f0150b6609cb8~mv2.png;5e7d77_d3f93736bdfe4dae9c241dce2c3aadd9,Салюты,РК8051,,990.0,,true,PERCENT,0.0,30,0.0,Диаметр,DROP_DOWN,0.8 дюйма,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
`)
	var out bytes.Buffer

	config.Cfg.ConvertParams.CheckUrl = true

	params := config.Cfg.ConvertParams
	params.Url = "https://www.skyindiamonds.ru"
	processDefaultRow(params)

	reader := bytes.NewReader(in)
	writer := bufio.NewWriter(&out)

	err := convert(reader, writer)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n", out.String())
}

func TestCheckUrl(t *testing.T) {
	ok, err := checkUrl("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ok)
}
