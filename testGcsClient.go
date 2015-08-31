package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Book struct {
	Name       string    `json:"name"`
	Author     string    `json:"author"`
	Pages      int       `json:"pages"`
	Year       int       `json:"year"`
	CreateTime time.Time `json:"createtime"`
}

// Target server
const BookURL = "https://testgcsserver.appspot.com/api/0.1/"
// const BookURL = "http://127.0.0.1:8080/api/0.1/"

const BookMaxPages = 1000

var BookName = []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF", "GGG", "HHH", "III", "JJJ"}
var BookAuthor = []string{"AuthorA", "AuthorB", "AuthorC", "AuthorD", "AuthorE", "AuthorF", "AuthorG", "AuthorH", "AuthorI", "AuthorJ"}

// Pring a Book
func (b Book) String() string {
	s := ""
	s += fmt.Sprintln("Name:", b.Name)
	s += fmt.Sprintln("Author:", b.Author)
	s += fmt.Sprintln("Pages:", b.Pages)
	s += fmt.Sprintln("Year:", b.Year)
	s += fmt.Sprintln("CreateTime:", b.CreateTime)
	return s
}

func queryAll() {
	// Send request
	resp, err := http.Get(BookURL + "books")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print status
	fmt.Println(resp.Status, resp.StatusCode)

	// Get body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode body
	var books map[string]Book = make(map[string]Book)
	if resp.StatusCode == http.StatusOK {
		// Decode as JSON
		if err := json.Unmarshal(body, &books); err != nil {
			fmt.Println(err, "in decoding JSON")
			return
		}
		for i, v := range books {
			fmt.Println("-------------------------------")
			fmt.Println("Key:", i)
			fmt.Println(v)
		}
		fmt.Println("Total", len(books), "books")
	} else {
		// Decode as text
		fmt.Printf("%s", body)
	}
}

func queryBook() {
	// Make URL
	var u *url.URL
	var err error
	if u, err = url.ParseRequestURI(BookURL + "books"); err != nil {
		fmt.Println(err, "in making URL")
		return
	}
	var q url.Values = u.Query()
	q.Add("Name", BookName[rand.Intn(len(BookName))])
	u.RawQuery = q.Encode()

	// Send request
	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print status
	fmt.Println(resp.Status, resp.StatusCode)

	// Get body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode body
	var books map[string]Book = make(map[string]Book)
	if resp.StatusCode == http.StatusOK {
		// Decode as JSON
		if err := json.Unmarshal(body, &books); err != nil {
			fmt.Println(err, "in decoding JSON")
			return
		}
		for i, v := range books {
			fmt.Println("-------------------------------")
			fmt.Println("Key:", i)
			fmt.Println(v)
		}
		fmt.Println("Total", len(books), "books")
	} else {
		// Decode as text
		fmt.Printf("%s", body)
	}
}

// Return 0: success
// Return 1: failed
func storeTen() int {
	// Send request
	resp, err := http.Post(BookURL+"storeTen", "", nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status, resp.StatusCode)
	if resp.StatusCode == http.StatusCreated {
		return 0
	} else {
		return 1
	}
}

// Return
// int = 0: success
//       1: failed
// string is the new book's unique key
func storeBook() (r int, key string) {
	// Return value
	r = 0
	key = ""

	// Make body
	book := Book{
		Name:       BookName[rand.Intn(len(BookName))],
		Author:     BookAuthor[rand.Intn(len(BookAuthor))],
		Pages:      rand.Intn(BookMaxPages),
		Year:       rand.Intn(time.Now().Year()),
		CreateTime: time.Now(),
	}
	b, err := json.Marshal(book)
	if err != nil {
		fmt.Println(err, "in encoding a book as JSON")
		r = 1
		return
	}

	// Send request
	resp, err := http.Post(BookURL+"books", "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)
		r = 1
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status, resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		r = 1
		return
	}
	url, err := resp.Location()
	if err != nil {
		fmt.Println(err, "in getting location from response")
		return
	}
	fmt.Println("Location is", url)

	// Get key from URL
	tokens := strings.Split(url.Path, "/")
	var keyIndexInTokens int = 0
	for i, v := range tokens {
		if v == "books" {
			keyIndexInTokens = i + 1
		}
	}
	if keyIndexInTokens >= len(tokens) {
		fmt.Println("Key is not given")
		return
	}
	key = tokens[keyIndexInTokens]
	if key == "" {
		fmt.Println("Key is empty")
		return
	}
	return
}

// Return 0: success
// Return 1: failed
func deleteBook(key string) int {
	// Send request
	pReq, err := http.NewRequest("DELETE", BookURL+"books/"+key, nil)
	if err != nil {
		fmt.Println(err, "in making request")
		return 1
	}
	resp, err := http.DefaultClient.Do(pReq)
	if err != nil {
		fmt.Println(err, "in sending request")
		return 1
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status, resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		return 0
	} else {
		return 1
	}
}

// Return 0: success
// Return 1: failed
func deleteAll() int {
	// Send request
	pReq, err := http.NewRequest("DELETE", BookURL+"books", nil)
	if err != nil {
		fmt.Println(err, "in making request")
		return 1
	}
	resp, err := http.DefaultClient.Do(pReq)
	if err != nil {
		fmt.Println(err, "in sending request")
		return 1
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status, resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		return 0
	} else {
		return 1
	}
}

// Return
// int = 0: success
//       1: failed
// string is the new book's unique key
func storeImage() (r int) {
	// Return value
	r = 0

	// Read file
	b, err := ioutil.ReadFile("Hydrangeas.jpg")
	if err != nil {
		fmt.Println(err, "in reading file")
		r = 1
		return
	}

	// Vernon debug
	fmt.Println("File length:", len(b))

	// Send request
	resp, err := http.Post(BookURL+"storeImage", "image/jpeg", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err)
		r = 1
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status, resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		// Get data from body
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err, "in reading body")
			r = 1
			return
		}
		fmt.Printf("%s\n", b)

		r = 1
		return
	}
	url, err := resp.Location()
	if err != nil {
		fmt.Println(err, "in getting location from response")
		return
	}
	fmt.Println("Location is", url)

	// Get key from URL
	// tokens := strings.Split(url.Path, "/")
	// var keyIndexInTokens int = 0
	// for i, v := range tokens {
	// 	if v == "books" {
	// 		keyIndexInTokens = i + 1
	// 	}
	// }
	// if keyIndexInTokens >= len(tokens) {
	// 	fmt.Println("Key is not given")
	// 	return
	// }
	// key = tokens[keyIndexInTokens]
	// if key == "" {
	// 	fmt.Println("Key is empty")
	// 	return
	// }
	return
}

func main() {
	storeImage()
	// // Random seed
	// rand.Seed(time.Now().Unix())

	// // Test suite
	// fmt.Println("========================================")
	// if storeTen() != 0 {
	// 	fmt.Println("Store books failed")
	// 	return
	// } else {
	// 	fmt.Println("Store 10 books")
	// }
	// fmt.Println("========================================")
	// r, key := storeBook()
	// if r != 0 {
	// 	fmt.Println("Store a book failed")
	// 	return
	// } else {
	// 	fmt.Println("Store a book in key", key)
	// }
	// fmt.Println("========================================")
	// queryAll()
	// fmt.Println("========================================")
	// queryBook()
	// fmt.Println("========================================")
	// if deleteBook(key) != 0 {
	// 	fmt.Println("Failed to delete book key", key)
	// 	return
	// } else {
	// 	fmt.Println("Delete book key", key)
	// }
	// fmt.Println("========================================")
	// if deleteAll() != 0 {
	// 	fmt.Println("Delete failed")
	// 	return
	// } else {
	// 	fmt.Println("Delete all")
	// }
}
