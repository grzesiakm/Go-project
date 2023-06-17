# Go-project
Web application to get the cheapest flight using web-scraping and golang

## Opis projektu

Celem projektu jest stworzenie prostej aplikacji webowej przy użyciu języka Go, która będzie używać web-scrapingu. 
Zadaniem aplikacji będzie wyszukanie użytkownikowi najtańszego lotu ograniczonego parametrami, które sam podał. Dane będą wyszukiwane w 5 wybranych liniach lotniczych, a następnie wyświetlane na stworzonej stronie.

W projekcie planujemy użyć:
- https://pkg.go.dev/github.com/PuerkitoBio/goquery
- https://pkg.go.dev/github.com/playwright-community/playwright-go
- https://pkg.go.dev/net/http

## Jak uruchomić projekt?

W katalogu **scraping** wykonaj poniższe komendy:
- `go build`
- `go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps`
- `go run main.go` lub uruchom stworzony plik wykonywalny (np. main.exe)

Następnie w przeglądarce otwórz http://127.0.0.1:8080/
 
## Zespół

- Marta Grzesiak
- Aleksandra Kuś
