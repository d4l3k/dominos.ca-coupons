package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const url = "https://order.dominos.ca/power/store/10057/coupon/%d?lang=en"

type Tags struct {
	Hidden bool
}

type Coupon struct {
	ID                                                                                       int
	Status                                                                                   int
	StoreID, BusinessDate, StoreAsOfTime, LanguageCode, MasterSortSeq, PulseSortSeq          string
	Local, Bundle                                                                            bool
	Code                                                                                     string
	Tags                                                                                     Tags
	Name, Description, Price, ImageCode, SizeLargeImageURL, SizeThumbNailImageURL, PulseCode string
}

func scrape(id int) (*Coupon, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(url, id), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.75 Safari/537.36")
	req.Header.Set("Referer", "https://order.dominos.ca/en/pages/order/?redirect=/section/Coupons/category/All/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Cookie", "returnUser=true; noCheeseUpSell=true")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var coupon Coupon
	coupon.ID = id
	if err := json.NewDecoder(resp.Body).Decode(&coupon); err != nil {
		return nil, err
	}
	return &coupon, nil
}

func main() {
	c := make(chan int)
	r := make(chan *Coupon)
	done := false
	go func() {
		for i := 1000; i < 10000; i++ {
			c <- i
		}
		close(c)
		done = true
	}()
	for i := 0; i < 100; i++ {
		go func() {
			for id := range c {
				coupon, err := scrape(id)
				if err != nil {
					log.Fatal(err)
				}
				r <- coupon
			}
			if done {
				close(r)
			}
		}()
	}
	for coupon := range r {
		if coupon.Status == 0 {
			fmt.Printf("%d ($%s): %s\n", coupon.ID, coupon.Price, coupon.Name)
		}
	}
}
