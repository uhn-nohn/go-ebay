package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/xml"
	"strconv"
)

type Item struct {
	ItemID string `xml:"itemId"`
	Title string `xml:"title"`
	Location string `xml:"location"`
	CurrentPrice float64 `xml:"sellingStatus>currentPrice"`
	ShippingPrice float64 `xml:"shippingInfo>shippingServiceCost"`
	BINprice float64 `xml:"listingInfo>buyItNowPrice"`
}

type ResponseXML struct {
	XMLName xml.Name `xml:"findItemsByKeywordsResponse"`
	Items []Item `xml:"searchResult>item"`
}


// build the url for HTTP GET 
func buildRequest(appid, keywords, n string) string {
	keywords = keywordConvert(keywords)

	url := "http://svcs.ebay.com/services/search/FindingService/v1?OPERATION-NAME=findItemsByKeywords&SERVICE-VERSION=1.0.0&SECURITY-APPNAME="
	url += appid
	url += "&GLOBAL-ID=EBAY-US&RESPONSE-DATA-FORMAT=XML&REST-PAYLOAD&keywords="
	url += keywords
	url += "&paginationInput.entriesPerPage="
	url += n
	//Comment out the next line if you do not want to limit results to Fixed Price and Auction with Buy It Now listings. 
	url += "&itemFilter(0).name=ListingType&itemFilter(0).value(0)=FixedPrice&itemFilter(0).value(1)=AuctionWithBIN"

	return url
}

// URL-encode any spaces in keywords
// all " " (spaces) should be converted to "%20"
func keywordConvert(keywords string) string {
	for i := 0; i < len(keywords); i++ {
    	if keywords[i] == ' ' {
    		keywords = keywords[:(i)] + "%20" + keywords[(i+1):]
    	}
    }
    return keywords
}

// print out information on each listing
func printListings(v ResponseXML) {
	for i := 0; i < len(v.Items); i++ {
    	fmt.Println(strconv.Itoa((i+1)) + ". " + v.Items[i].Title)
    	if(v.Items[i].BINprice != 0){
    		fmt.Println("Buy It Now price: $" + strconv.FormatFloat(v.Items[i].BINprice, 'f', 2, 64))
    	}else{
    		fmt.Println("Current price: $" + strconv.FormatFloat(v.Items[i].CurrentPrice, 'f', 2, 64))
    	}
    	if(v.Items[i].ShippingPrice != 0){
    		fmt.Println("Shipping: $" + strconv.FormatFloat(v.Items[i].ShippingPrice, 'f', 2, 64))
    	}else{
    		fmt.Println("Shipping: Free")
    	}
    	fmt.Println("Location: " + v.Items[i].Location)
    	fmt.Println("-------------------------------------")
    }
}

func sendAndProcessRequest(url string) ResponseXML {
	r, _ := http.Get(url)
	response, _ := ioutil.ReadAll(r.Body)
	v := ResponseXML{}
	err := xml.Unmarshal([]byte(response), &v)
	if err != nil {
   		fmt.Printf("error: %v", err)
    }
    return v
}


func main() {
	appID := "" // YOUR APPID HERE
	query := "" // SEARCH QUERY HERE
	num_items := "10"

	url := buildRequest(appID, query, num_items)

	v := sendAndProcessRequest(url)

    printListings(v)    
}