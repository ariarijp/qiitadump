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

const perPage int = 20

func getItems(host *string, endpoint *string, token *string, page int) ([]interface{}, error) {
	endpointURL := fmt.Sprintf("https://%s%s", *host, *endpoint)
	u, err := url.Parse(endpointURL)
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
	optHost := flag.String("host", "qiita.com", "Host")
	optEndpoint := flag.String("endpoint", "/api/v2/authenticated_user/items", "Endpoint")
	optWithoutPrivate := flag.String("without-private", "true", "Dump without private items")
	optLimit := flag.Int("limit", 20, "Limit the number of items")
	flag.Parse()

	if *optToken == "" {
		log.Fatal("Access token is required. Please set your valid access token.")
	}

	withoutPrivate := true
	if *optWithoutPrivate == "false" {
		withoutPrivate = false
	}

	num := 1
	for page := 1; ; page++ {
		rawItems, err := getItems(optHost, optEndpoint, optToken, page)
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

			if num >= *optLimit {
				break
			}

			num++
		}

		if num >= *optLimit || len(rawItems) < perPage {
			break
		}
	}
}
