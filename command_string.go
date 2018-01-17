// Code generated by "stringer -type Command"; DO NOT EDIT.

package medtronic

import "strconv"

const _Command_name = "acknaksetClockbolussetAbsoluteTempBasalsuspendbuttonwakeupsetPercentTempBasalclockpumpIDbatteryreservoirhistoryPagerfRemoteSrfRemoteACTrfRemoteBglucoseUnitscarbRatiosinsulinSensitivitiesmodelbasalRatesbasalPatternAbasalPatternBtempBasalhistoryPageCountglucoseTargetssettingsstatus"

var _Command_map = map[Command]string{
	6:   _Command_name[0:3],
	21:  _Command_name[3:6],
	64:  _Command_name[6:14],
	66:  _Command_name[14:19],
	76:  _Command_name[19:39],
	77:  _Command_name[39:46],
	91:  _Command_name[46:52],
	93:  _Command_name[52:58],
	105: _Command_name[58:77],
	112: _Command_name[77:82],
	113: _Command_name[82:88],
	114: _Command_name[88:95],
	115: _Command_name[95:104],
	128: _Command_name[104:115],
	129: _Command_name[115:124],
	134: _Command_name[124:135],
	136: _Command_name[135:144],
	137: _Command_name[144:156],
	138: _Command_name[156:166],
	139: _Command_name[166:186],
	141: _Command_name[186:191],
	146: _Command_name[191:201],
	147: _Command_name[201:214],
	148: _Command_name[214:227],
	152: _Command_name[227:236],
	157: _Command_name[236:252],
	159: _Command_name[252:266],
	192: _Command_name[266:274],
	206: _Command_name[274:280],
}

func (i Command) String() string {
	if str, ok := _Command_map[i]; ok {
		return str
	}
	return "Command(" + strconv.FormatInt(int64(i), 10) + ")"
}
