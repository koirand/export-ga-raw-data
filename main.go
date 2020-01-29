package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/api/analytics/v3"
	"google.golang.org/api/option"
)

const credentialsPath = "./credentials.json"

func main() {
	s, err := getService()
	if err != nil {
		log.Fatalln(err)
	}

	accounts, err := s.Management.Accounts.List().Do()
	if err != nil {
		log.Fatalln(err)
	}

	for _, account := range accounts.Items {
		properties, err := s.Management.Webproperties.List(account.Id).Do()
		if err != nil {
			log.Fatalln(err)
		}
		for _, property := range properties.Items {
			profiles, err := s.Management.Profiles.List(account.Id, property.Id).Do()
			if err != nil {
				log.Fatalln(err)
			}
			for _, profile := range profiles.Items {
				results, err := s.Data.Ga.Get("ga:"+profile.Id, "7daysAgo", "today", "ga:pageviews").
					Output("json").
					Dimensions("ga:clientId,ga:pagePath,ga:dateHourMinute").
					Do()
				if err != nil {
					log.Fatalln(err)
				}
				for _, row := range results.Rows {
					fmt.Println(strings.Join(row, ","))
				}
			}
		}
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
