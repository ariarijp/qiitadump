package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/koron/go-dproxy"
)

func getItems(token *string, page int, perPage int) ([]interface{}, error) {
	u, err := url.Parse("http://qiita.com/api/v2/authenticated_user/items")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("per_page", strconv.Itoa(perPage))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+*token)

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v interface{}
	json.Unmarshal(byteArray, &v)

	rawItems, err := dproxy.New(v).Array()
	if err != nil {
		return nil, err
	}

	return rawItems, nil
}

func main() {
	optToken := flag.String("token", "", "Access token")
	optWithoutPrivate := flag.String("without-private", "true", "Dump without private item")
	flag.Parse()

	if *optToken == "" {
		log.Fatal("Access token is required. Please set your valid access token.")
	}

	withoutPrivate := true
	if *optWithoutPrivate == "false" {
		withoutPrivate = false
	}

	perPage := 20
	for page := 1; ; page++ {
		rawItems, err := getItems(optToken, page, perPage)
		if err != nil {
			log.Fatal(err)
		}

		for _, rawItem := range rawItems {
			item := rawItem.(map[string]interface{})

			if withoutPrivate && item["private"].(bool) {
				continue
			}

			bytes, err := json.Marshal(item)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(bytes))
		}

		if len(rawItems) < perPage {
			break
		}
	}
}
