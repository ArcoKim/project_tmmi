package collect

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-rod/rod"
)

func Flo() {
	fmt.Println("FLO Start :", now())

	browser := rod.New().MustConnect()
	defer browser.Close()

	chart := browser.MustPage("https://www.music-flo.com/browse")
	chart.MustWaitLoad()
	chart.MustElement(".btn_list_more").MustClick()

	rows := chart.MustElements("tbody > tr")

	var wg sync.WaitGroup
	for idx, row := range rows {
		wg.Add(1)
		go func(idx int, row *rod.Element) {
			defer wg.Done()

			var tmap map[string]interface{}
			temp := *row.MustElement("button.btn_listen").MustAttribute("data-rake")
			json.Unmarshal([]byte(temp), &tmap)

			songNo := tmap["trackId"].(string)
			name := row.MustElement("strong.tit__text").MustText()
			artist := row.MustElement("span.artist__link").MustText()

			musicId := getMusic(name, artist)
			if musicId != nil {
				if !existOnChart("flo", songNo) {
					updateMusic("flo", songNo, *musicId)
				}
				insertRank(idx+1, "flo", *musicId)
				return
			}

			song := browser.MustPage("https://www.music-flo.com/detail/track/" + songNo + "/details")
			song.MustWaitLoad()

			album := song.MustElement("p.album > a").MustText()
			albumId := getAlbum(album, artist)
			if albumId == nil {
				image := song.MustElement("div.link_thumbnail > img").MustAttribute("src")
				albumId = insertAlbum(album, artist, *image)
			}

			lyric := song.MustElement("div.lyrics").MustText()

			musicId = insertMusic(name, *albumId, lyric, nil, nil, &songNo, nil, nil)
			insertRank(idx+1, "flo", *musicId)
		}(idx, row)
	}
	wg.Wait()

	fmt.Println("FLO End :", now())
}
