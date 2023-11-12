package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
)

func declareVolumeRoutes(cfg *config.Config, router *gin.RouterGroup) {
	router.GET("/mixer", func(c *gin.Context) {
		cnt, err := parseAmixerContent()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, cnt)
	})

	mixerRoutes := router.Group("/mixer/:mixer")
	mixerRoutes.Use(MixerHandler)

	mixerRoutes.GET("", func(c *gin.Context) {
		cc := c.MustGet("mixer").(*CardControl)

		c.JSON(http.StatusOK, cc.ToCardControlState())
	})
	mixerRoutes.POST("/values", func(c *gin.Context) {
		cc := c.MustGet("mixer").(*CardControl)

		var valuesINT []interface{}

		err := c.BindJSON(&valuesINT)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errmsg": fmt.Sprintf("Unable to parse values: %s", err.Error())})
			return
		}

		var values []string
		for _, v := range valuesINT {
			if t, ok := v.(int64); ok {
				values = append(values, strconv.FormatInt(t, 10))
			} else if t, ok := v.(float64); ok {
				values = append(values, fmt.Sprintf("%f", t))
			} else if t, ok := v.(bool); ok {
				if t {
					values = append(values, "on")
				} else {
					values = append(values, "off")
				}
			}
		}

		err = cc.CsetAmixer(values...)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": fmt.Sprintf("Unable to set values: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, values)
	})
}

func MixerHandler(c *gin.Context) {
	mixers, err := parseAmixerContent()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": err.Error()})
		return
	}

	mixer := c.Param("mixer")

	for _, m := range mixers {
		if strconv.FormatInt(m.NumID, 10) == mixer || m.Name == mixer {
			c.Set("mixer", m)
			c.Next()
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errmsg": "Mixer not found"})
}

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
		cc.Name = value
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
		Items: cc.Items,
	}

	// Convert values
	for _, v := range cc.Values {
		if cc.Type == "INTEGER" {
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

	if cc.DBScale.Min != 0 {
		ccs.Min = cc.DBScale.Min
		ccs.Unit = "dB"
	} else if cc.Min != 0 {
		ccs.Min = float64(cc.Min)
	}

	if cc.DBScale.Step != 0 {
		ccs.Step = cc.DBScale.Step
		ccs.Unit = "dB"
	} else if cc.Step != 0 {
		ccs.Step = float64(cc.Step)
	} else {
		ccs.Step = 1.0
	}

	if cc.Max != 0 {
		ccs.Max = ccs.Min + ccs.Step*float64(cc.Max-cc.Min)
	}

	return ccs
}

type CardControldBScale struct {
	Min  float64
	Step float64
	Mute int64
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
	Min     float64
	Max     float64
	Step    float64
	Unit    string        `json:"unit,omitempty"`
	Current []interface{} `json:"values,omitempty"`
	Items   []string      `json:"items,omitempty"`
}

func parseAmixerContent() ([]*CardControl, error) {
	cmd := exec.Command("amixer", "-c1", "contents")

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
	return ret, err
}

func (cc *CardControl) CsetAmixer(values ...string) error {
	opts := []string{
		"-c1",
		"cset",
		fmt.Sprintf("numid=%d", cc.NumID),
	}
	opts = append(opts, values...)
	cmd := exec.Command("amixer", opts...)

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
