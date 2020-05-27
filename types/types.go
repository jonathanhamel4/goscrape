package types

type Movie struct {
	Title string `selector:".title_wrapper h1"`
	Country string `selector:"#titleDetails > .txt-block:first-of-type a"`
	Language string `selector:"#titleDetails > .txt-block:nth-of-type(2) a"`
	Rating string `selector:"span[itemprop='ratingValue']"`
	RatingCount string `selector:"span[itemprop='ratingCount']"`
	Summary string `selector:"div.summary_text"`
	StoryLine string `selector:"#titleStoryLine div > p > span"`
	Runtime string `selector:"#titleDetails > .txt-block:nth-of-type(8) > time"`
	Genres []string `selector:".title_wrapper a[href]:not([title])"`
	ImdbUrl string
}

type Genre string
