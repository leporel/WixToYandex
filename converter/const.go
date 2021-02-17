package converter

var wixProductFields = make(map[int]string)

const skuFilter = "abcdefghijklmnopqrstuvwxyz0123456789"

/*
Wix fields 06.02.2021

handleId
fieldType
name
description
productImageUrl
collection
sku
ribbon
price
surcharge
visible
discountMode
discountValue
inventory
weight
productOptionName1
productOptionType1
productOptionDescription1
productOptionName2
productOptionType2
productOptionDescription2
productOptionName3
productOptionType3
productOptionDescription3
productOptionName4
productOptionType4
productOptionDescription4
productOptionName5
productOptionType5
productOptionDescription5
productOptionName6
productOptionType6
productOptionDescription6
additionalInfoTitle1
additionalInfoDescription1
additionalInfoTitle2
additionalInfoDescription2
additionalInfoTitle3
additionalInfoDescription3
additionalInfoTitle4
additionalInfoDescription4
additionalInfoTitle5
additionalInfoDescription5
additionalInfoTitle6
additionalInfoDescription6
customTextField1
customTextCharLimit1
customTextMandatory1
customTextField2
customTextCharLimit2
customTextMandatory2
*/

var yandexColNumber = map[string]int{
	"id":                           0,
	"Статус товара":                1,
	"Доставка":                     2,
	"Стоимость доставки":           3,
	"Срок доставки":                4,
	"Самовывоз":                    5,
	"Стоимость самовывоза":         6,
	"Срок самовывоза":              7,
	"Купить в магазине без заказа": 8,
	"Ссылка на товар на сайте магазина": 9,
	"Производитель":                     10,
	"Название":                          11,
	"Категория":                         12,
	"Цена":                              13,
	"Цена без скидки":                   14,
	"Валюта":                            15,
	"Ссылка на картинку":                16,
	"Описание":                          17,
	"Характеристики товара":             18,
	"Условия продажи":                   19,
	"Гарантия производителя":            20,
	"Страна происхождения":              21,
	"Штрихкод":                          22,
	"bid":                               23,
	"Уцененный товар":                   24,
	"Причина уценки":                    25,
	"Кредитная программа":               26,
}

var defaultYandexRow = []string{
	"",          //	id*
	"В наличии", //	Статус товара
	"Нет",       //	Доставка
	"",          //	Стоимость доставки
	"",          //	Срок доставки
	"Есть",      //	Самовывоз
	"",          //	Стоимость самовывоза
	"",          //	Срок самовывоза
	"Можно",     //	Купить в магазине без заказа
	"",          //	Ссылка на товар на сайте магазина*
	"",          //	Производитель
	"",          //	Название*
	"",          //	Категория*
	"",          //	Цена*
	"",          //	Цена без скидки
	"",          //	Валюта*
	"",          //	Ссылка на картинку*
	"",          //	Описание
	"",          //	Характеристики товара
	"",          //	Условия продажи
	"Нет",       //	Гарантия производителя
	"",          //	Страна происхождения
	"",          //	Штрихкод
	"",          //	bid
	"",          //	Уцененный товар
	"",          //	Причина уценки
	"",          //	Кредитная программа
}
