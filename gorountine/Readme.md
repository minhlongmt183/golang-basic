# Goroutines
---  
## Kien thuc co ban
---  
Tiến trình (process): có thể hiểu đơn giản là một chương trình đang chạy 
trong máy tính. Mỗi tiến trình sẽ có một luồng chính (main thread) để chạy 
chương trình và được hệ điều hành cấp pháp cho một không gian bộ nhớ nhất định. 
Khi main thread ngừng hoạt động đồng nghĩa với việc chương trình bị tắt.  
Luồng (thread): thread hay còn được gọi là tiểu trình là một luồng trong 
tiến trình đang chạy. Các luồng được chạy song song trong tiến trình và có 
thể truy cập đến vùng nhớ được cung cấp bởi tiến trình. Những thread sẽ được 
cấp pháp riêng một vùng nhớ stack để lưu trữ biến riêng của thread đó.

## Goroutines vs system threads
---  
- Golang sử dụng goroutine để xử lý đồng thời nhiều tác vụ.  
- Goroutines là hàm hoặc phương thức chạy đồng thời với các hàm hoặc phương thức khác.  
- Khởi tạo goroutines sẽ tốn ít chi phí hơn khởi tạo thread so với các ngôn 
ngữ khác.  
- Trên thread sẽ có một kích thước vùng nhớ stack cố định. Vùng nhớ này 
  chủ yếu được sử dụng để lưu trữ những tham số, biến cục bộ và địa chỉ trả 
  về khi chúng ta gọi hàm.

- Vì kích thước cố định của stack nên đẫn đến hai vấn đề:
  - Gặp hiện tượng stack overflow với những chương trình gọi hàm đệ quy sâu.  
  - Lãng phí vùng nhớ đối với chương trình đơn giản.  
  
Với Goroutines thì vẫn đề này đã được khắc phục bằng cách cấp pháp linh hoạt 
vùng nhớ stack:  
- Một Goroutines sẽ được bắt đầu bằng một vùng nhớ nhỏ.
- Khi chương trình chạy nếu không gian stack hiện tại không đủ, Goroutines sẽ tự động tăng không gian stack 
- Do chi phí việc khơi tạo nhỏ nên ta có thể dễ dàng giải phóng hàng ngàn goroutines

| Thread                                                                                | Goroutines                                                                | 
|---------------------------------------------------------------------------------------|---------------------------------------------------------------------------| 
| Thread được quản lý bởi hệ điều hành và phụ thuộc vào số nhân của CPU                 | Goroutines được quản lý bởi go runtime và không phụ thuộc vào số nhân CPU |
| Thread có kích cỡ vùng nhớ stack cố định                                              | 	Goroutines có kích cỡ vùng nhớ stack tùy theo chương trình               |  
| Giao tiếp giữa các thread khá khó. Có đỗ trễ lớn trong việc tương tác giữa các thread | Goroutines sử dụng channels để tương tác với nhau với độ trễ thấp         |
| Thread có định danh                                                                   | Goroutine không có định danh |                                             |
| Khơi tạo và giải phóng thread tốn nhiều thời gian | Goroutines được khởi tạo và giải phóng bởi go runtime nên rất nhanh |

## sync.WaitGroup  
---  
```go
package main
import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	fmt.Println("Application start")

	wg.Add(1)
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Goroutines: ", i)
		}

		wg.Done()
	}()

	fmt.Println("Application end")
	wg.Wait()
}

```
Đầu tiên ta tiến hành tạo một `sync.WaitGroup` gọi hàm `Add(1)` vì ở đây mình 
chỉ gọi 1 goroutines nên sẽ thêm 1 nếu chương trình cần nhiều goroutines hơn 
thì các bạn sẽ thêm số lượng tương ứng vào WaitGroup. Sau đó trong hàm 
goroutines mình sẽ gọi đến hàm Done() để báo hiệu là goroutines đã chạy xong 
và ở hàm main mình có thêm hàm `Wait()` để chương trình chờ cho goroutines 
chạy xong rồi mới kết thúc.

## Channel
---  
Channel trong Go là một đường ống kết nối các goroutines để chúng có thể 
chia sẻ dữ liệu cho nhau. Để tạo channel ta sử dụng cú pháp:  
```go
make(chan <type>)
```  
- Gửi dữ liệu vào channel:  
```go
channelName <-
```  
- Nhận dữ liệu từ channel:
```go
<- channelName
```
- Mặc định quá trình gửi và nhận giữa các goroutines sẽ bị block đến khi cả 
  2 goroutines đã sẵn sàng để gửi và nhận.  
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool)

	fmt.Println("Application start")

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 5; i++ {
			fmt.Println("Goroutines: ", i)
		}
		done <- true
	}()

	fmt.Println("Application end")
	<-done
}
```  
Ở đây mình có tạo một channel trên là `done` có kiểu boolean trong goroutine 
mình có truyền một tín hiệu `true` vào channel thể hiện là goroutine đã chạy 
xong và ở hàm main mình nhận dữ liệu đó ở cuối chương trình để đợi goroutine 
đó chạy xong rồi mới kết thúc.  

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool)

	fmt.Println("Application start")

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Goroutines: ", i)
		}
	}()

	fmt.Println("Application end")
	time.Sleep(time.Second)
	<-done
}
```  
- Ở đây xuất hiện một lỗi đó là `fatal error: all goroutines are asleep - 
  deadlock!.`  
- `"quá trình gửi và nhận giữa các goroutines sẽ bị block đến khi cả 2 
   goroutines đã sẵn sàng để gửi và nhận"` nên ở đây mình chỉ nhận dữ liệu ở 
   hàm main mà trong goroutine kia mình không gửi gì qua channel nên chương 
   trình sẽ bị lỗi. Tương tự như vậy:  
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan string)

	fmt.Println("Application start")

	done <- "Done"
    fmt.Println(<-done)
	fmt.Println("Application end")
	time.Sleep(time.Second)
}

```  
- Khi chạy chương trình ta vẫn gặp đúng lỗi đó vì mặc dù ta đã thêm đầu nhận 
  của channel rồi nhưng code được chạy tuần tự khi gặp dòng `done <- "Done"` 
  thì chương trình đã gặp lỗi rồi nên channel vẫn chưa có đầu nhận dữ liệu.  
- Ta có thể giải quyết vấn đề này bằng channel buffering. Ta có đoạn code 
  như sau:
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan string, 1)

	fmt.Println("Application start")

	done <- "Done"

	fmt.Println("Application end")
	time.Sleep(time.Second)
}

```  
- khởi tạo channel buffering bằng cách truyền tham số thứ 2. Channel 
buffering cho phép ta có thể giới hạn giá trị mà channel nhận mà không cần 
  đầu nhận tương ứng của giá trị đó. 
- Ngoài ra nếu chúng ta sử dụng channel như một tham số của hàm thì ta có 
  thể định nghĩa xem tham số channel đó chỉ nhận hay gửi dữ liệu như sau:  

```go
package main

import (
	"fmt"
	"time"
)

func sendValue(number string, channel chan<- string) {
	for {
		channel <- number
	}
}

func receiveValue(channel <-chan string) {
	for v := range channel {
		fmt.Println(v)
	}
}
func main() {
	channel := make(chan string, 64)
	go sendValue("Hello", channel)
	go sendValue("Edisc", channel)

	go receiveValue(channel)
	time.Sleep(time.Second)
}

```  
Ouput bạn sẽ thấy chữ `Hello` và `Edisc` sẽ được in ra liên tục. Ở đây cú pháp 
channel `chan<- string` thể hiện tham số channel này chỉ dùng để gửi dữ liệu 
còn channel <-chan string thể hiện chỉ để nhận dữ liệu.