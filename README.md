A V A T A R
===========

Create avatars based on names. The colors are chosen based on the first
character. You can save to disk ```ToDisk``` or send back over HTTP 
```ToHTTP```.

This fork is based off github.com/argylelabcoat/avatar, which is forked from andrew's https://github.com/ae0000/avatar. 

All this fork cares about is ditching the logging and `os.Exit` bits, instead returning `error` values. This makes it safer to use as a library.

Example
-------

```
package main

import (
	"github.com/argylelabcoat/avatar"
)

func main() {
	avatar.ToDisk("AE", "../ae.png")

    // Which is the same as
    avatar.ToDisk("Andrew Edwards", "../ae.png")
}


```
[![Example](https://raw.githubusercontent.com/argylelabcoat/avatar/master/ae.png)](https://raw.githubusercontent.com/argylelabcoat/avatar/master/ae.png)

You can pass in a single character as well

```
// Single initial as well..
avatar.ToDisk("Jet", "../j.png")
```
[![Example](https://raw.githubusercontent.com/argylelabcoat/avatar/master/j.png)](https://raw.githubusercontent.com/argylelabcoat/avatar/master/j.png)


HTTP example
------------

Using [go-chi](https://github.com/go-chi/chi) (highly recommended HTTP router)

```
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/argylelabcoat/avatar"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	// Get the png based on the initials, You would use it like this:
	//    <img src="http://localhost:3000/avatar/ae/png" width="150">
	r.Get("/avatar/{initials}.png", func(w http.ResponseWriter, r *http.Request) {
		initials := chi.URLParam(r, "initials")

		avatar.ToHTTP(initials, w)
	})

	http.ListenAndServe(":3000", r)
}

```
TODO
----

- [x] HTTP example
- [x] Caching
- [x] Custom colors
- [ ] Add unique colors that are missing (T-Z,0-9)