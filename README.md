# tooladapter

> **DEPRECATED**: This repository has been merged into [toolfoundation](https://github.com/jonwraymond/toolfoundation). Please use `github.com/jonwraymond/toolfoundation/adapter` instead.

## Migration

See [MIGRATION.md](MIGRATION.md) for detailed migration instructions.

### Quick Migration

Replace your imports:

```go
// Before
import "github.com/jonwraymond/tooladapter"
import "github.com/jonwraymond/tooladapter/adapters"

// After
import "github.com/jonwraymond/toolfoundation/adapter"
import "github.com/jonwraymond/toolfoundation/adapter/adapters"
```

## Why This Change?

The tooladapter functionality has been consolidated into toolfoundation as part of the ApertureStack unification effort. This provides:

- Single import path for foundational tool types
- Better cohesion between canonical types and adapters
- Simplified dependency management

## Timeline

- This repository will remain available for existing users
- No new features will be added
- Critical bug fixes only until v1.0 of toolfoundation

## License

See LICENSE file.
