package main

type DimmingStruct struct {
	Brightness float32 `json:"brightness"`
}

type ColourTemperature struct {
	Mirek int `json:"mirek"`
}

type LightTemperature struct {
	ColourTemperature ColourTemperature    `json:"color_temperature"`
	Dimming_Struct    DimmingStruct `json:"dimming"`
}

type ColourVals struct {
	X_Val float32 `json:"x"`
	Y_Val float32 `json:"y"`
}

type ColourStruct struct {
	ColourData ColourVals `json:"xy"`
}

type LightColour struct {
	Colour_Struct  ColourStruct  `json:"color"`
	Dimming_Struct DimmingStruct `json:"dimming"`
}

type HueConfigData struct {
	Hue_IP_address       string `json:"hueIP"`
	Hue_app_key          string `json:"appKey"`
	Hue_grouped_light_id string `json:"groupedLight"`
	Port                 string `json:"port"`
}

type HueLightValue struct {
	LightDataType string  `json:"light"`
	Mirek_value   int     `json:"mirek"`
	X_Data        float32 `json:"x"`
	Y_Data        float32 `json:"y"`
	Brightness    float32 `json:"brightness"`
}

type HueAuth struct {
	DeviceType        string `json:"devicetype"`
	GenerateClientKey bool   `json:"generateclientkey"`
}