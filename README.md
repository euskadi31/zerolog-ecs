# Elastic Common Schema (ECS) for Zerolog

Apply [ECS](https://www.elastic.co/guide/en/ecs/1.12/index.html) json structure to [Zerolog](https://github.com/rs/zerolog) library.

Example:

```
package main

import (
    "github.com/rs/zerolog/log"
    "github.com/euskadi31/zerolog-ecs"
    ecsaws "github.com/euskadi31/zerolog-ecs/feature/aws"
)

func main() {
    zerologecs.Configure(zerologecs.WithServiceName("hello-bin"))
    ecsaws.Configure()

    log.Info().Msg("hello")
}

```

## License

zerolog-ecs is licensed under the [MIT license](LICENSE.md).
