# Quiet HN

https://github.com/gophercises/quiet_hn

Concurrently gets top stories from HackerNews.

- Works concurrently!
- **Always** gets exactly the top `num_stories` amount of stories! \*\*
- Order of stories is exactly the same as HackerNews!
- Results are cached!

```sh
go run .
go run . --port 3030       # default: 3000
go run . --num_stories 100 # default: 30
```

\*\* unless there is less than `num_stories` stories in the top 450 stories / links.
