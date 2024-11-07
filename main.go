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
        book.Name = e.ChildText("h2")
        book.Price = e.ChildText(".woocommerce-Price-amount.amount")

        books = append(books, book)
    })

    c.OnHTML("a.next.page-numbers", func(e *colly.HTMLElement) {
    nextPage := e.Attr("href")
    fmt.Println("Visiting:", nextPage)
    e.Request.Visit(nextPage) // Recursive visit to the next page
    })

    startURL := "https://diwanegypt.com/product-category/books/english-adults/"
    c.Visit(startURL)


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
    writer.Flush() 
}
