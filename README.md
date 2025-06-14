# Go Language Learning

This repository contains my journey learning the Go programming language through hands-on examples and projects.

## What's Inside

- Basic Go syntax and concepts
- Struct examples with practical applications
- Simple programs demonstrating Go features

## Understanding Structs in Go

A **struct** is a collection of fields grouped together to represent a single concept or entity.

### Basic Syntax
```go
type Person struct {
    Name     string
    Age      int
    Health   int
}
```

### When to Use Structs

- **Grouping related data**: Instead of passing multiple separate variables, bundle them together
- **Modeling real-world entities**: People, cars, game characters, etc.
- **Creating custom types**: When you need something more complex than basic types (int, string, bool)
- **Organizing code**: Structs can have methods attached to them

### Example Use Cases
```go
// Game character
type Player struct {
    Name     string
    Health   int
    Strength int
}

// Configuration settings
type Config struct {
    Host string
    Port int
    Debug bool
}

// Database record
type User struct {
    ID    int
    Email string
    CreatedAt time.Time
}
```

**Rule of thumb**: If you find yourself passing the same 3+ variables to multiple functions, consider making them a struct.

## Running the Code

```bash
go run cmd/myapp/filename.go
```
```

This README is concise, beginner-friendly, and gives a clear understanding of structs and their practical applications!