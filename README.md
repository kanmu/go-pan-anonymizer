# go-pan-anonymizer

Anonymize PAN (Primary Account Number = Credit card number) in text.

## Example

```
package main

import (
	"io"
	"os"

	anonymizer "github.com/kanmu/go-pan-anonymizer"
	"golang.org/x/text/transform"
)

func main() {
	io.Copy(os.Stdout, transform.NewReader(os.Stdin, anonymizer.DefaultAnonymizer()))
}
```
