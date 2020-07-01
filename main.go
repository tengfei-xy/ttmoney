package main

import (
		"net/url"
		"net/http"
		"fmt"
		"io/ioutil"
		"github.com/PuerkitoBio/goquery"
		"encoding/json"
		"strings"
		"os"
		"time"
		"math/rand"
)
// Json Result 
type JResult struct {
	ErrCode		int		 `json:"ErrCode"`
	ErrMsg		string 	 `json:"ErrMsg"`
	Datas		[]JRData `json:"Datas"`
}
// Json Result
type JRData struct{
	ID				string	`json:"_id"`
	CODE			string	`json:"CODE"`
	Name			string	`json:"NAME"`
	SHORTNAME		string	`json:"SHORTNAME"`
}

// Json Excel Line Head
type JSData struct {
	ID				string // 基金ID
	FName			string // 基金名称
	FCreate			string // 基金创建时间
	FLink			string // 基金链接
	ErrGZ			string // 基金跟踪误差率
	ErrPJ			string // 基金平均误差率
	ErrLink			string // 基金误差网页
	Rate			string // 基金费率
	Way				string // 基金交易方式

	CName			string // 公司名称
	CLink			string // 公司链接
	CTime			string // 公司创建时间
	CScale			string // 公司规模

	Include			bool
}
// Json Hude Company
type JCompany	struct{
	CName			string
	CLink			string
	CTime			string
	CScale			string
}

// 基本结构体

type tt_search struct {
	Search_link			string
	Search_key			string
	Search_result		JResult
	Invalue				[4]string
}
type tt_Data struct {
	Count				int
	Data				[]JSData				
}
const d int =10
//公司、规模、成立时间、基金成立时间、跟踪误差率、同类平均误差率、费率、交易方式


func main(){
	var tt tt_search
	client := &http.Client{}
	var key string	
	fmt.Println("输入基金名称:")
	fmt.Scanln(&key)

	// 获取千亿公司
	c := GetHugeCompany(client)

	tt.Search_link 		= "http://fundsuggest.eastmoney.com/FundSearch/api/FundSearchPageAPI.ashx?m=1&pageindex=0&pagesize=200&key="
	tt.Search_key 		= url.QueryEscape(key)
	tt.Invalue  		= [4]string{"分级","等权","增强","优选"}
	tt.Search_link 		+= tt.Search_key

	tt.GetSearchResult(client)
	fd := tt.SecectData()
	fd.InitFond(client,&c)
	fd.InitError(client)
	fd.Output()

	fmt.Print("\n筛选基金完成")
	fmt.Scanln(&key)

}
func GetHugeCompany(client  * http.Client) []JCompany {
	link := `http://fund.eastmoney.com/company/default.html`
	fmt.Println("千亿公司链接 ",link)
	r, err := http.NewRequest("GET", link, nil)
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;")
	r.Header.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Connection", "keep-alive")
	r.Header.Add("Host", "fund.eastmoney.com")
	r.Header.Add("Pragma", "no-cache")
	r.Header.Add("Referer", "http://fund.eastmoney.com/110003.html")
	r.Header.Add("Upgrade-Insecure-Requests", "1" )
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
	if err != nil{
		panic(err)
	}
	//if res,err:= os.Open("hugecompany"); err !=nil{
	if res, err := client.Do(r);err !=nil{
		panic(err)
	}else{
		if doc,err := goquery.NewDocumentFromReader(res.Body);err !=nil{
		//if byteBody,err := ioutil.ReadAll(res.Body);err !=nil{
			panic(err)
		}else{
			defer res.Body.Close()		

			var count int
			//fmt.Println(string(byteBody))
			doc.Find("body").Find("p.td-gm")
			doc.Find("body").Find("p.td-gm").Each(func(i int ,s * goquery.Selection){
				if strings.Index(s.Text(),",") != -1{
					count += 1
					return
				}
			})
			base:= doc.Find("body").Find("table#gspmTbl>tbody>tr")
			var company []JCompany = make ([]JCompany,count)
			for i :=0; i < count;i++{

				// 公司规模（千亿）
				money := doc.Find("body").Find("p.td-gm").Eq(i).Text()
				company[i].CScale		= strings.Replace(money[:strings.Index(money,".")+3],",","",1)

				// 公司名称
				company[i].CName	= base.Find("td.td-align-left>a").Eq(i).Text()

				// 公司成立时间
				company[i].CTime	= base.Find("td.menu-link").Next().Eq(i).Text()

				// 公司代码
				code,_ := base.Find("td.td-align-left>a").Eq(i).Attr("href")
				// 输出结果：80041198
				company[i].CLink	= code[strings.LastIndex(code,"/")+1:strings.Index(code,".")]
				// 输出结果：/Company/80041198.html
				//company[i].CLink	= code

				// 输出
				fmt.Printf("公司名称:%s, 公司代码:%s, 公司成立时间:%s, 公司规模(亿元):%s\n",company[i].CName,company[i].CLink,company[i].CTime,company[i].CScale)
			}
			Delay(d)
			return company
		}
	}
}
func (tt * tt_search)GetSearchResult(client  * http.Client){
	fmt.Println("\n关键词链接 ",tt.Search_link)
	r, err := http.NewRequest("GET", tt.Search_link, nil)
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;")
	r.Header.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Connection", "keep-alive")
	r.Header.Add("Pragma", "no-cache")
	r.Header.Add("Host", "fundsuggest.eastmoney.com")
	r.Header.Add("Referer", "http://fund.eastmoney.com/data/fundsearch.html?spm=search&key=" + tt.Search_key )
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")

	if err != nil{
		panic(err)
	}

	// 开始模拟搜索关键词
	res, err := client.Do(r)
	//res,err:= os.Open("result")
	if err !=nil{
		panic(err)
	}

	// 开始读取数据
	byteBody,err := ioutil.ReadAll(res.Body)
	if err !=nil{
		fmt.Println("读取JSON错误")
		panic(err)
	}

	// 输出并解析html
	if err := json.Unmarshal(byteBody,&tt.Search_result) ;err != nil{
		fmt.Println("解析JSON错误")
		panic(err)
	}
	Delay(d)
	defer res.Body.Close()
	
}

func (tt * tt_search)SecectData() tt_Data{
	var fd tt_Data
	var data []JSData	= make([]JSData,len( tt.Search_result.Datas))
	var Count			= 0
	var t 				= 0
	for j,i := range tt.Search_result.Datas{
		for _,y := range tt.Invalue{
			if strings.Index(i.Name,y)!=-1 {
				fmt.Printf("%d 剔除%s: %s\n",j+1,y,i.Name)

				t = t|1
				break
			}
		}
		if t==0{

			// 基金ID
			data[Count].ID = i.ID

			// 基金名称
			data[Count].FName = i.Name
			fmt.Printf("%d 保留: %s\n",j+1,i.Name)
			// 基金链接
			data[Count].FLink = `http://fund.eastmoney.com/`+i.ID+".html"
			Count +=1
		}
		t=0

	}
	fd.Count = Count
	fd.Data = data
	fmt.Printf("发现基金数:%d,保留基金数量:%d,过滤关键词基金数量:%d\n",len(fd.Data),fd.Count,len(fd.Data)-fd.Count)
	//fmt.Printf("预计需要%d秒\n",Count*7*3)
	return fd
}
// 基金成立时间、交易方式、费率
func (fd * tt_Data) InitFond(client  * http.Client,c * []JCompany){
	fmt.Println("\n开始抓取基金信息")
	for i:=0 ;i<fd.Count; i++ {

		var text string	
		if fd.Data[i].FLink == ""{
			fmt.Print("跳过")
			continue
		}
		r, err := http.NewRequest("GET", fd.Data[i].FLink, nil)
		r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;")
		r.Header.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
		r.Header.Add("Cache-Control", "no-cache")
		r.Header.Add("Proxy-Connection", "keep-alive")
		r.Header.Add("Host", `fund.eastmoney.com`)
		r.Header.Add("Referer", fd.Data[i].FLink)
		r.Header.Add("Upgrade-Insecure-Requests", "1")
		r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
		r.Header.Add("Pragma", "no-cache")
		if err != nil{
			panic(err)
		}
		res, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		//byteBody,err := ioutil.ReadAll(res.Body)
		//fmt.Print(string(byteBody))
		//res,_ := os.Open("index")
		//doc,err := goquery.NewDocumentFromReader(res)
		doc,err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			panic(err)
		}

		defer res.Body.Close()

		// 基金的公司代码
		base := doc.Find("body")
		text,_ = base.Find("div.fundDetail-footer>ul>li>a").Eq(3).Attr("href")
		for _,j := range  *c{
			if strings.Index(text,j.CLink)!=-1{
				//fmt.Printf("test:%s   j.Clink:%s\n",text,j.CLink)

				// 公司名称
				fd.Data[i].CName			= j.CName

				// 公司链接
				fd.Data[i].CLink 			= text

				// 公司成立时间
				fd.Data[i].CTime			= j.CTime

				// 公司规模
				fd.Data[i].CScale			= j.CScale

				// 基金成立时间
				text						= base.Find("div.infoOfFund>table>tbody>tr").Next().Find("td").Eq(0).Text()
				fd.Data[i].FCreate			= strings.Split(text,"：")[1]

				// 基金交易状态/方式
				text						= base.Find("span.staticCell").Eq(0).Text()
				if text == "" || text == "封闭期" || (strings.Index(text,"暂停") != -1) { fmt.Printf("过滤 %s 原因:%s\n 链接:%s\n",fd.Data[i].FName,text,fd.Data[i].FLink); continue }
				text						= base.Find("span.staticCell").Text()
				if (strings.Index(text,"不可购买") != -1) || (strings.Index(text,"不开放") != -1) {fmt.Printf("过滤 %s 原因:%s 链接%s\n",fd.Data[i].FName,text,fd.Data[i].FLink); continue }

				fd.Data[i].Way				= strings.TrimSuffix(text," ")
 
				// 基金费率
				text 						= base.Find("span.nowPrice").Text()
				fd.Data[i].Rate				= text

				// 基金误差链接
				text,_ 						= base.Find("td.specialData>a").Eq(1).Attr("href")
				fd.Data[i].ErrLink			= text

				// 包括
				fd.Data[i].Include			= true

				// 输出信息
				fmt.Printf("名称:%s基,基金链接:%s,基金公司:%s\n",fd.Data[i].FName,fd.Data[i].FLink,fd.Data[i].CName)
				
				Delay(d)
				break 

			}
		}
	}
}

// 跟踪误差率,同类误差率
func (fd * tt_Data) InitError(client  * http.Client){
	fmt.Println()
	for i:=0 ;i<fd.Count; i++{
		if !fd.Data[i].Include{
			continue
		}
		if fd.Data[i].ErrLink == "" {
			fd.Data[i].ErrGZ = "--"
			fd.Data[i].ErrPJ = "--"
			continue
		}

		fmt.Printf("基金名称%s,请求误差链接:%s\n",fd.Data[i].FName,fd.Data[i].ErrLink)
		r, err := http.NewRequest("GET", fd.Data[i].ErrLink, nil)
		r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;")
		r.Header.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
		r.Header.Add("Cache-Control", "no-cache")
		r.Header.Add("Proxy-Connection", "keep-alive")
		r.Header.Add("Host", "fundf10.eastmoney.com")
		r.Header.Add("Referer", fd.Data[i].ErrLink)
		r.Header.Add("Upgrade-Insecure-Requests", "1" )
		r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36")
		r.Header.Add("Pragma", "no-cache")
		if err != nil{
			panic(err)
		}
		res, err := client.Do(r)
		if err != nil {
			panic(err)
		}
		//byteBody,err := ioutil.ReadAll(res.Body)
		//fmt.Print(string(byteBody))
		if err != nil {
			panic(err)
		}
		// res,_ := os.Open("error")
		//doc,err := goquery.NewDocumtFromReader(res)
		doc,err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			panic(err)
		}
		base := doc.Find("body").Find("div#jjzsfj").Eq(0)
		tbase := base.Find("table.fxtb>tbody>tr").Eq(1).Find("td")
		fd.Data[i].ErrGZ = tbase.Eq(1).Text()
		fd.Data[i].ErrPJ = tbase.Eq(2).Text()
		defer res.Body.Close()
		Delay(d)
	}
}

func (fd * tt_Data) Output() {
	filename := "基金筛选结果"
	fileext := ".csv"
	filename += fileext
	fmt.Printf("\n输出%s\n",filename)

	file,err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE  | os.O_TRUNC,0755)
	fb := true
	if err !=nil {
		fmt.Print("打开文件错误",err)
		fb=false
	}
	defer file.Close()

	for i:=0 ;i<fd.Count; i++{
		if !fd.Data[i].Include{
			continue
		}
		var content string
		//fmt.Printf("基金名称:%s,基金创建时间:%s,基金链接:%s,基金跟踪误差率:%s,同类基金平均误差率:%s,基金费率:%s,基金交易方式:%s,公司名称:%s,公司链接:%s,公司创建时间:%s,公司规模:%s\n",
		//fd.Data[i].ID
		if i ==0{
			content = fmt.Sprintf("基金名称,基金链接,基金创建时间,基金跟踪误差率,同类基金平均误差率,基金费率,基金交易方式,公司名称,公司链接,公司创建时间,公司规模\n")
		}else {
			content = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			fd.Data[i].FName,
			fd.Data[i].FLink,
			fd.Data[i].FCreate,
			fd.Data[i].ErrGZ,
			fd.Data[i].ErrPJ,
			fd.Data[i].Rate,
			fd.Data[i].Way,
			fd.Data[i].CName,
			fd.Data[i].CLink,
			fd.Data[i].CTime,
			fd.Data[i].CScale)
		}
		if fb {
			if _,err := file.WriteString(content); err != nil{
				fmt.Print("写入文件数据错误\n",err)
			}
		}
		fmt.Print(content)
	}
}
func Delay(i int){
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(i)) * time.Second)
}