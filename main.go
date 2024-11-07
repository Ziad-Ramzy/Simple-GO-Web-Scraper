package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"

    "github.com/gocolly/colly"
)

func main() {
    type Book struct {
        Url, Image, Name, Price string
    }

    var books []Book

    c := colly.NewCollector()

    c.OnHTML("li.product", func(e *colly.HTMLElement) {
        book := Book{}
        book.Url = e.ChildAttr("a", "href")
        book.Image = e.ChildAttr("img", "src")
        book.Name = e.ChildText(".title")
        book.Price = e.ChildText(".woocommerce-Price-amount")

        books = append(books, book)
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting: ", r.URL)
    })

    c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong: ", err)
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Page visited: ", r.Request.URL)
    })

    c.OnHTML("a", func(e *colly.HTMLElement) {
        fmt.Printf("Found link: %v\n", e.Attr("href"))
    })

    c.OnScraped(func(r *colly.Response) {
        fmt.Println(r.Request.URL, "scraped!")

        file, err := os.Create("books.csv")
        if err != nil {
            log.Fatalln("Failed to create output CSV file", err)
        }
        defer file.Close()

        writer := csv.NewWriter(file)
        defer writer.Flush()

        headers := []string{"Url", "Image", "Name", "Price"}
        writer.Write(headers)

        for _, book := range books {
            record := []string{book.Url, book.Image, book.Name, book.Price}
            writer.Write(record)
        }
    })

    baseURL := "https://diwanegypt.com/product-category/books/english-adults/"
    maxPages := 167
    
    for i := 1; i <= maxPages; i++ {
    var url string
    if i == 1 {
        url = baseURL
    } else {
        url = fmt.Sprintf("%s/page/%d/", baseURL, i)
    }
    c.Visit(url)
}

}
