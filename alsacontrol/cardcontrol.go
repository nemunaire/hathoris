package alsa

import (
	"bufio"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type CardControl struct {
	NumID     int64
	Interface string
	Name      string
	Type      string
	Access    string
	NValues   int64
	Min       int64
	Max       int64
	Step      int64
	DBScale   CardControldBScale
	Values    []string
	Items     []string
}

func (cc *CardControl) parseAmixerField(key, value string) (err error) {
	switch key {
	case "numid":
		cc.NumID, err = strconv.ParseInt(value, 10, 64)
	case "iface":
		cc.Interface = value
	case "name":
		cc.Name = strings.TrimPrefix(strings.TrimSuffix(value, "'"), "'")
	case "type":
		cc.Type = value
	case "access":
		cc.Access = value
	case "values":
		cc.NValues, err = strconv.ParseInt(value, 10, 64)
	case "min":
		cc.Min, err = strconv.ParseInt(value, 10, 64)
	case "max":
		cc.Max, err = strconv.ParseInt(value, 10, 64)
	case "step":
		cc.Step, err = strconv.ParseInt(value, 10, 64)
	}

	return
}

func (cc *CardControl) ToCardControlState() *CardControlState {
	ccs := &CardControlState{
		NumID: cc.NumID,
		Type:  cc.Type,
		Name:  cc.Name,
		RW:    strings.HasPrefix(cc.Access, "rw"),
		Items: cc.Items,
	}

	if cc.DBScale.Min != 0 || cc.DBScale.Step != 0 {
		ccs.DBScale = &cc.DBScale
	}

	// Convert values
	for _, v := range cc.Values {
		if cc.Type == "INTEGER" || cc.Type == "ENUMERATED" {
			if tmp, err := strconv.ParseFloat(v, 10); err == nil {
				ccs.Current = append(ccs.Current, tmp)
			}
		} else if cc.Type == "BOOLEAN" {
			if v == "on" {
				ccs.Current = append(ccs.Current, true)
			} else {
				ccs.Current = append(ccs.Current, false)
			}
		}
	}

	ccs.Min = cc.Min
	ccs.Max = cc.Max

	return ccs
}

type CardControldBScale struct {
	Min  float64
	Step float64
	Mute int64 `json:",omitempty"`
}

func (cc *CardControldBScale) parseAmixerField(key, value string) (err error) {
	switch key {
	case "min":
		cc.Min, err = strconv.ParseFloat(strings.TrimSuffix(value, "dB"), 10)
	case "step":
		cc.Step, err = strconv.ParseFloat(strings.TrimSuffix(value, "dB"), 10)
	case "mute":
		cc.Mute, err = strconv.ParseInt(value, 10, 64)
	}

	return
}

type CardControlState struct {
	NumID   int64
	Name    string
	Type    string
	RW      bool `json:"RW,omitempty"`
	Min     int64
	Max     int64
	DBScale *CardControldBScale `json:",omitempty"`
	Current []interface{}       `json:"values,omitempty"`
	Items   []string            `json:"items,omitempty"`
}

func ParseAmixerContent(cardId string) ([]*CardControl, error) {
	cardIdType := "-D"
	if _, err := strconv.Atoi(cardId); err == nil {
		cardIdType = "-c"
	}

	cmd := exec.Command("amixer", cardIdType, cardId, "-M", "contents")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var ret []*CardControl

	fscanner := bufio.NewScanner(stdout)
	for fscanner.Scan() {
		line := fscanner.Text()

		if strings.HasPrefix(line, "  ; Item #") {
			cc := ret[len(ret)-1]
			cc.Items = append(cc.Items, strings.TrimSuffix(line[strings.Index(line, "'")+1:], "'"))
		} else if strings.HasPrefix(line, "  :") {
			cc := ret[len(ret)-1]
			line = strings.TrimPrefix(line, "  : ")

			kv := strings.SplitN(line, "=", 2)
			if kv[0] == "values" {
				cc.Values = strings.Split(kv[1], ",")
			}
		} else if strings.HasPrefix(line, "  |") || strings.HasPrefix(line, "    |") {
			cc := ret[len(ret)-1]
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "| dBscale-") {
				line = strings.TrimPrefix(line, "| dBscale-")

				fields := strings.Split(line, ",")

				var scale CardControldBScale
				for _, field := range fields {
					kv := strings.SplitN(field, "=", 2)
					scale.parseAmixerField(kv[0], kv[1])
				}
				cc.DBScale = scale
			}
		} else {
			var cc *CardControl

			if strings.HasPrefix(line, "numid=") {
				cc = &CardControl{}
				ret = append(ret, cc)
			} else {
				cc = ret[len(ret)-1]
				line = strings.TrimPrefix(line, "  ; ")
			}

			fields := strings.Split(line, ",")

			for _, field := range fields {
				kv := strings.SplitN(field, "=", 2)
				cc.parseAmixerField(kv[0], kv[1])
			}
		}
	}

	err = cmd.Wait()

	// Sort mixers by NumID
	sort.Sort(ByNumID(ret))

	return ret, err
}

func (cc *CardControl) CsetAmixer(cardId string, values ...string) error {
	opts := []string{
		"-c",
		cardId,
		"-M",
		"cset",
		fmt.Sprintf("numid=%d", cc.NumID),
	}
	opts = append(opts, strings.Join(values, ","))
	cmd := exec.Command("amixer", opts...)

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

type ByNumID []*CardControl

func (a ByNumID) Len() int           { return len(a) }
func (a ByNumID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNumID) Less(i, j int) bool { return a[i].NumID < a[j].NumID }
