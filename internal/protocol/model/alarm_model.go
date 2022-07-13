package model

type AlarmModel struct {
	Identifier        uint8  `json:"identifier"`
	Name              string `json:"name"`
	MinMainCpuVersion string `json:"minMainCpuVersion"`
	MinGrpsCpuVersion string `json:"minGprsCpuVersion"`
}

const (
	AMT_4010_SMART = 0x41
	AMT_2018_EG    = 0x1E
	ANM_24_NET     = 0x24
)

var AlarmModels = make(map[int]AlarmModel)

func init() {
	AlarmModels[AMT_4010_SMART] = AlarmModel{
		Identifier:        0x41,
		Name:              "AMT 4010 SMART",
		MinMainCpuVersion: "3.0",
		MinGrpsCpuVersion: "2.0",
	}
	AlarmModels[AMT_2018_EG] = AlarmModel{
		Identifier:        0x1E,
		Name:              "AMT 2018 EG",
		MinMainCpuVersion: "7.0",
		MinGrpsCpuVersion: "4.0",
	}
	AlarmModels[AMT_2018_EG] = AlarmModel{
		Identifier:        0x24,
		Name:              "ANM 24 NET",
		MinMainCpuVersion: "3.0",
		MinGrpsCpuVersion: "",
	}
}

func GetAlarmModel(model int) AlarmModel {
	if model, ok := AlarmModels[model]; ok {
		return model
	}
	return AlarmModel{}
}
