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

	/*
	
	com stream o firefox consegue imediatamente ouvir a musica
	o que é bom pq se tiver com pouca rede a musica continua a dar a medida
	a medida q carrega

	sim consigo fazer o mesmo com content-lenght e ranges
	mas n me vou dar ao trabalho disso... 
	
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

	/*
	uso api https://developer.mozilla.org/en-US/docs/Web/API/NetworkInformation
	para caso que seje algum mobile data or smth

	e passar para /stream

	tmb n uso firefox (no tele) por isso fds ig  
	posso usar essa api
	*/

	rg.GET("/stream", func(ctx *gin.Context) {
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