<img width="150" height="150" alt="procyon-logo" src="https://go.codnect.io/procyon-logo.png" />

# Procyon Framework

Procyon is an application framework for Go that helps you build modern,
production-ready applications with minimal boilerplate.

## Installation

To install Procyon, run the following command:

```bash
go get -u go.codnect.io/procyon/...
```

## Overview

Building production-ready Go applications often requires combining multiple
libraries and integrating them into a cohesive application.

As applications grow, managing infrastructure code and keeping these components
working together can introduce unnecessary complexity and boilerplate.

Procyon simplifies application development by providing the essential
capabilities needed to build modern Go applications, allowing you to focus on
your business logic instead of infrastructure.

## Quick Start

Create a new application:

```go
package main

import (
    "os"
	
    "go.codnect.io/procyon"
)

func main() {
    if err := procyon.Run(); err != nil {
        os.Exit(1)
    }
}
```

## Documentation

Read the full documentation at https://go.codnect.io/procyon.

## License

Procyon Framework is released under version 2.0 of the Apache License.