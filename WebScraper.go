package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// initialize a data structure to keep the scraped data
type Product struct {
	Url, Image, Name, Price string
}

func main() {

	// instantiate a new collector object
	coll := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
	)

	// initialize the slice of structs that will contain the scraped data
	var prods []Product

	// OnHTML callback
	coll.OnHTML("li.product", func(e *colly.HTMLElement) {

		// initialize a new Product instance
		product := Product{}

		// scrape the target data
		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Price = e.ChildText(".price")

		// add the product instance with scraped data to the list of products
		prods = append(prods, product)

	})

	// store the data to a CSV after extraction
	coll.OnScraped(func(r *colly.Response) {

		// open the CSV file
		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err)
		}
		defer file.Close()

		// initialize a file writer
		writer := csv.NewWriter(file)

		// write the CSV headers
		headers := []string{
			"Url",
			"Image",
			"Name",
			"Price",
		}
		writer.Write(headers)

		// write each product as a CSV row
		for _, product := range prods {
			// convert a Product to an array of strings
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}

			// add a CSV record to the output file
			writer.Write(record)
		}
		defer writer.Flush()
	})

	// open the target URL
	coll.Visit("https://www.scrapingcourse.com/ecommerce")

}
