package main

import (
	"log"
	"os"
	"strings"

	"encoding/csv"

	"golang.org/x/net/context"
	"google.golang.org/api/analytics/v3"
	"google.golang.org/api/option"
)

const credentialsPath = "./credentials.json"
const outputPath = "./ga.csv"
const viewID = "208412389"
const startDate = "90daysAgo"
const endDate = "yesterday"
const itemsPerPage int64 = 1000

func main() {
	// アウトプットファイル作成
	out, err := newCsvFile(outputPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()

	// サービス取得
	s, err := getService()
	if err != nil {
		log.Fatalln(err)
	}

	var page int64 = 0
	for {
		// rawデータ取得
		results, err := s.Data.Ga.Get("ga:"+viewID, startDate, endDate, "ga:pageviews").
			Output("json").
			Dimensions(strings.Join([]string{
				"ga:clientId",
				"ga:pagePath",
				"ga:pageTitle",
				"ga:dateHourMinute",
			}, ",")).
			StartIndex(1 + itemsPerPage*page).
			Do()
		if err != nil {
			log.Fatalln(err)
		}

		// 出力
		for _, row := range results.Rows {
			out.Write(row)
		}
		out.Flush()

		if results.NextLink == "" {
			break
		}
		page++
	}
}

func getService() (*analytics.Service, error) {
	ctx := context.Background()
	as, err := analytics.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}
	return as, nil
}

type CsvFile struct {
	File   *os.File
	Writer *csv.Writer
}

func (c *CsvFile) Close() {
	c.File.Close()
}

func (c *CsvFile) Write(record []string) {
	c.Writer.Write(record)
}

func (c *CsvFile) Flush() {
	c.Writer.Flush()
}
func newCsvFile(fileName string) (*CsvFile, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return &CsvFile{
		File:   file,
		Writer: csv.NewWriter(file),
	}, nil
}
