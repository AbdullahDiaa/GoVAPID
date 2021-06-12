# GOVAPID

[![GoDoc][godoc-image]][godoc-url]
[![codecov](https://codecov.io/gh/abdullahdiaa/govapid/branch/main/graph/badge.svg?token=XY7LZ4XCK3)](https://codecov.io/gh/abdullahdiaa/govapid)
[![Build Status](https://www.travis-ci.com/AbdullahDiaa/GoVAPID.svg?branch=main)](https://www.travis-ci.com/AbdullahDiaa/GoVAPID)
[![Go Report Card](https://goreportcard.com/badge/github.com/AbdullahDiaa/govapid)](https://goreportcard.com/report/github.com/AbdullahDiaa/govapid)

> A micro-package to generate VAPID keys required for sending web push notifications, this package uses standard library dependencies only.

## Usage

```go
package main

import (
	"fmt"
	"github.com/AbdullahDiaa/govapid"
)

func main() {
	VAPIDkeys, err := govapid.GenerateVAPID()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Public Key:%s\nPrivate Key:%s", VAPIDkeys.Public, VAPIDkeys.Private)
}

```

## Documentation

You can view detailed documentation here: [GoDoc][godoc-url].

## Contributing

There are many ways to contribute:
- Fix and [report bugs](https://github.com/AbdullahDiaa/GoVAPID/issues/new)
- [Improve documentation](https://github.com/AbdullahDiaa/GoVAPID/issues?q=is%3Aopen+label%3Adocumentation)
- [Review code and feature proposals](https://github.com/AbdullahDiaa/GoVAPID/pulls)


## Changelog

View the [changelog](/CHANGELOG.md) for the latest updates and changes by
version.

## License

[Apache License 2.0][licence-url]

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

[godoc-image]: https://godoc.org/github.com/AbdullahDiaa/GoVAPID?status.svg
[godoc-url]: https://godoc.org/github.com/AbdullahDiaa/GoVAPID
[licence-url]: https://github.com/AbdullahDiaa/GoVAPID/blob/main/LICENSE