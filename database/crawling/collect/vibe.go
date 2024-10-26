package collect

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-rod/rod"
)

func Vibe() {
	fmt.Println("Vibe Start :", now())

	browser := rod.New().MustConnect()
	defer browser.Close()

	chart := browser.MustPage("https://vibe.naver.com/chart/total")
	chart.MustWaitLoad()

	rows := chart.MustElements("tbody > tr")

	var wg sync.WaitGroup
	for idx, row := range rows {
		wg.Add(1)
		go func(idx int, row *rod.Element) {
			defer wg.Done()

			songNo := *row.MustElement("a.link_text > span").MustAttribute("id")
			name := row.MustElement("a.link_text > span").MustText()
			artist := row.MustElement("a.link_artist > span").MustText()

			musicId := getMusic(name, artist)
			if musicId != nil {
				if !existOnChart("vibe", songNo) {
					updateMusic("vibe", songNo, *musicId)
				}
				insertRank(idx+1, "vibe", *musicId)
				return
			}

			song := browser.MustPage("https://vibe.naver.com/track/" + songNo)
			song.MustWaitLoad()

			album := strings.TrimSpace(song.MustElement("div.text_area > div > a").MustText()[9:])
			albumId := getAlbum(album, artist)
			if albumId == nil {
				image := song.MustElement("div.summary_thumb > img").MustAttribute("src")
				albumId = insertAlbum(album, artist, *image)
			}

			lyric := song.MustElement("div.lyrics > p").MustText()

			musicId = insertMusic(name, *albumId, lyric, nil, nil, nil, nil, &songNo)
			insertRank(idx+1, "vibe", *musicId)
		}(idx, row)
	}
	wg.Wait()

	fmt.Println("Vibe End :", now())
}
