# md2jira

A renderer extension for [yuin/goldmark](https://github.com/yuin/goldmark) to convert CommonMark to Jira text formatting

## Usage

Import it into your own code

```go
import "github.com/bradfordboyle/md2jira/jira"

renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(jira.NewRenderer(), 1000)))
```

Use it from the shell

```console
h1. md2jira

A renderer extension for [yuin/goldmark|https://github.com/yuin/goldmark] to convert CommonMark to Jira text formatting

h2. Usage

Import it into your own code

{code:go}
import "github.com/bradfordboyle/go-md2jira/jira"

renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(jira.NewRenderer(), 1000)))
{code}
```
