package broadcast

const BroadcastingProtocolVersion byte = 4

type OutboundMessageType byte

const (
	OutboundRegisterCommandApplication   OutboundMessageType = 1
	OutboundUnregisterCommandApplication OutboundMessageType = 9
	OutboundRequestEntryList             OutboundMessageType = 10
	OutboundRequestTrackData             OutboundMessageType = 11
	OutboundChangeHUDPage                OutboundMessageType = 49
	OutboundChangeFocus                  OutboundMessageType = 50
	OutboundInstantReplayRequest         OutboundMessageType = 51
	OutboundSaveManualReplayHighlight    OutboundMessageType = 60
)

type InboundMessageType byte

const (
	InboundRegistrationResult InboundMessageType = 1
	InboundRealtimeUpdate     InboundMessageType = 2
	InboundRealtimeCarUpdate  InboundMessageType = 3
	InboundEntryList          InboundMessageType = 4
	InboundTrackData          InboundMessageType = 5
	InboundEntryListCar       InboundMessageType = 6
	InboundBroadcastingEvent  InboundMessageType = 7
)

var CarModels = map[byte]string{
	// GT3 - 2018
	0:  "Porsche 911 GT3 R 2018",
	1:  "Mercedes-AMG GT3 2015",
	2:  "Ferrari 488 GT3 2018",
	3:  "Audi R8 LMS 2015",
	4:  "Lamborghini Huracán GT3 2015",
	5:  "McLaren 650S GT3 2015",
	6:  "Nissan GT-R Nismo GT3 2018",
	7:  "BMW M6 GT3 2017",
	8:  "Bentley Continental GT3 2018",
	9:  "Porsche 911 II GT3 Cup 2017",
	10: "Nissan GT-R Nismo GT3 2015",
	11: "Bentley Continental GT3 2015",
	12: "Aston Martin Vantage V12 GT3 2013",
	13: "Lamborghini Gallardo G3 Reiter 2017",
	14: "Emil Frey Jaguar G3 2012",
	15: "Lexus RCF GT3 2016",
	16: "Lamborghini Huracán GT3 Evo 2019",
	17: "Honda NSX GT3 2017",
	18: "Lamborghini Huracán ST 2015",

	// GT3 - 2019
	19: "Audi R8 LMS Evo 2019",
	20: "Aston Martin V8 Vantage GT3 2019",
	21: "Honda NSX GT3 Evo 2019",
	22: "McLaren 720S GT3 2019",
	23: "Porsche 911 II GT3 R 2019",

	// GT3 - 2020
	24: "Ferrari 488 GT3 Evo 2020",
	25: "Mercedes-AMG GT3 2020",

	// GTC (Challengers Pack)
	26: "Ferrari 488 Challenge Evo 2020",
	27: "BMW M2 Cup 2020",
	28: "Porsche 992 GT3 Cup 2021",
	29: "Lamborghini Huracán ST Evo2 2021",

	// GT3 - 2021
	30: "BMW M4 GT3 2021",

	// GT3 - 2022
	31: "Audi R8 LMS Evo II 2022",

	// GT3 - 2023
	32: "Ferrari 296 GT3 2023",
	33: "Lamborghini Huracán GT3 Evo2 2023",
	34: "Porsche 992 GT3 R 2023",
	35: "McLaren 720S GT3 Evo 2023",

	// GT3 - 2024
	36: "Ford Mustang GT3 2024",

	// GT4
	50: "Alpine A110 GT4 2018",
	51: "Aston Martin Vantage AMR GT4 2018",
	52: "Audi R8 LMS GT4 2016",
	53: "BMW M4 GT4 2018",
	55: "Chevrolet Camaro GT4 R 2017",
	56: "Ginetta G55 GT4 2012",
	57: "KTM X-BOW GT4 2016",
	58: "Maserati Gran Turismo MC GT4 2016",
	59: "McLaren 570s GT4 2016",
	60: "Mercedes AMG GT4 2016",
	61: "Porsche 718 Cayman GT4 MR 2019",

	// GT2
	80: "Audi R8 LMS GT2 2021",
	82: "KTM X-BOW GT2 2021",
	83: "Maserati GT2 2023",
	84: "Mercedes-AMG GT2 2023",
	85: "Porsche 991 II GT2 RS CS Evo 2023",
	86: "Porsche 935 2019",
}

func GetCarModelName(modelID byte) string {
	if name, exists := CarModels[modelID]; exists {
		return name
	}
	return "Unknown Car"
}
