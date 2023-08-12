# Cappadocia

A simple filesystem watcher that runs a command when a file changes.

Motivating use case: I have a GraphQL schema that's defined in Python using [strawberry](https://strawberry.rocks). I want to automatically generate a GraphQL schema from this Python file whenever it changes, so that it can be used by a React frontend easily.

To do so, this is just


```bash
cappadocia watch schema.py strawberry "export-schema package.module:schema > schema.graphql"

# simplest example
cappadocia watch "*.md" "echo" "hello world"
> Watching 1 files matching *.md
> hello world
> hello world

# also equivalent
cappadocia watch "*.md" echo hello world
```