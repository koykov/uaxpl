package main

import (
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/koykov/entry"
	"gopkg.in/yaml.v2"
)

type deviceTuple map[string]deviceBrand

type deviceBrand struct {
	Name   string        `yaml:"-"`
	Regex  string        `yaml:"regex"`
	Device string        `yaml:"device"`
	Model  string        `yaml:"model,omitempty"`
	Models []deviceModel `yaml:"models,omitempty"`
}

type deviceModel struct {
	Regex  string `yaml:"regex"`
	Model  string `yaml:"model"`
	Device string `yaml:"device"`
}

type deviceModule struct{}

var (
	// Regexp replaces
	staticReplDevice = map[string]string{
		`TCL`: `TCL[^\s]`,
	}
	reReplDevice = map[string]string{
		`MOT|DROID ?(?:Build|[a-z0-9]+)|portalmmm/2.0 (?:E378i|L6|L7|v3)|XOOM [^;/]*Build|XT1941-2|XT1924-9|XT1925-10|XT1965-6|XT1970-5|XT1799-2|XT1021|XT2171-3|XT2071-4|XT2175-2|XT2125-4|XT2143-1|XT2153-1|XT2201-2|XT2137-2|XT1710-08|XT180[3-5]|XT194[23]-1|XT1929-15|(?:XT|MZ|MB|ME)[0-9]{3,4}[a-z]?(?:\(Defy\)|-0[1-5])?(?:[;]? Build|\))`: `MOT[^T][^O]|[^A][^N]DROID ?(?:Build|[a-z0-9]+)|portalmmm\/2.0 (?:E378i|L6|L7|v3)|XOOM [^;\/]*Build|XT1941-2|XT1924-9|XT1925-10|XT1965-6|XT1970-5|XT1799-2|XT1021|XT2171-3|XT2071-4|XT2175-2|XT2125-4|XT2143-1|XT2153-1|XT2201-2|XT2137-2|XT1710-08|XT180[3-5]|XT194[23]-1|XT1929-15|(?:XT|MZ|MB|ME)[0-9]{3,4}[a-z]?(?:\(Defy\)|-0[1-5])?(?:[;]? Build|\))`,
		`CUBOT|QUEST(?: LITE)?|(?:NOTE 7|NOTE 20(?: PRO)?|RAINBOW 2|KING_KONG_3|KINGKONG [57]|KINGKONG 5 Pro|ZORRO 001|KINGKONG[_ ]MINI2?|KING KONG CS|R15_PRO|S550 Pro| CHEETAH| [RX]19|X19 S|MAX 3)(?:[);/ ]|$)`:                                                                                                                                `CUBOT|QUEST(?: LITE)?|(?:[^REDMITGA ]+NOTE 7|NOTE 20(?: PRO)?|RAINBOW 2|KING_KONG_3|KINGKONG [57]|KINGKONG 5 Pro|ZORRO 001|KINGKONG[_ ]MINI2?|KING KONG CS|R15_PRO|S550 Pro| CHEETAH| [RX]19|X19 S|MAX 3)(?:[);/ ]|$)`,
		`Samsung|Maple_2011`: `Samsung.*Smart\s*TV|Tizen.*TV|Maple_2011`,
		`LG|Web0S; Linux/SmartTV|portalmmm/2\.0 (?:KE|KG|KP|L3)|(?:VX[0-9]+|L-0[12]D|L-07C|P713|(LM-[A-Z][0-9]{3}[A-Z]{0,3}(?:\(FGN?\)|\(G\)|\.FN|\.F(?:GN?)?)?)|LM-(X625N1|G820QM[0127])|G Watch)(?:[);/ ]|$)|NetCast|RS98[78]|RS500|VS(?:50[01]|835|876|880|98[05-8]|99[056]|(410|415|425|500|880)PP)|VS820|VS870 4G|Nexsus 5|LML(41[34]D|211B|713D|212V)L|VK810 4G|VK815|X5-LG|L-05[DE]|L-03K|L-01[LJK]|L-52A|L-[45]1A|901LG|A001LG|80[12]LG|402LG|LGT32|L-02K|L-01F|LM-V40|DM-01K|DM-02H|LMX130IM|VS810PP|VS425|VM670`:                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           `[^\d]LG|Web0S; Linux/SmartTV|portalmmm/2\.0 (?:KE|KG|KP|L3)|(?:VX[0-9]+|L-0[12]D|L-07C|P713|(LM-[A-Z][0-9]{3}[A-Z]{0,3}(?:\(FGN?\)|\(G\)|\.FN|\.F(?:GN?)?)?)|LM-(X625N1|G820QM[0127])|G Watch)(?:[);/ ]|$)|NetCast|RS98[78]|RS500|VS(?:50[01]|835|876|880|98[05-8]|99[056]|(410|415|425|500|880)PP)|VS820|VS870 4G|Nexsus 5|LML(41[34]D|211B|713D|212V)L|VK810 4G|VK815|X5-LG|L-05[DE]|L-03K|L-01[LJK]|L-52A|L-[45]1A|901LG|A001LG|80[12]LG|402LG|LGT32|L-02K|L-01F|LM-V40|DM-01K|DM-02H|LMX130IM|VS810PP|VS425|VM670`,
		`(?:WIKO[_ -])|(?:WIKO-)?CINK_[a-z0-9_]+|(?:WIKO-)?HIGHWAY_[a-z0-9_]+|(?:WIKO-)?(?:Cink([^/;]*)|Highway([^/;]*)|Iggy|Stairway|Rainbow ?(Jam|Lite|Up)?|Darkside|Darkmoon|Darkfull|Darknight|Freddy|FEVER|Jerry[2]?|Sublim|Ozzy|JIMMY|Barry|Birdy|Bloom|Getaway|Goa|Kite|Robby|Lenny[2-9]?|Slide|Sunset[2-9]?|Sunny[2-9]?|Tommy[2-9]?|PULP(?: Fab)?|Wax|HARRY|Ridge(?: Fab)?|U FEEL(?: Prime| Lite)?|U PULSE(?: LITE)?|View2 (?:Go|Plus)|View XL|View Prime)(?: Plus)?(?: 4G)?(?: Build|$)?|W_(?:C8[06]0|K[346]00|P130|P200|C200SN)|WC300|W-(?:V7[23]0|V745|V755|V8[035]0|P[36]11|K130|K211|K3[68]0|K420|K510S?|K560|K610|K630)-(?:BYT|EEA|FRA|OPE|SUN|SWI|TIM|TMV|TVM|DU|T[HV]|ID)|WIM Lite|W-V851|W-V6[08]0|W-V750BN|W-K200|W-U300|U307AS|CINK PEAX 2|HIGHWAY PURE|C210AE|v3913|U614AS|V3991AN`:                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              `(?:WIKO[_ -])|(?:WIKO-)?CINK_[a-z0-9_]+|(?:WIKO-)?HIGHWAY_[a-z0-9_]+|(?:WIKO-)?(?:Cink([^/;]*)|Highway([^/;]*)|Iggy|Stairway|Rainbow ?(Jam|Lite|Up)?|Darkside|Darkmoon|Darkfull|Darknight|Freddy|FEVER|Jerry[2]?|Sublim|Ozzy|JIMMY|Barry|Birdy|Bloom|Getaway|Goa[^nalplus]|Kite|Robby|Lenny[2-9]?|Slide|Sunset[2-9]?|Sunny[2-9]?|Tommy[2-9]?|PULP(?: Fab)?|Wax|HARRY|Ridge[^\splus](?: Fab)?|U FEEL(?: Prime| Lite)?|U PULSE(?: LITE)?|View2 (?:Go|Plus)|View XL|View Prime)(?: Plus)?(?: 4G)?(?: Build|$)?|W_(?:C8[06]0|K[346]00|P130|P200|C200SN)|WC300|W-(?:V7[23]0|V745|V755|V8[035]0|P[36]11|K130|K211|K3[68]0|K420|K510S?|K560|K610|K630)-(?:BYT|EEA|FRA|OPE|SUN|SWI|TIM|TMV|TVM|DU|T[HV]|ID)|WIM Lite|W-V851|W-V6[08]0|W-V750BN|W-K200|W-U300|U307AS|CINK PEAX 2|HIGHWAY PURE|C210AE|v3913|U614AS|V3991AN`,
		`Alcatel|Alc[a-z0-9]+|One[ _]?Touch|idol3|TIMXL|(?:(?:4003[AJ]|4009[ADEFIKMSX]|4013[DEJKMX]|4014[ADEKMX]|4015[ADNTX]|4016[ADX]|4017[ADEFSX]|4018[ADEFMX]|4024[DEX]|4027[ADNX]|4028[AEJS]|4032[ADEX]|4034[ADEFGLTX]|4035[ADXY]|4045[ADELX]|4047[ADFGNX]|4049[DEGMX]|4060[SW]|4087U|4114E|5001[ADJTU]|5002[ABDEFJHLMRSX]|5003[ADGU]|5006[ADG]|5007[AGUO]|5008[ADTUY]|5009[AD]|5010[DEGSUX]|5011A|5012[DFG]|5015[ADEX]|5016[AJX]|5017[ABDEOX]|5019D|5022[DEX]|5023[EF]|5024[ADFIJ]|5025[DEG]|5026[ADJ]|5027B|5028[ADY]|5029[ADEFY]|5030[ADEFJMU]|5032W|5033(?:[YX]R|FP|[ADEFGJMOQSTXY])|5034D|5036D|5038[ADEX]|5039[DY]|5041[CD]|5042[ADEFGTWX]|5044[ADGIKOPSTY]|5045[ADFGIJTXY]|5046[ADGIJSTUY]|5047[DIUY]|5048[AIUY]|5049[EGSWZ]|5050[ASXY]|5051[ADEJMTWX]|5052[ADY]|5053[ADKY]|5054[ADNSTWX]|5056[ADEGIJMNTUWX]|5057M|5058[AIY]|5059[ADIJSTXYZ]|5060[ADJ]|5061[AKU]|5065[ADNWX]|5070D|5080[ADFQUX]|5085[ABCDGHIJNOQY]|5086[ADY]|5090[AIY]|5095[IKY]|5098[OS]|5099[ADIUY]|5116J|5145A|6016[ADEX]|6025[DH]|6036[AXY]|6037[BKY]|6039[AHJKSY]|6042D|6043[AD]|6044D|6045[BFIKOXY]|6050[AFY]|6055[ABDHIKPUYZ]|6056H|6058[ADX]|6060[SXC]|6062W|6070K|7040[ADEFKRT]|7041[DX]|7042A|7043[AEKY]|7044[AX]|7045Y|7048[ASWX]|7053D|7055A|7070X|7071[ADX]|8030Y|8050[DEGX]|8063|8088[MQX]|8094M|9001[DIX]|9002X|9003[AX]|9005X|9007[ATX]|9008[ADIJNTUX]|9009[AG]|9010X|9013[AX]|9020A|9022X|9024O|9026X|9027[FQTWX]|9029Z|9032[WXZ]|9203A|A(?:46[46]BG|57[04]BL|576CC|621BL|50[13]DL|57[17]VL)|I213|I216[AX]|80(?:82|6[78]|5[12]))(?:_(?:EEA|RU|TR))?|P3[36]0X|AM-H200|4063F)(?:[);/ ]|$)`: `Alcatel|Alc[^o][^r][^\s][a-z0-9]+|One[ _]?Touch|idol3|TIMXL|(?:(?:4003[AJ]|4009[ADEFIKMSX]|4013[DEJKMX]|4014[ADEKMX]|4015[ADNTX]|4016[ADX]|4017[ADEFSX]|4018[ADEFMX]|4024[DEX]|4027[ADNX]|4028[AEJS]|4032[ADEX]|4034[ADEFGLTX]|4035[ADXY]|4045[ADELX]|4047[ADFGNX]|4049[DEGMX]|4060[SW]|4087U|4114E|5001[ADJTU]|5002[ABDEFJHLMRSX]|5003[ADGU]|5006[ADG]|5007[AGUO]|5008[ADTUY]|5009[AD]|5010[DEGSUX]|5011A|5012[DFG]|5015[ADEX]|5016[AJX]|5017[ABDEOX]|5019D|5022[DEX]|5023[EF]|5024[ADFIJ]|5025[DEG]|5026[ADJ]|5027B|5028[ADY]|5029[ADEFY]|5030[ADEFJMU]|5032W|5033(?:[YX]R|FP|[ADEFGJMOQSTXY])|5034D|5036D|5038[ADEX]|5039[DY]|5041[CD]|5042[ADEFGTWX]|5044[ADGIKOPSTY]|5045[ADFGIJTXY]|5046[ADGIJSTUY]|5047[DIUY]|5048[AIUY]|5049[EGSWZ]|5050[ASXY]|5051[ADEJMTWX]|5052[ADY]|5053[ADKY]|5054[ADNSTWX]|5056[ADEGIJMNTUWX]|5057M|5058[AIY]|5059[ADIJSTXYZ]|5060[ADJ]|5061[AKU]|5065[ADNWX]|5070D|5080[ADFQUX]|5085[ABCDGHIJNOQY]|5086[ADY]|5090[AIY]|5095[IKY]|5098[OS]|5099[ADIUY]|5116J|5145A|6016[ADEX]|6025[DH]|6036[AXY]|6037[BKY]|6039[AHJKSY]|6042D|6043[AD]|6044D|6045[BFIKOXY]|6050[AFY]|6055[ABDHIKPUYZ]|6056H|6058[ADX]|6060[SXC]|6062W|6070K|7040[ADEFKRT]|7041[DX]|7042A|7043[AEKY]|7044[AX]|7045Y|7048[ASWX]|7053D|7055A|7070X|7071[ADX]|8030Y|8050[DEGX]|8063|8088[MQX]|8094M|9001[DIX]|9002X|9003[AX]|9005X|9007[ATX]|9008[ADIJNTUX]|9009[AG]|9010X|9013[AX]|9020A|9022X|9024O|9026X|9027[FQTWX]|9029Z|9032[WXZ]|9203A|A(?:46[46]BG|57[04]BL|576CC|621BL|50[13]DL|57[17]VL)|I213|I216[AX]|80(?:82|6[78]|5[12]))(?:_(?:EEA|RU|TR))?|P3[36]0X|AM-H200|4063F)(?:[);/ ]|$)`,
	}
)

func (m deviceModule) Validate(input, _ string) error {
	if len(input) == 0 {
		return fmt.Errorf("param -input is required")
	}
	return nil
}

func (m deviceModule) Compile(w moduleWriter, input, target string) (err error) {
	if len(target) == 0 {
		target = "device_repo.go"
	}

	hd, _ := os.UserHomeDir()
	input = strings.ReplaceAll(input, "~", hd)

	files, err := filepath.Glob(input + "/*.yml")
	if err != nil {
		return err
	}

	var (
		bufST []deviceBrand
		bufDR []string
		bufDM []string
		bufRE []string
		buf   buf
	)

	_, _ = w.WriteString("import (\n\"regexp\"\n)\n\n")
	_, _ = w.WriteString("var (\n")

	_, _ = w.WriteString("__dr_idx = [" + strconv.Itoa(len(files)) + "][]deviceTuple{\n")

	for i := 0; i < len(files); i++ {
		bufST = bufST[:0]
		bufDR = bufDR[:0]

		var body []byte
		if body, err = ioutil.ReadFile(files[i]); err != nil {
			return
		}
		if len(body) == 0 {
			err = errors.New("nothing to parse")
			return
		}

		tuples := make(deviceTuple)
		if err = yaml.Unmarshal(body, &tuples); err != nil {
			return
		}

		for name, brand := range tuples {
			brand.Name = name
			bufST = append(bufST, brand)
		}
		sort.Slice(bufST, func(i, j int) bool {
			return bufST[i].Name < bufST[j].Name
		})

		for j := 0; j < len(bufST); j++ {
			brand := &bufST[j]
			var (
				type64   entry.Entry64
				brand64  entry.Entry64 // brand name index
				matchRI  = int32(-1)   // regex index
				match64  entry.Entry64 // string index
				modelSI  = int32(-1)   // single model index
				models64 entry.Entry64 // models index
			)
			brand64 = buf.add(brand.Name)
			type64 = buf.add(brand.Device)

			rs := brand.Regex
			if rs1, ok := staticReplDevice[rs]; ok {
				rs = rs1
			}
			if !isRegex(rs) {
				match64 = buf.add(brand.Regex)
			} else {
				rs = normalizeRegex(rs)
				if rs1, ok := reReplDevice[rs]; ok {
					rs = rs1
				}
				if _, err = regexp.Compile(rs); err == nil {
					bufRE = append(bufRE, rs)
					matchRI = int32(len(bufRE) - 1)
				} else {
					log.Printf("regexp error '%s' on '%s'", err, rs)
				}
			}

			var (
				type641  entry.Entry64
				matchRI1 = int32(-1)   // regex index
				match641 entry.Entry64 // string index
				model641 entry.Entry64 // model name index
			)
			if len(brand.Models) > 0 {
				meLO := uint32(len(bufDM))
				for k := 0; k < len(brand.Models); k++ {
					type641.Reset()
					matchRI1 = int32(-1)
					match641.Reset()
					model641.Reset()

					model := &brand.Models[k]
					rs1 := model.Regex
					if !isRegex(rs1) {
						match641 = buf.add(model.Regex)
					} else {
						rs1 = normalizeRegex(rs1)
						if rs1, ok := reReplDevice[rs]; ok {
							rs = rs1
						}
						if _, err = regexp.Compile(rs1); err == nil {
							bufRE = append(bufRE, rs1)
							matchRI1 = int32(len(bufRE) - 1)
						} else {
							log.Printf("regexp error '%s' on '%s'", err, rs)
						}
					}
					type641 = buf.add(model.Device)
					model641 = buf.add(model.Model)
					bufDM = append(bufDM, fmt.Sprintf("modelTuple{type64:%s,matchRI:%s,match64:%s,model64:%s}",
						hex(type641), hex(matchRI1), hex(match641), hex(model641)))
				}
				meHI := uint32(len(bufDM))
				models64.Encode(meLO, meHI)
			} else if len(brand.Model) > 0 {
				model641 = buf.add(brand.Model)
				bufDM = append(bufDM, fmt.Sprintf("modelTuple{type64:%s,matchRI:%s,match64:%s,model64:%s}",
					hex(type641), hex(matchRI1), hex(match641), hex(model641)))
				modelSI = int32(len(bufDM)) - 1
			}

			bufDR = append(bufDR, fmt.Sprintf("deviceTuple{type64:%s,brand64:%s,matchRI:%s,match64:%s,modelSI:%s,models64:%s},",
				hex(type64), hex(brand64), hex(matchRI), hex(match64), hex(modelSI), hex(models64)))
		}

		_, _ = w.WriteString("// " + filepath.Base(files[i]) + "\n")
		_, _ = w.WriteString("{\n")
		for j := 0; j < len(bufDR); j++ {
			_, _ = w.WriteString(bufDR[j])
			_ = w.WriteByte('\n')
		}
		_, _ = w.WriteString("},\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__dr_dm = []modelTuple{\n")
	for i := 0; i < len(bufDM); i++ {
		_, _ = w.WriteString(bufDM[i])
		_, _ = w.WriteString(",\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__dr_re = []*regexp.Regexp{\n")
	for i := 0; i < len(bufRE); i++ {
		re := bufRE[i]
		// if _, ok := ciExclude[re]; !ok {
		// 	re = "(?i)" + re
		// }
		_, _ = w.WriteString("regexp.MustCompile(`(?i)" + re + "`),\n")
	}
	_, _ = w.WriteString("}\n")

	_, _ = w.WriteString("__dr_buf = []byte{\n")
	for i := 0; i < len(buf.buf); i++ {
		if i > 0 && i%16 == 0 {
			_ = w.WriteByte('\n')
		}
		_, _ = w.WriteString(fmt.Sprintf("0x%02x, ", buf.buf[i]))
	}
	_, _ = w.WriteString("\n}\n")

	_, _ = w.WriteString(")\n")

	source := w.Bytes()
	var fmtSource []byte
	if fmtSource, err = format.Source(source); err != nil {
		return
	}

	err = ioutil.WriteFile(target, fmtSource, 0644)

	return
}
