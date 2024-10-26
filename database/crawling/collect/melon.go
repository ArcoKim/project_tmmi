package collect

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Song struct {
	rank   int
	name   string
	artist string
}

func now() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func Melon() {
	fmt.Println("Melon Start :", now())

	c := colly.NewCollector(colly.Async(true))
	detail := c.Clone()
	store := make(map[string]Song)

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr.lst50, tr.lst100", func(idx int, se *colly.HTMLElement) {
			songNo := se.Attr("data-song-no")
			name := se.ChildText("div.rank01 > span > a")
			artist := se.ChildText("div.rank02 > a")

			musicId := getMusic(name, artist)
			if musicId != nil {
				if !existOnChart("melon", songNo) {
					updateMusic("melon", songNo, *musicId)
				}
				insertRank(idx+1, "melon", *musicId)
			} else {
				store[songNo] = Song{idx + 1, name, artist}
				detail.Visit("https://www.melon.com/song/detail.htm?songId=" + songNo)
			}
		})
	})

	detail.OnHTML("div#conts", func(e *colly.HTMLElement) {
		songNo := e.ChildAttr("#btnLike", "data-song-no")
		song := store[songNo]

		album := e.ChildText("div.meta > dl.list > dd > a")
		albumId := getAlbum(album, song.artist)
		if albumId == nil {
			image := e.ChildAttr("a.image_typeAll > img", "src")
			albumId = insertAlbum(album, song.artist, image)
		}

		lyric, _ := e.DOM.Find("div.lyric").Html()
		lyric = strings.Replace(lyric, "<br/>", "\n", -1)
		lyric = strings.Replace(lyric, "<!-- height:auto; 로 변경시, 확장됨 -->", "", 1)
		lyric = strings.TrimSpace(lyric)

		musicId := insertMusic(song.name, *albumId, lyric, &songNo, nil, nil, nil, nil)
		insertRank(song.rank, "melon", *musicId)
	})

	c.Visit("https://www.melon.com/chart/day/index.htm")
	c.Wait()

	fmt.Println("Melon End :", now())
}
