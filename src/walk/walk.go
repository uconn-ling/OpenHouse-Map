package walk

import (
  "io/ioutil"
  "log"
  "regexp"
)

type Country struct {
  pathToName string
  pathToFlag string
  people map[string]Person
}

type Person struct {
  pathToName string
  pathToPic string
}

var countries map[string]Country

// return a list of countries that contain at least one person
func GetData (dataDir string, countryPrefix string) map[string]Country {

  countries = make(map[string]Country)

  files, err := ioutil.ReadDir("./" + dataDir)
  if err != nil {
    log.Fatal(err)
  }

  for _, file := range files {

    if !file.IsDir() || file.Name()[:len(countryPrefix)] != countryPrefix {
      continue
    }

    // remove 'country_' from directory name
    exonym := file.Name()[len(countryPrefix):]

    var country Country
    country.people = make(map[string]Person)

    // read all files in this country folder
    files2, err2 := ioutil.ReadDir("./" + dataDir + "/" + file.Name())
    if err2 != nil {
      log.Fatal(err2)
    }

    cname, _ := regexp.Compile(exonym + "-name\\.tex")
    cflag, _ := regexp.Compile(exonym + "-flag\\.[a-z]+")
    pname, _ := regexp.Compile("([a-z]+)\\.tex") // decide which identifiers we allow
    ppic,  _ := regexp.Compile("([a-z]+)\\.(png|jpg|jpeg)") // decide which endings we allow

    for _, file2 := range files2 {
      log.Println(file2.Name())

      switch {
        case cname.MatchString(file2.Name()):
          country.pathToName = file2.Name()
        case cflag.MatchString(file2.Name()):
          country.pathToFlag = file2.Name()
        case pname.MatchString(file2.Name()):
          var key = pname.FindStringSubmatch(file2.Name())[1]
          p, ok := country.people[key]
          if !ok {
            p = Person{}
          }
          p.pathToName = file2.Name()
          country.people[key] = p
        case ppic.MatchString(file2.Name()):
          var key = ppic.FindStringSubmatch(file2.Name())[1]
          p, ok := country.people[key]
          if !ok {
            p = Person{}
          }
          p.pathToPic = file2.Name()
          country.people[key] = p
        default:
          log.Printf("Unexpected file found: " + file2.Name() )
      }

    }

    // check that each person has a picture and a name
    for key, p := range country.people {
      if p.pathToPic == "" {
        log.Fatal("Not enough data: person '" + key + "' lacks picture!")
      }
      if p.pathToName == "" {
        log.Fatal("Not enough data: person '" + key + "' lacks name!")
      }
    }
    // check that each country has a flag and a name
    if country.pathToFlag == "" {
      log.Fatal("Not enough data: country '" + exonym + "' lacks flag!")
    }
    if country.pathToName == "" {
      log.Fatal("Not enough data: country '" + exonym + "' lacks endonym!")
    }
    // check that this country has people in it
    if len(country.people) == 0 {
      countries[exonym] = country
    }
  }
  return countries
}
