package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var mu sync.Mutex
var domainCountMap = make(map[string]int)

func setupRoutes() {
	http.HandleFunc("/", homePageForm)
	http.HandleFunc("/shortenURL", shortenURL)
	http.HandleFunc("/redirectURL/", redirectURL)
}

func main() {
	setupRoutes()

	fmt.Println("URL Shortener is running on :3030")
	http.ListenAndServe(":3030", nil)
}

func homePageForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shortenURL", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title><center>URL Shortener</center></title>
		</head>
		<body>
			<h2><center>URL Shortener</center></h2>
			<form method="post" action="/shortenURL" align="center">
				<input type="url" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form>
		</body>
		</html>
	`)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	_, err := RedisClient().Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	recordDomain(originalURL)

	shortKey := ""

	cachedValue, err2 := RedisClient().Get(originalURL).Result()
	if err2 != nil {
		log.Println(err2)
	}

	if cachedValue == "" {
		shortKey = generateShortKey()
		err1 := RedisClient().Set(originalURL, shortKey, 0).Err()
		if err1 != nil {
			return
		}
		RedisClient().Set(shortKey, originalURL, 0)

	} else {
		shortKey = cachedValue
	}

	shortenedURL := fmt.Sprintf("http://localhost:3030/short/%s", shortKey)
	shortURL := fmt.Sprintf("spurl/%s", shortKey)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title><center>URL Shortener</center></title>
		</head>
		<body>
			<h2><center>URL Shortener</center></h2>
			<p><center>Original URL: `, originalURL, `</center></p>
			<p><center>Shortened URL (Click on this link to visit Original URL): <a href="`, shortenedURL, `">`, shortURL, `</a></center></p>
		</body>
		</html>
	`)

}

func RedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return redisClient
}

func recordDomain(longURL string) {
	domain, _ := extractDomain(longURL)
	mu.Lock()
	defer mu.Unlock()
	domainCountMap[domain]++
}

func extractDomain(URL string) (string, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	hostname := parsedURL.Hostname()
	parts := strings.Split(hostname, ".")
	domainName := ""
	if len(parts) >= 2 {
		domainName = parts[len(parts)-2]
	} else {
		domainName = hostname
	}
	domainName = capitalizeFirstLetter(domainName)
	return domainName, nil

}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	strShortKey := string(shortKey)
	alreadyUsed, _ := RedisClient().Get(strShortKey).Result()
	if alreadyUsed != "" {
		generateShortKey()
	}
	return strShortKey
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/redirectURL/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	originalURL, found := RedisClient().Get(shortKey).Result()
	if found != nil {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
