package collect

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func Bugs() {
	fmt.Println("Bugs Start :", now())

	c := colly.NewCollector(colly.Async(true))
	detail := c.Clone()
	store := make(map[string]Song)

	c.OnHTML("table.trackList > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(idx int, se *colly.HTMLElement) {
			songNo := se.Attr("trackid")
			name := se.ChildText("p.title > a")
			artist := se.ChildText("p.artist > a")

			musicId := getMusic(name, artist)
			if musicId != nil {
				if !existOnChart("bugs", songNo) {
					updateMusic("bugs", songNo, *musicId)
				}
				insertRank(idx+1, "bugs", *musicId)
			} else {
				store[songNo] = Song{idx + 1, name, artist}
				detail.Visit("https://music.bugs.co.kr/track/" + songNo)
			}
		})
	})

	detail.OnHTML("#container", func(e *colly.HTMLElement) {
		url := strings.Split(e.Request.URL.String(), "/")
		songNo := url[len(url)-1]
		song := store[songNo]

		album := e.ChildText("div.basicInfo > table > tbody > tr:nth-child(3) > td > a")
		if album == "" {
			album = e.ChildText("div.basicInfo > table > tbody > tr:nth-child(2) > td > a")
		}

		albumId := getAlbum(album, song.artist)
		if albumId == nil {
			image := e.ChildAttr("div.photos > ul > li.big > a > img", "src")
			albumId = insertAlbum(album, song.artist, image)
		}

		lyric := e.ChildText("xmp")

		musicId := insertMusic(song.name, *albumId, lyric, nil, nil, nil, &songNo, nil)
		insertRank(song.rank, "bugs", *musicId)
	})

	c.Visit("https://music.bugs.co.kr/chart/track/day/total")
	c.Wait()

	fmt.Println("Bugs End :", now())
}
