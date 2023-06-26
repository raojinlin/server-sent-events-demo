package main

import (
	"fmt"
	"net/http"
	"time"
)

type Event struct {
	Id      string
	Name    string
	Data    []string
	Retry   uint
	Comment string
}

func (e *Event) String() string {
	event := ""
	if e.Comment != "" {
		event += fmt.Sprintf(":%s\n", e.Comment)
	}

	if e.Id != "" {
		event += "id: " + e.Id + "\n"
	}

	if e.Name != "" {
		event += "event: " + e.Name + "\n"
	}

	for _, data := range e.Data {
		event += "data: " + data + "\n"
	}

	if e.Retry > 0 {
		event += fmt.Sprintf("retry: %d\n", e.Retry)
	}

	return event + "\n"
}

func (e *Event) Bytes() []byte {
	return []byte(e.String())
}

func stream(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	counter := 100
	req.Context().Deadline()

	ctx := req.Context()
	go func() {
		<-ctx.Done()
		// 客户端关闭连接，停止生成事件
		counter = -1
	}()

	for {
		if counter <= 0 {
			break
		}

		time.Sleep(1 * time.Second)
		event := Event{
			Name:    "message",
			Data:    []string{fmt.Sprintf("{\"type\": \"counter\", \"value\":%d}", counter)},
			Id:      fmt.Sprintf("%d", counter),
			Comment: "this is comment",
			Retry:   5000,
		}

		_, err := w.Write(event.Bytes())
		if err != nil {
			break
		}

		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		counter--
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
<html>
<head><title>events</title></head>
<body>
</body>
<script type="text/javascript">
  const sourceEvent = new EventSource('/stream');
  sourceEvent.onmessage = e => {
   const ele = document.createElement('div');
   ele.textContent = e.data
    document.body.append(ele);
  };
</script>
</html>
`)
}
func main() {
	http.HandleFunc("/events", stream)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
