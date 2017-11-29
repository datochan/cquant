package cmd

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"math/rand"
	"github.com/urfave/cli"
	"github.com/kniren/gota/series"
	"github.com/kniren/gota/dataframe"

	"github.com/datochan/gcom/cnet"
	"github.com/datochan/gcom/utils"
	"github.com/datochan/gcom/logger"
	"github.com/datochan/ctdx"
	"cquant/comm"
)

func Calendar(c *cli.Context) error {
	logger.Info("开始更新股市交易日历...")
	configure := c.App.Metadata["configure"].(*comm.Configure)

	content := cnet.HttpRequest(configure.Datayes.Urls.Calendar, "", "", configure.Datayes.Token, "")
	df := dataframe.ReadCSV(strings.NewReader(utils.ConvertTo(content, "gbk", "utf8")))

	var curDateAry []string
	var prevDateAry []string

	for _, row := range df.Maps() {
		strDate := row["calendarDate"].(string)
		strPrevDate := row["prevTradeDate"].(string)

		curDateAry = append(curDateAry, strings.Replace(strDate, "-", "", -1))
		prevDateAry = append(prevDateAry, strings.Replace(strPrevDate, "-", "", -1))
	}

	df = df.Mutate(series.New(curDateAry, series.String, "calendarDate"))
	df = df.Mutate(series.New(prevDateAry, series.String, "prevTradeDate"))

	calendarPath := fmt.Sprintf("%s%s", configure.App.DataPath, configure.Tdx.Files.Calendar)
	utils.WriteCSV(calendarPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, &df)

	logger.Info("股市交易日历更新完毕...")
	return nil
}

func Basics(c *cli.Context) error {
	logger.Info("准备更新市场基础数据...")
	configure := c.App.Metadata["configure"].(*comm.Configure)

	tdxClient := ctdx.NewDefaultTdxClient(configure)
	defer tdxClient.Close()

	tdxClient.Conn()
	tdxClient.UpdateStockBase()

	// 更新结束后会向管道中发送一个通知
	<- tdxClient.Finished
	logger.Info("市场基础数据更新完毕...")
	tdxClient.Close()
	return nil
}

func Bonus(c *cli.Context) error {
	logger.Info("准备更新高送转数据...")
	configure := c.App.Metadata["configure"].(*comm.Configure)

	tdxClient := ctdx.NewDefaultTdxClient(configure)
	defer tdxClient.Close()

	tdxClient.Conn()
	tdxClient.UpdateStockBonus()

	// 更新结束后会向管道中发送一个通知
	<- tdxClient.Finished
	logger.Info("高送转数据更新完毕...")

	tdxClient.Close()

	return nil
}

func Days(c *cli.Context) error {
	logger.Info("准备更新日线数据...")
	configure := c.App.Metadata["configure"].(*comm.Configure)

	tdxClient := ctdx.NewDefaultTdxClient(configure)
	defer tdxClient.Close()

	tdxClient.Conn()
	tdxClient.UpdateDays()

	// 更新结束后会向管道中发送一个通知
	<- tdxClient.Finished
	logger.Info("日线数据更新完毕...")

	tdxClient.Close()

	return nil
}

func Mins(c *cli.Context) error {
	logger.Info("准备更新五分钟线数据...")
	configure := c.App.Metadata["configure"].(*comm.Configure)

	tdxClient := ctdx.NewDefaultTdxClient(configure)
	defer tdxClient.Close()

	tdxClient.Conn()
	tdxClient.UpdateMins()

	// 更新结束后会向管道中发送一个通知
	<- tdxClient.Finished
	logger.Info("五分钟线数据更新完毕...")

	tdxClient.Close()

	return nil
}


func updateST(configure *comm.Configure) error {
	isAppend := false

	stockSTPath := fmt.Sprintf("%s%s", configure.App.DataPath, configure.Datayes.Files.StockSt)
	colTypes := map[string]series.Type{ "date": series.Int, "code": series.String, "name": series.String,
		"flag": series.String}

	stockItemDF := utils.ReadCSV(stockSTPath, dataframe.WithTypes(colTypes))

	start := "19980101"
	if nil == stockItemDF.Err {
		// 获取最后一条记录的日期
		isAppend = true

		idx := utils.FindInStringSlice("date", stockItemDF.Names())
		start = stockItemDF.Elem(stockItemDF.Nrow()-1, idx).String()
		start = utils.AddDays(start, 1)
	}

	nowDate, _ := strconv.Atoi(utils.Today())
	endDate, _ := strconv.Atoi(start)

	if endDate > nowDate {
		errMsg := fmt.Errorf("ST信息已是最新,无需继续更新")
		logger.Error("%v", errMsg)
		return errMsg
	}

	strEnd := utils.AddDays(start, 6*30)
	endDate, _ = strconv.Atoi(strEnd)
	if endDate > nowDate {
		strEnd = utils.Today()
	}

	logger.Info("\t正在获取 %s 至 %s 之间的ST信息", start, strEnd)

	stUrl := fmt.Sprintf(configure.Datayes.Urls.StockSt, start, strEnd)
	content := cnet.HttpRequest(stUrl, "", "", configure.Datayes.Token, "")
	if len(content) <= 20 {
		errMsg := fmt.Errorf("获取ST数据失败, 获取内容为:%s", content)
		logger.Error("%v", errMsg)
		return errMsg
	}

	df := dataframe.ReadCSV(strings.NewReader(utils.ConvertTo(content, "gbk", "utf8")),
		dataframe.WithTypes(map[string]series.Type{ "ticker": series.String}))
	df.SetNames("date", "code", "name", "flag")

	var curDateAry []int

	for _, row := range df.Maps() {
		strDate := row["date"].(string)
		nDate, _ := strconv.Atoi(strings.Replace(strDate, "-", "", -1))
		curDateAry = append(curDateAry, nDate)
	}

	df = df.Mutate(series.New(curDateAry, series.Int, "date"))

	sortedDf := df.Arrange(dataframe.Sort("date"))
	if !isAppend {
		utils.WriteCSV(stockSTPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, &sortedDf)
	} else {
		utils.WriteCSV(stockSTPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, &sortedDf, dataframe.WriteHeader(false))
	}

	return nil
}

func ST(c *cli.Context) error {
	logger.Info("开始更新股市ST股票记录...")
	configure := c.App.Metadata["configure"].(*comm.Configure)
	for nil == updateST(configure) {
		idx := time.Duration(rand.Intn(4)) + 1
		time.Sleep(time.Second*idx)
	}

	logger.Info("股市ST股票记录更新结束...")
	return nil
}

func Report(c *cli.Context) error {
	fmt.Println("Report: ", c.Args().First())
	return nil
}

