package systems

import (
	"time"

	canpackets "github.com/Liquid-Propulsion/canpackets/go"
	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/database/timeseries"
	"github.com/nakabonne/tstorage"
)

type SensorKey struct {
	NodeID   uint8
	SensorID uint8
}

type SensorData struct {
	Time       time.Time
	SensorData uint32
	SensorType canpackets.SensorType
}

type SensorUpdate struct {
	Time       time.Time
	NodeID     uint8
	SensorID   uint8
	SensorData uint32
	SensorType canpackets.SensorType
}

type SensorsSystem struct {
	data    map[SensorKey]SensorData
	channel chan SensorUpdate
}

func NewSensorsSystem() *SensorsSystem {
	sensors := new(SensorsSystem)
	sensors.data = make(map[SensorKey]SensorData)
	sensors.channel = make(chan SensorUpdate)
	return sensors
}

func (sensors *SensorsSystem) Run() {
	channel := canbackend.CurrentCANBackend.SensorDataChannel()
	for {
		data := <-channel
		sensors.channel <- SensorUpdate{
			Time:       time.Now(),
			NodeID:     uint8(data.NodeId),
			SensorID:   data.SensorId,
			SensorData: data.SensorData,
			SensorType: data.SensorType,
		}
		sensors.data[SensorKey{uint8(data.NodeId), data.SensorId}] = SensorData{
			Time:       time.Now(),
			SensorData: data.SensorData,
			SensorType: data.SensorType,
		}
		timeseries.Database.InsertRows([]tstorage.Row{
			{
				Metric: "sensor",
				Labels: []tstorage.Label{
					{Name: "node_id", Value: string(data.NodeId)},
					{Name: "sensor_id", Value: string(data.SensorId)},
					{Name: "sensor_type", Value: string(data.SensorType)},
				},
				DataPoint: tstorage.DataPoint{
					Timestamp: time.Now().UnixNano(),
					Value:     float64(data.SensorData),
				},
			},
		})
	}
}

func (sensors *SensorsSystem) GetLatestSensorData(nodeID uint8, sensorID uint8) (SensorData, error) {
	if val, ok := sensors.data[SensorKey{nodeID, sensorID}]; ok {
		return val, nil
	}
	return SensorData{}, nil
}

func (sensors *SensorsSystem) GetSensorData(startTime int64, endTime int64, nodeID uint8, sensorID uint8) ([]*tstorage.DataPoint, error) {
	return timeseries.Database.Select("sensor", []tstorage.Label{
		{Name: "node_id", Value: string(nodeID)},
		{Name: "sensor_id", Value: string(sensorID)},
	}, startTime, endTime)
}