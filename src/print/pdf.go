package print

import (
	"fmt"
	// "log"
	"math"
	"path"
	// "github.com/fogleman/gg"
	"github.com/fogleman/gg"
	"github.com/jung-kurt/gofpdf"
	"github.com/uconn-ling/openHouseMap/src/data"
)

var fileName string = "./output/pdf/openHouseMap.pdf"
var inputDataFolder string = path.Join("./inputData/")
var hashPicturesFolder string = path.Join(inputDataFolder, "pngs")

var imagesPerLine = 3
var linesPerPage = 4

var letterWidthMM = 215.9
var letterheightMM = 279.4

var borderLeftRight = 10.0
var borderTopBottom = 10.0

var heightOfHeader = 10.0
var heightOfPersons = 35.0
var heightOfName = 10.0
var countryTipaScalingFactor, personTipaScalingFactor float64

var verticalBorderTopBottomPicture = 3.0
var spacebetweenPictureAndImg = 2.0

var heightOfPersonsInklNameAndSpace = (heightOfPersons + (1 * verticalBorderTopBottomPicture) + heightOfName + spacebetweenPictureAndImg)

var headerBorder = 2.0

var borderBetweenFlagAndCountryName = 5.0

var verticalBoxBorder = 3.0

var countryNameFlagHeight = (heightOfHeader - (2 * headerBorder))

var defaultBoxWidth = (letterWidthMM - (2 * borderLeftRight))

var singleImgBoxWidth = defaultBoxWidth / 3

type Page struct {
	Countries []string
}

func CreatePdf(ct map[string]data.Country) {

	countryTipaScalingFactor, personTipaScalingFactor = getNameScalingFactor(ct)
	pages := optimalSort(ct)

	fmt.Printf("%v", pages)

	pdf := gofpdf.New("P", "mm", "Letter", "")
	for _, page := range pages {
		xStart := borderLeftRight
		yStart := borderTopBottom
		pdf.AddPage()

		// pdf.Line((letterWidthMM / 2), 0.0, (letterWidthMM / 2), letterheightMM)
		for idx, country := range page.Countries {

			/* ENDONYM */
			countryImgPath := path.Join(hashPicturesFolder, ct[country].Endonym.Hash+".png")
			newCountryNameWidth := ct[country].Endonym.Rendered.Width * (countryNameFlagHeight / ct[country].Endonym.Rendered.Height)
			// newCountryNameWidth := ct[country].Endonym.Rendered.Width * countryTipaScalingFactor
			newCountryNameHeight := countryNameFlagHeight
			// newCountryNameHeight := ct[country].Endonym.Rendered.Height * countryTipaScalingFactor

			/* FLAG */
			countryFlagPath := path.Join(inputDataFolder, "country_"+country, ct[country].Flag.Path)
			flagIm, err := gg.LoadImage(countryFlagPath)
			if err != nil {
				panic(err)
			}
			flagIw, flagIh := flagIm.Bounds().Dx(), flagIm.Bounds().Dy()
			newCountryFlagWidth := (float64(flagIw) * countryNameFlagHeight) / float64(flagIh)

			/* HEADER */
			minboxWidth := (4 * borderBetweenFlagAndCountryName) + (2 * newCountryFlagWidth) + newCountryNameWidth

			lines := float64(calcLines(ct[country].PersonCount, imagesPerLine))

			calcX := float64(xStart)
			calcY := float64(yStart + ((heightOfHeader + heightOfPersonsInklNameAndSpace + verticalBorderTopBottomPicture + verticalBoxBorder) * float64(idx)))

			// fmt.Println(fmt.Sprintf("%f", heightOfPersonsInklNameAndSpace))
			// fmt.Println(fmt.Sprintf("%f", calcY))

			boxheight := (heightOfHeader + verticalBorderTopBottomPicture + (lines * heightOfPersonsInklNameAndSpace))
			boxWidth := defaultBoxWidth
			if lines == 1 {
				boxWidth = (float64(ct[country].PersonCount) * singleImgBoxWidth)
			}

			if boxWidth < minboxWidth {
				boxWidth = minboxWidth
			}
			// fmt.Println(country)
			// fmt.Println(fmt.Sprintf("%f", boxheight))
			countryNameStartX := ((calcX + (boxWidth / 2)) - (newCountryNameWidth / 2))
			countryNameStartY := calcY + headerBorder
			// countryNameStartY := calcY + ((heightOfHeader - (2*headerBorder)) / 2) - (newCountryNameHeight / 2)

			pdf.SetLineWidth(0.5)
			pdf.SetDrawColor(100, 100, 100)

			pdf.ClipRoundedRect(calcX, calcY, boxWidth, boxheight, 2, true)
			pdf.SetFillColor(0, 0, 0)
			pdf.Rect(calcX, calcY, (letterWidthMM - (2 * borderLeftRight)), heightOfHeader, "F")
			pdf.Image(countryImgPath,
				countryNameStartX,
				countryNameStartY,
				newCountryNameWidth,
				newCountryNameHeight,
				false, "", 0, "")
			pdf.Image(countryFlagPath,
				(countryNameStartX - newCountryFlagWidth - borderBetweenFlagAndCountryName),
				(calcY + headerBorder),
				newCountryFlagWidth,
				countryNameFlagHeight,
				false, "", 0, "")
			pdf.Image(countryFlagPath,
				(countryNameStartX + newCountryNameWidth + borderBetweenFlagAndCountryName),
				(calcY + headerBorder),
				newCountryFlagWidth,
				countryNameFlagHeight,
				false, "", 0, "")

			count := 0
			personCount := 0
			remainder := ct[country].PersonCount % imagesPerLine
			persons := sortPersons(ct[country].People)

			/* PEOPLE */
			for _, personKey := range persons {
				personMap := ct[country].People[personKey]
				personLines := calcLines((count + 1), imagesPerLine)

				personStartX := calcX + (float64(personCount) * singleImgBoxWidth)
				// fmt.Println(country)
				// fmt.Println(fmt.Sprintf("%f", lines))
				// fmt.Println(fmt.Sprintf("%f", float64(personLines)))
				if lines == float64(personLines) {
					// fmt.Println(strconv.Itoa(ct[country].PersonCount))
					// fmt.Println(strconv.Itoa(count))
					if remainder == 1 {
						personStartX = calcX + (boxWidth / 2) - (singleImgBoxWidth / 2)
					} else if remainder == 2 {
						personStartX = calcX + (boxWidth / 2) - singleImgBoxWidth + (singleImgBoxWidth * float64(personCount))
					}
				}

				personStartY := calcY + heightOfHeader + verticalBorderTopBottomPicture + (heightOfPersonsInklNameAndSpace * float64(personLines-1))
				// fmt.Println(fmt.Sprintf("%f", personStartY))

				/* FACE */
				personImgPath := path.Join(inputDataFolder, "country_"+country, personMap.Picture.Path)
				im, err := gg.LoadImage(personImgPath)
				if err != nil {
					panic(err)
				}
				iw, ih := im.Bounds().Dx(), im.Bounds().Dy()
				newPersonImgWidth := (float64(iw) * heightOfPersons) / float64(ih)

				personCenter := personStartX + (singleImgBoxWidth / 2)

				pdf.Image(personImgPath,
					(personCenter - (newPersonImgWidth / 2)),
					personStartY,
					newPersonImgWidth,
					heightOfPersons,
					false, "", 0, "")

				/* NAME */
				personNamePath := path.Join(hashPicturesFolder, personMap.Name.Hash+".png")
				newPersonNameWidth := personMap.Name.Rendered.Width * personTipaScalingFactor
				newPersonNameHeight := personMap.Name.Rendered.Height * personTipaScalingFactor

				pdf.Image(personNamePath,
					(personCenter - (newPersonNameWidth / 2)),
					(personStartY + heightOfPersons + spacebetweenPictureAndImg),
					newPersonNameWidth,
					newPersonNameHeight,
					false, "", 0, "")

				count++
				personCount++
				if personCount > (imagesPerLine - 1) {
					personCount = 0
				}
			}
			pdf.ClipEnd()

		}
	}

	//
	// pt := gofpdf.PointType{X: 0, Y: 0}
	// st := gofpdf.SizeType{Wd: 6, Ht: 6}
	// tpl := gofpdf.CreateTpl(pt, st, "P", "mm", "", func(tpl *gofpdf.Tpl) {
	// 	tpl.Image(imgName, 10, 10, 30, 0, false, "", 0, "")
	// })
	//
	//
	//
	// pdf.AddPage()
	// pdf.SetFont("Arial", "B", 16)
	// pdf.Cell(40, 10, "Hello, world")
	// pdf.UseTemplate(tpl)
	// pdf.Image(imgName, 10, 10, 30, 0, false, "", 0, "")
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		fmt.Println(err)
	}
}

func sortPersons(p map[string]data.Person) []string {

	var PersonNames []string
	for personName := range p {
		PersonNames = append(PersonNames, personName)
	}

	var sortedKeys = make([]string, len(p))

	for i := range sortedKeys {
		lastFoundKey := "zzzzz"
		foundIdx := -1
		for idx, personName := range PersonNames {
			if personName < lastFoundKey {
				lastFoundKey = personName
				foundIdx = idx
			}
		}
		sortedKeys[i] = lastFoundKey
		PersonNames = remove(PersonNames, foundIdx)
	}
	return sortedKeys
}

func sortCountries(ct map[string]data.Country) []string {

	var countryNames []string
	for countryName := range ct {
		countryNames = append(countryNames, countryName)
	}

	var sortedKeys = make([]string, len(ct))

	for i := range sortedKeys {
		lastFoundKey := ""
		lastFoundPersonCount := 0
		foundIdx := -1
		for idx, countryName := range countryNames {
			if ct[countryName].PersonCount > lastFoundPersonCount {
				lastFoundKey = countryName
				lastFoundPersonCount = ct[countryName].PersonCount
				foundIdx = idx
			}
		}
		sortedKeys[i] = lastFoundKey
		countryNames = remove(countryNames, foundIdx)
	}
	return sortedKeys
}

func optimalSort(ct map[string]data.Country) []Page {
	sortedkeys := sortCountries(ct)
	var optimizedKeys = []Page{}

	for {
		usedFullLines := 0
		p := Page{}
		iterationKeys := sortedkeys
		for _, countryName := range iterationKeys {
			lines := calcLines(ct[countryName].PersonCount, imagesPerLine)
			if (usedFullLines + lines) <= linesPerPage {
				usedFullLines = (usedFullLines + lines)
				p.Countries = append(p.Countries, countryName)
				sortedkeys = remove(sortedkeys, indexOf(countryName, sortedkeys))
			}
		}
		optimizedKeys = append(optimizedKeys, p)
		if len(sortedkeys) <= 0 {
			break
		}
	}

	return optimizedKeys
}

func calcLines(personCount int, imagesPerLine int) int {
	return int(math.Ceil(float64(personCount) / float64(imagesPerLine)))
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

func getNameScalingFactor (countries map[string]data.Country) (cFactor float64, pFactor float64) {
	/* calculate the scaling factor for rendered transcriptions */

	cMax := 0.0
	pMax := 0.0
	var cHeight, pHeight float64

	for _, c := range countries {
		cHeight = c.Endonym.Rendered.Height
		if cHeight > cMax {
			cMax = cHeight
		}

		for _, p := range c.People {
			pHeight = p.Name.Rendered.Height
			if pHeight > pMax {
				pMax = pHeight
			}
		}
	}

	if cMax > countryNameFlagHeight {
		cFactor = countryNameFlagHeight / cMax
	} else {
		cFactor = 1.0
	}

	if pMax > heightOfName {
		pFactor = heightOfName / pMax
	} else {
		pFactor = 1.0
	}
	return
}
