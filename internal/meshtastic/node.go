package meshtastic

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"time"
)

const (
	SeenByLimit   = 10
	NeighborLimit = 100
)

func cleanFloat(f float32) float32 {
	if f != f {
		// IEEE 754 says that only NaNs satisfy f != f
		return 0
	}
	// limit floats to 3 decimal places
	return float32(math.Round(float64(f*1000)) / 1000)
}

type NeighborInfo struct {
	Snr     float32 `json:"snr,omitempty"`
	Updated int64   `json:"updated"`
}

type Node struct {
	// User
	LongName  string `json:"longName"`
	ShortName string `json:"shortName"`
	HwModel   string `json:"hwModel"`
	Role      string `json:"role"`
	// MapReport
	FwVersion        string `json:"fwVersion,omitempty"`
	Region           string `json:"region,omitempty"`
	ModemPreset      string `json:"modemPreset,omitempty"`
	HasDefaultCh     bool   `json:"hasDefaultCh,omitempty"`
	OnlineLocalNodes uint32 `json:"onlineLocalNodes,omitempty"`
	LastMapReport    int64  `json:"lastMapReport,omitempty"`
	// Position
	Latitude  int32  `json:"latitude"`
	Longitude int32  `json:"longitude"`
	Altitude  int32  `json:"altitude,omitempty"`
	Precision uint32 `json:"precision,omitempty"`
	// DeviceMetrics
	BatteryLevel      uint32  `json:"batteryLevel,omitempty"`
	Voltage           float32 `json:"voltage,omitempty"`
	ChUtil            float32 `json:"chUtil,omitempty"`
	AirUtilTx         float32 `json:"airUtilTx,omitempty"`
	Uptime            uint32  `json:"uptime,omitempty"`
	LastDeviceMetrics int64   `json:"lastDeviceMetrics,omitempty"`
	// EnvironmentMetrics
	Temperature            float32 `json:"temperature,omitempty"`
	RelativeHumidity       float32 `json:"relativeHumidity,omitempty"`
	BarometricPressure     float32 `json:"barometricPressure,omitempty"`
	Lux                    float32 `json:"lux,omitempty"`
	WindDirection          uint32  `json:"windDirection,omitempty"`
	WindSpeed              float32 `json:"windSpeed,omitempty"`
	WindGust               float32 `json:"windGust,omitempty"`
	Radiation              float32 `json:"radiation,omitempty"`
	Rainfall1              float32 `json:"rainfall1,omitempty"`
	Rainfall24             float32 `json:"rainfall24,omitempty"`
	LastEnvironmentMetrics int64   `json:"lastEnvironmentMetrics,omitempty"`
	// NeighborInfo
	Neighbors map[uint32]*NeighborInfo `json:"neighbors,omitempty"`
	// key=mqtt topic, value=first seen/last position update
	SeenBy map[string]int64 `json:"seenBy"`
}

func NewNode(topic string) *Node {
	return &Node{
		SeenBy: map[string]int64{topic: time.Now().Unix()},
	}
}

func (node *Node) ClearDeviceMetrics() {
	node.BatteryLevel = 0
	node.Voltage = 0
	node.ChUtil = 0
	node.AirUtilTx = 0
	node.Uptime = 0
	node.LastDeviceMetrics = 0
}

func (node *Node) ClearEnvironmentMetrics() {
	node.Temperature = 0
	node.RelativeHumidity = 0
	node.BarometricPressure = 0
	node.Lux = 0
	node.WindDirection = 0
	node.WindSpeed = 0
	node.WindGust = 0
	node.Radiation = 0
	node.Rainfall1 = 0
	node.Rainfall24 = 0
	node.LastEnvironmentMetrics = 0
}

func (node *Node) ClearMapReportData() {
	node.FwVersion = ""
	node.Region = ""
	node.ModemPreset = ""
	node.HasDefaultCh = false
	node.OnlineLocalNodes = 0
	node.LastMapReport = 0
}

func (node *Node) IsValid() bool {
	if len(node.SeenBy) == 0 {
		return false
	}
	if len(node.LongName) == 0 {
		return false
	}
	if node.Latitude == 0 && node.Longitude == 0 {
		return false
	}
	return true
}

func (node *Node) Prune(seenByTtl, neighborTtl, metricsTtl, mapReportTtl int64) {
	now := time.Now().Unix()
	// SeenBy
	for topic, lastSeen := range node.SeenBy {
		if lastSeen+seenByTtl < now {
			delete(node.SeenBy, topic)
		}
	}
	for len(node.SeenBy) > SeenByLimit {
		var toDelete string
		for topic, lastSeen := range node.SeenBy {
			if len(toDelete) == 0 || lastSeen < node.SeenBy[toDelete] {
				toDelete = topic
			}
		}
		delete(node.SeenBy, toDelete)
	}
	// Neighbors
	for neighborNum, neighbor := range node.Neighbors {
		if neighbor.Updated+neighborTtl < now {
			delete(node.Neighbors, neighborNum)
		}
	}
	if len(node.Neighbors) == 0 {
		node.Neighbors = nil
	}
	for len(node.Neighbors) > NeighborLimit {
		var toDelete uint32
		for neighborNum, neighbor := range node.Neighbors {
			if toDelete == 0 || neighbor.Updated < node.Neighbors[toDelete].Updated {
				toDelete = neighborNum
			}
		}
		delete(node.Neighbors, toDelete)
	}
	// DeviceMetrics
	if node.LastDeviceMetrics > 0 && node.LastDeviceMetrics+metricsTtl < now {
		node.ClearDeviceMetrics()
	}
	// EnvironmentMetrics
	if node.LastEnvironmentMetrics > 0 && node.LastEnvironmentMetrics+metricsTtl < now {
		node.ClearEnvironmentMetrics()
	}
	// MapReport
	if node.LastMapReport > 0 && node.LastMapReport+mapReportTtl < now {
		node.ClearMapReportData()
	}
}

func (node *Node) UpdateDeviceMetrics(batteryLevel uint32, voltage, chUtil, airUtilTx float32, uptime uint32) {
	node.BatteryLevel = batteryLevel
	node.Voltage = cleanFloat(voltage)
	node.ChUtil = cleanFloat(chUtil)
	node.AirUtilTx = cleanFloat(airUtilTx)
	node.Uptime = uptime
	node.LastDeviceMetrics = time.Now().Unix()
}

func (node *Node) UpdateEnvironmentMetrics(temperature, relativeHumidity, barometricPressure, lux float32, windDirection uint32, windSpeed, windGust, radiation, rainfall1, rainfall24 float32) {
	node.Temperature = cleanFloat(temperature)
	node.RelativeHumidity = cleanFloat(relativeHumidity)
	node.BarometricPressure = cleanFloat(barometricPressure)
	node.Lux = cleanFloat(lux)
	node.WindDirection = windDirection
	node.WindSpeed = cleanFloat(windSpeed)
	node.WindGust = cleanFloat(windGust)
	node.Radiation = cleanFloat(radiation)
	node.Rainfall1 = cleanFloat(rainfall1)
	node.Rainfall24 = cleanFloat(rainfall24)
	node.LastEnvironmentMetrics = time.Now().Unix()
}

func (node *Node) UpdateMapReport(fwVersion, region, modemPreset string, hasDefaultCh bool, onlineLocalNodes uint32) {
	node.FwVersion = fwVersion
	node.Region = region
	node.ModemPreset = modemPreset
	node.HasDefaultCh = hasDefaultCh
	node.OnlineLocalNodes = onlineLocalNodes
	node.LastMapReport = time.Now().Unix()
}

func (node *Node) UpdateNeighborInfo(neighborNum uint32, snr float32) {
	if node.Neighbors == nil {
		node.Neighbors = make(map[uint32]*NeighborInfo)
	}
	node.Neighbors[neighborNum] = &NeighborInfo{
		Snr:     cleanFloat(snr),
		Updated: time.Now().Unix(),
	}
}

func (node *Node) UpdatePosition(latitude, longitude, altitude int32, precision uint32) {
	node.Latitude = latitude
	node.Longitude = longitude
	node.Altitude = altitude
	node.Precision = precision
}

func (node *Node) UpdateSeenBy(topic string) {
	node.SeenBy[topic] = time.Now().Unix()
}

func (node *Node) UpdateUser(longName, shortName, hwModel, role string) {
	node.LongName = longName
	node.ShortName = shortName
	node.HwModel = hwModel
	node.Role = role
}

type NodeDB map[uint32]*Node

func (db NodeDB) Prune(seenByTtl, neighborTtl, metricsTtl, mapReportTtl int64) {
	for nodeNum, node := range db {
		node.Prune(seenByTtl, neighborTtl, metricsTtl, mapReportTtl)
		if len(node.SeenBy) == 0 {
			delete(db, nodeNum)
		}
	}
}

func (db NodeDB) GetValid() NodeDB {
	valid := make(NodeDB)
	for nodeNum, node := range db {
		if node.IsValid() {
			valid[nodeNum] = node
		}
	}
	return valid
}

func (db *NodeDB) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(db)
}

func (db NodeDB) WriteFile(path string) error {
	dir, file := filepath.Split(path)
	f, err := os.CreateTemp(dir, file)
	if err != nil {
		return err
	}
	err = json.NewEncoder(f).Encode(db)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err == nil {
		err = os.Chmod(f.Name(), 0644)
	}
	if err == nil {
		err = os.Rename(f.Name(), path)
	}
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}
