package main

import (
	"fmt"
	"sync"

	"simplepractice.com/wp-post-migrator/builder"
	"simplepractice.com/wp-post-migrator/migrators"
)

func main() {
	var wg sync.WaitGroup
	i := 0
	for i < 5 {
		builder.DeleteAllEntries(builder.MODEL_RESOURCE_POSTS)
		i++
	}
	i = 0
	for i < 4 {
		wg.Add(1)
		go migrators.MigrateResources("100", fmt.Sprintf("%d", i * 100), &wg)
		i++
	}
	wg.Wait()
}
