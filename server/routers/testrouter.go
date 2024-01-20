package routers

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TestGroup(rg *gin.RouterGroup) {
	rg.GET("/a", func(ctx *gin.Context) {
		ctx.JSON(200, "hello world")
	})

	// la pq eu consigo defenir o tamanho do range que quero
	// eu preciso de limitar ainda a velocidade
	// eu posso tentar usar o buffer de audio e sincronizarlos
	// stream e sempre mais rapido tmb
	// mas eu fico sem habilidades de seek eu acho 

	// eu tenho de ver isto 

	/* 
	rg.GET("/stream", func(ctx *gin.Context) {
		chanStream := make(chan []byte)
		ctx.Header("Content-Type", "audio/mp3")
		ctx.Header("Content-Disposition", "inline; filename=audio")

		go func() {
			file, err := os.ReadFile("./audio/audio.opus")
			//ctx.Header("Content-Length", string(bytes.Count(file, file)))

			chanStream <- file
			if err != nil {
				ctx.String(500, "Error streaming file %s", err)
			}
			defer close(chanStream)
		}()

		ctx.Stream(func(w io.Writer) bool {
			for buf := range chanStream {
				w.Write(buf)
			}
			return false
		})
	}) 
	*/

	rg.GET("/stream2", func(ctx *gin.Context) {
		// de alguma maneira isto funciona...
		ctx.Header("Content-Type", "audio/ogg")
		ctx.Header("Content-Disposition", "inline; filename=audio.opus")

		data, err := os.ReadFile("./audio/audio.opus")
		if err != nil {
			ctx.String(500, "Error opening file: %s", err)
			return
		}

		ctx.Header("Content-Length", strconv.Itoa(len(data)))

		ctx.Stream(func(w io.Writer) bool {
			var totalsent int64
			var chunkSize int64 = 512 << 10

			for totalsent < int64(len(data)) {
				remaining := int64(len(data)) - totalsent
				if remaining < chunkSize {
					chunkSize = remaining
				}

				n, err := w.Write(data[totalsent : totalsent+chunkSize])
				if err != nil {
					fmt.Println("[stream] Err: ", err)
					return false
				}

				totalsent += int64(n)
			}

			return false
		})
	})


	// think abt this
	// entendi como usar isto ig
	rg.StaticFile("/audio", "./audio/audio.opus")
}