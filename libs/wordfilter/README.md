###Wordfilter (敏感文字过滤)
```go
package main

import (
	"fmt"

	"github.com/k4s/tea/libs/Wordfilter"
)

func main() {
	fmt.Println("中国共产党习近平主席总书记万岁，fuckyou!")
	//过滤前
	//中国共产党习近平主席总书记万岁，fuckyou!

	fmt.Println(Wordfilter.Wordfilter("中国共产党习近平主席总书记万岁.fuckyou!"))
	//过滤后
	//**********总书记万岁.**!
}

```