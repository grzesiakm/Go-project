# Go-project
Web application to get the cheapest flight using web-scraping and golang

## Opis projektu

Celem projektu jest stworzenie prostej aplikacji webowej przy użyciu języka Go, która będzie używać web-scrapingu. 
Zadaniem aplikacji będzie wyszukanie użytkownikowi najtańszego lotu ograniczonego parametrami, które sam podał. Dane będą eksportowane w postaci plików JSON, które zostaną odpowiednio przygotowane do wyświetlenia.

W projekcie planujemy użyć:
- https://pkg.go.dev/github.com/PuerkitoBio/goquery
- https://pkg.go.dev/github.com/playwright-community/playwright-go
- https://pkg.go.dev/encoding/json
- https://pkg.go.dev/net/http
- https://pkg.go.dev/time

## Jak uruchomić projekt?

W katalogu **scraping** wykonaj poniższe komendy:
- `go get github.com/playwright-community/playwright-go`
- `go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps`
- `go run .`
 
## Zespół

- Marta Grzesiak
- Aleksandra Kuś
