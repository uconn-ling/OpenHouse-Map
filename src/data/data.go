package data

type Country struct {
	Endonym 		Transcription
	// PathToName  string
	// TipaString  string
	// TipaHash    string
	// TipaWidth 	float64
	// TipaHeight 	float64
	Flag 				Picture
	// PathToFlag  string
	// FlagWidth 	float64
	// FlagHeight 	float64
	People      map[string]Person
	PersonCount int
}

type Person struct {
	Name					Transcription
	// PathToName    string
	// TipaString    string
	// TipaHash      string
	// TipaWidth 		float64
	// TipaHeight 		float64
	Picture 			Picture
	// PathToFacePic string
	// PicWidth 			float64
	// PicHeight 		float64
}

type Picture struct {
	Path 					string
	Width 				float64
	Height 				float64
}

type Transcription struct {
	Path		 			string
	Value					string
	Hash					string
	Rendered 			Picture
}
