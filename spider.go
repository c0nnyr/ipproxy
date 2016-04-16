package main

import (
	"github.com/PuerkitoBio/goquery"
	"ipproxy/common"
	"log"
	"strconv"
	"strings"
)

// spider
func crawl_xicidaili(ch chan<- int, ind int) {
	// www.xicidaili.com
	url_base := "http://www.xicidaili.com/wn/"
	cur_url := url_base + strconv.Itoa(ind)
	doc, err := goquery.NewDocument(cur_url)
	log.Println("==============READING ", cur_url)
	if err != nil {
		log.Fatal(err)
		ch <- 1
		return
	}
	if doc.Find("#ip_list > tbody > tr").Length() <= 1 {
		ch <- 0
		return
	}
	doc.Find("#ip_list > tbody > tr").Each(handling_xicidaili_tr)
	ch <- 1
}

func handling_xicidaili_tr(i int, s *goquery.Selection) {
	if i != 0 {
		// Map将每一项处理后返回，Each仍旧是返回selection，方便链式
		items := s.Find("td").Map(func(i2 int, s2 *goquery.Selection) string {
			return s2.Text()
		})
		//item := new(common.IPProxyItem)//仍旧是非指针，只是没有顺便初始化
		item := common.IPProxyItem{}
		item.IP = strings.TrimSpace(items[2])
		item.Country = strings.TrimSpace(items[4])
		item.Hide_Type = strings.TrimSpace(items[5])
		item.Connection_Type = strings.TrimSpace(items[6])
		//log.Printf("%s, %s, %s, %s\n", item.IP, item.Country, item.Hide_Type, item.Connection_Type)
		common.UpsertItemToDB(item)
	}
}

func main() {
	common.ParrallelRun(crawl_xicidaili, 40)
}
