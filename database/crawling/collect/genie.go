package collect

import (
	"fmt"

	"github.com/gocolly/colly"
)

func Genie() {
	fmt.Println("Genie Start :", now())

	c := colly.NewCollector(colly.Async(true))
	detail := c.Clone()
	store := make(map[string]Song)

	adds := 1
	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr.list", func(idx int, se *colly.HTMLElement) {
			songNo := se.Attr("songid")
			name := se.ChildText("td.info > a.title")
			artist := se.ChildText("td.info > a.artist")

			musicId := getMusic(name, artist)
			if musicId != nil {
				if !existOnChart("genie", songNo) {
					updateMusic("genie", songNo, *musicId)
				}
				insertRank(idx+adds, "genie", *musicId)
			} else {
				store[songNo] = Song{idx + adds, name, artist}
				detail.Visit("https://www.genie.co.kr/detail/songInfo?xgnm=" + songNo)
			}
		})
		adds += 50
	})

	detail.OnHTML("#body-content", func(e *colly.HTMLElement) {
		songNo := e.ChildAttr("#add_my_album_top", "songid")
		song := store[songNo]

		album := e.ChildText("ul > li:nth-child(2) > span.value > a")
		albumId := getAlbum(album, song.artist)
		if albumId == nil {
			image := "https:" + e.ChildAttr("div.photo-zone > a", "href")
			albumId = insertAlbum(album, song.artist, image)
		}

		lyric := e.ChildText("#pLyrics > p")

		musicId := insertMusic(song.name, *albumId, lyric, nil, &songNo, nil, nil, nil)
		insertRank(song.rank, "genie", *musicId)
	})

	c.Visit("https://www.genie.co.kr/chart/top200?ditc=D&rtm=N&pg=1")
	c.Wait()

	c.Visit("https://www.genie.co.kr/chart/top200?ditc=D&rtm=N&pg=2")
	c.Wait()

	fmt.Println("Genie End :", now())
}
