# Link Parser

https://github.com/gophercises/link

Finds all links (HTML anchor tags), inclusive of the text in them.

- Works on local files!
- Works on public websites!
- Handles nested links!
- Ignores comments!

```sh
go run . # runs it on "https://jinwei.dev"
go run . --url https://google.com.sg # runs it on "https://google.com.sg"

go run . --file ex1.html # runs it on local file "ex1.html"
```
