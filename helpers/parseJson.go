package helpers

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

type AllPolicies struct {
	ID           int                 `json:"id"`
	Name         string              `json:"name"`
	Class        string              `json:"class"`
	DisplayName  string              `json:"displayName"`
	Presentation []Presentation_json `json:"presentation,omitempty"`
	//View 		 []string 			 `json:"view,omitempty"`
	ExplainText string   `json:"explainText"`
	Category    string   `json:"category"`
	SupportedOn string   `json:"supportedOn"`
	Values      []Values `json:"values"`
	HasManual   bool     `json:"hasManual"`
}

type AllPoliciesBson struct {
	//ID           bson.ObjectId `bson:"_id"`
	ID           string   `bson:"id"`
	IDtmp        int      `bson:"IDtmp"`
	Name         string   `bson:"name"`
	Class        string   `bson:"class"`
	DisplayName  string   `bson:"displayName"`
	ExplainText  string   `bson:"explainText"`
	Category     string   `bson:"category"`
	SupportedOn  string   `bson:"supportedOn"`
	Values       []Values `bson:"values"`
	GpName       string   `bson:"gpname,omitempty"`
	GpType       string   `bson:"gptype,omitempty"` // usr, def, sub
	Dependencies string   `bson:"dependencies,omitempty"`
	HasManual    bool     `bson:"manual,omitempty"`
}

/*type Presentation_json struct {
	Chardata string `json:",chardata,omitempty"`
	ID       string `json:"id,omitempty"`
	CheckBox []struct {
		Text           string `json:",chardata,omitempty"`
		RefId          string `json:"refId,omitempty"`
		DefaultChecked string `json:"defaultChecked,omitempty"`
	} `json:"checkBox,omitempty"`
	ComboBox []struct {
		Text       string   `json:",chardata,omitempty"`
		RefId      string   `json:"refId,omitempty"`
		NoSort     string   `json:"noSort,omitempty"`
		Label      string   `json:"label,omitempty"`
		Suggestion []string `json:"suggestion,omitempty"`
	} `json:"comboBox,omitempty"`
	DropdownList []struct {
		Text        string `json:",chardata,omitempty"`
		RefId       string `json:"refId,omitempty"`
		DefaultItem string `json:"defaultItem,omitempty"`
		NoSort      string `json:"noSort,omitempty"`
	} `json:"dropdownList,omitempty"`
	Text    []string `json:"text,omitempty"`
	ListBox []struct {
		Text  string `json:",chardata,omitempty"`
		RefId string `json:"refId,omitempty"`
	} `json:"listBox,omitempty"`
	DecimalTextBox []struct {
		Text         string `json:",chardata,omitempty"`
		RefId        string `json:"refId,omitempty"`
		DefaultValue string `json:"defaultValue,omitempty"`
		SpinStep     string `json:"spinStep,omitempty"`
	} `json:"decimalTextBox,omitempty"`
	LongDecimalTextBox []struct {
		Text         string `json:",chardata,omitempty"`
		RefId        string `json:"refId,omitempty"`
		DefaultValue string `json:"defaultValue,omitempty"`
		SpinStep     string `json:"spinStep,omitempty"`
	} `json:"longDecimalTextBox,omitempty"`
	TextBox []struct {
		Text         string `json:",chardata,omitempty"`
		RefId        string `json:"refId,omitempty"`
		Label        string `json:"label,omitempty"`
		DefaultValue string `json:"defaultValue,omitempty"`
	} `json:"textBox,omitempty"`
	MultiTextBox []struct {
		Text  string `json:",chardata,omitempty"`
		RefId string `json:"refId,omitempty"`
	} `json:"multiTextBox,omitempty"`
}*/

type Values struct {
	Type          string `json:"type,omitempty" bson:"type,omitempty"`
	ValueName     string `json:"valueName,omitempty" bson:"valueName,omitempty"`
	DisplayName   string `json:"displayName,omitempty" bson:"displayName,omitempty"`
	Key           string `json:"key,omitempty" bson:"key,omitempty"`
	Required      string `json:"required,omitempty" bson:"required,omitempty"`
	MaxValue      string `json:"maxValue,omitempty" bson:"maxValue,omitempty"`
	MinValue      string `json:"minValue,omitempty" bson:"minValue,omitempty"`
	Value         string `json:"value,omitempty" bson:"value,omitempty"`
	DisabledValue string `json:"disabledValue,omitempty" bson:"disabledValue,omitempty"`
	EnabledValue  string `json:"enabledValue,omitempty" bson:"enabledValue,omitempty"`
	TrueValue     string `json:"trueValue,omitempty" bson:"trueValue,omitempty"`
	FalseValue    string `json:"falseValue,omitempty" bson:"falseValue,omitempty"`
	ValuePrefix   string `json:"valuePrefix,omitempty" bson:"valuePrefix,omitempty"`
	SelectedValue string `json:"selectedvalue,omitempty" bson:"selectedvalue,omitempty"`
	Showing       string `json:"showing,omitempty" bson:"showing,omitempty"`
	Dublicate     bool   `json:"dublicate,omitempty" bson:"dublicate,omitempty"`
	Notes         string `json:"notes,omitempty" bson:"notes,omitempty"`
	Manual        bool   `json:"manual,omitempty" bson:"manual,omitempty"`
}

func unique(vals []Values) []Values {
	keys := make(map[string]bool)
	list := []Values{}
	for _, entry := range vals {
		if _, value := keys[entry.ValueName]; !value {
			keys[entry.ValueName] = true
			list = append(list, entry)
		}
	}
	return list
}

func AllgpToBson(c *mgo.Collection, result []AllPolicies) {
	r := AllPoliciesBson{}
	for _, pol := range joinValues(result) {
		r.ID = bson.NewObjectId().Hex()
		r.IDtmp = pol.ID
		r.Category = pol.Category
		r.Class = pol.Class
		r.DisplayName = pol.DisplayName
		r.SupportedOn = pol.SupportedOn
		r.ExplainText = pol.ExplainText
		r.Values = pol.Values
		r.Name = pol.Name

		err := c.Insert(&r)
		if err != nil {
			log.Println(err)
		}
	}
}

func joinValues(jsonGP []AllPolicies) []AllPolicies {
	var vals []Values
	var resultGP []AllPolicies
	for _, pol := range jsonGP {
		vals = recVals(pol.Values)
		pol.Values = vals
		resultGP = append(resultGP, pol)
	}
	return resultGP
}

func recVals(vals []Values) []Values {
	var newVals []Values
	var v = make(map[string]string)
	var d = make(map[string]string)
	for _, val := range vals {
		if _, ok := v[val.ValueName]; ok {
			v[val.ValueName] += "|" + val.Value
		} else {
			v[val.ValueName] = val.Value
		}
		if _, ok := d[val.ValueName]; ok {
			d[val.ValueName] += "|" + val.DisplayName
		} else {
			d[val.ValueName] = val.DisplayName
		}
	}
	for _, val := range vals {
		if strings.Contains(v[val.ValueName], "|") {
			val.Value = v[val.ValueName]
			val.Dublicate = true
		}
		if strings.Contains(d[val.ValueName], "|") {
			val.DisplayName = d[val.ValueName]
		}
		newVals = append(newVals, val)
	}

	return unique(newVals)
}
