package data

import (
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/uconn-ling/openHouseMap/src/utils"

)

var countries map[string]Country

// return a list of countries that contain at least one person
func GetData(dataDir string, countryPrefix string) map[string]Country {

	countries = make(map[string]Country)

	files, err := ioutil.ReadDir("./" + dataDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		// fmt.Println("### " + file.Name() + "  " + strconv.Itoa((len(file.Name()))) + " " + strconv.Itoa(len(countryPrefix)))
		if !file.IsDir() || len(file.Name()) <= len(countryPrefix) || file.Name()[:len(countryPrefix)] != countryPrefix {
			continue
		}

		// remove 'country_' from directory name
		exonym := file.Name()[len(countryPrefix):]

		var country Country
		country.People = make(map[string]Person)
		// country.Endonym =
		// country.Flag =

		// read all files in this country folder
		countryFolderpath := path.Join("./" + dataDir + "/" + file.Name())
		files2, err2 := ioutil.ReadDir(countryFolderpath)
		if err2 != nil {
			log.Fatal(err2)
		}

		cname, _ := regexp.Compile(exonym + "-name\\.tex")
		cflag, _ := regexp.Compile(exonym + "-flag\\.[a-z]+")
		pname, _ := regexp.Compile("([a-z]+)\\.tex")           // decide which identifiers we allow
		ppic, _ := regexp.Compile("([a-z]+)\\.(png|jpg|jpeg)") // decide which endings we allow

		for _, file2 := range files2 {
			// log.Println(file2.Name())

			switch {
			case cname.MatchString(file2.Name()):
				// country.PathToName = file2.Name()
				country.Endonym.Path = file2.Name()
				// country.TipaString = strings.Replace(readTexFile(path.Join(countryFolderpath, file2.Name())), "\n", "", -1)
				country.Endonym.Value = strings.Replace(readTexFile(path.Join(countryFolderpath, file2.Name())), "\n", "", -1)
				// country.TipaHash = utils.HashString(country.TipaString)
				country.Endonym.Hash = utils.HashString(country.Endonym.Value)
			case cflag.MatchString(file2.Name()):
				// country.PathToFlag = file2.Name()
				country.Flag.Path = file2.Name()
			case pname.MatchString(file2.Name()):
				var key = pname.FindStringSubmatch(file2.Name())[1]
				p, ok := country.People[key]
				if !ok {
					p = Person{}
				}
				// p.PathToName = file2.Name()
				p.Name.Path = file2.Name()
				// p.TipaString = strings.Replace(readTexFile(path.Join(countryFolderpath, file2.Name())), "\n", "", -1)
				p.Name.Value = strings.Replace(readTexFile(path.Join(countryFolderpath, file2.Name())), "\n", "", -1)
				// p.TipaHash = utils.HashString(p.TipaString)
				p.Name.Hash = utils.HashString(p.Name.Value)
				country.People[key] = p
			case ppic.MatchString(file2.Name()):
				var key = ppic.FindStringSubmatch(file2.Name())[1]
				p, ok := country.People[key]
				if !ok {
					p = Person{}
				}
				// p.PathToFacePic = file2.Name()
				p.Picture.Path = file2.Name()
				country.People[key] = p
			default:
				log.Printf("Unexpected file found: " + file2.Name())
			}
		}

		// check that each person has a picture and a name
		for key, p := range country.People {
			if p.Picture.Path == "" {
				log.Fatal("Not enough data: person '" + key + "' lacks picture!")
			}
			if p.Name.Path == "" {
				log.Fatal("Not enough data: person '" + key + "' lacks name!")
			}
		}
		// check that each country has a flag and a name
		if country.Flag.Path == "" {
			log.Fatal("Not enough data: country '" + exonym + "' lacks flag!")
		}
		if country.Endonym.Path == "" {
			log.Fatal("Not enough data: country '" + exonym + "' lacks endonym!")
		}
		// check that this country has people in it
		if len(country.People) > 0 {
			country.PersonCount = len(country.People)
			countries[exonym] = country
		}
	}

	return countries
}

func readTexFile(filePath string) string {
	byteStr, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return string(byteStr)
}
