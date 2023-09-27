# URL Shortener

This is a URL shortening service implemented in Go. It allows you to shorten long URLs into short, easy-to-share links. Additionally, it provides features to track and display statistics on the usage of shortened URLs.

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Dependencies](#dependencies)

## Overview

The URL Shortener is built using the Go programming language and relies on the Redis database for storing and retrieving shortened URLs. It provides a web interface for shortening URLs and also supports features such as redirection and tracking of the top domains.

## Getting Started

To run the URL Shortener service locally, follow these steps:

1. **Clone the Repository:**

Clone this GitHub repository to your local machine:

```
git clone https://github.com/snehapatil1/url-shortener.git
```

2. **Install Dependencies:**

Make sure you have Go installed. Also, you need to install the required Go packages listed in the `go.mod` file. You can use the following command to install them:

```
go mod tidy
```

3. **Start the Redis Server:**

Ensure that you have a Redis server running on your machine. If not, you can download and install Redis from the official website: https://redis.io/download/.

4. **Configure Redis:**

Modify the Redis configuration in the `main.go` file if your Redis server is running on a different address or port. By default, it's set to `localhost:6379`.

5. **Build and Run:**

You can build and run the URL Shortener service using the following command:

```
go run main.go
```

The service will be accessible at `http://localhost:3030`.

## Usage

### Shortening URLs

1. Access the URL Shortener web interface at http://localhost:3030. The homepage will look something like this:

![Homepage](https://github.com/snehapatil1/url-shortener/blob/master/images/homepage.png)
2. Enter a long URL into the input field and click the "Shorten" button. On clicking Shorten button, you will get output as below:

![Shortened URL](https://github.com/snehapatil1/url-shortener/blob/master/images/shortened_url.png)
3. You will receive a shortened URL, which you can use to access the original long URL.

### Redirection

- To access the original long URL from a shortened URL, click on the given shortened URL or use the following format:
http://localhost:3030/redirectURL/{short_key}

### Tracking Top Domains

- To view the top domains statistics, access the following URL: http://localhost:3030/getTopDomains/

- This page will display the top 3 domains based on the number of times they have been shortened.
E.g.

![Top Domains](https://github.com/snehapatil1/url-shortener/blob/master/images/top_domains.png)

## Dependencies

- All dependencies are managed via Go Modules and are listed in the go.mod file.
