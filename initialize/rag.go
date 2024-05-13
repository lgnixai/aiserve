package initialize

import (
	"context"
	"fmt"
	"os"

	`github.com/gin-gonic/gin`
	"github.com/henomis/lingoose/assistant"
	embedder "github.com/henomis/lingoose/embedder/ollama"
	"github.com/henomis/lingoose/index"
	"github.com/henomis/lingoose/index/vectordb/jsondb"
	`github.com/henomis/lingoose/llm/ollama`
	"github.com/henomis/lingoose/rag"
	"github.com/henomis/lingoose/thread"
)

// download https://raw.githubusercontent.com/hwchase17/chat-your-data/master/state_of_the_union.txt

func AddRag(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.String(400, "Bad request", err)
		return
	}

	// 保存文件到指定位置
	dst := "uploads/" + file.Filename
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.String(500, "Internal server error")
		return
	}

	llm := ollama.New().WithEndpoint("http://localhost:11434/api").WithModel("hub/darkstorm2150/Ooh-Ollama:latest")

	r := rag.NewSubDocument(
		index.New(
			jsondb.New().WithPersist("./data/"+file.Filename+".json"),
			embedder.New(),
		),
		llm,
	).WithTopK(3)

	//_, err = os.Stat(file.Filename + "db.json")

	err = r.AddSources(context.Background(), dst)
	if err != nil {
		panic(err)
	}

	c.String(200, fmt.Sprintf("File %s uploaded successfully", file.Filename))

}
func ragget(c *gin.Context) {
	llm := ollama.New().WithEndpoint("http://localhost:11434/api").WithModel("hub/darkstorm2150/Ooh-Ollama:latest")

	r := rag.NewSubDocument(
		index.New(
			jsondb.New().WithPersist("db.json"),
			embedder.New(),
		),
		llm,
	).WithTopK(3)

	_, err := os.Stat("db.json")
	if os.IsNotExist(err) {
		err = r.AddSources(context.Background(), "./1.txt")
		if err != nil {
			panic(err)
		}
	}

	a := assistant.New(
		llm.WithTemperature(0),
	).WithRAG(r).WithThread(
		thread.New().AddMessages(
			thread.NewUserMessage().AddContent(
				thread.NewTextContent("用中文回答，leven 居住在哪个国 家"),
			),
		),
	)

	err = a.Run(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("----")
	fmt.Println(a.Thread())
	fmt.Println("----")
}
