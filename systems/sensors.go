package systems

import (
	"context"
	"errors"
	"log"
	"time"

	canpackets "github.com/Liquid-Propulsion/canpackets/go"
	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/database/timeseries"
	"github.com/Liquid-Propulsion/mainland-server/types"
	"github.com/d5/tengo/v2"
	"github.com/nakabonne/tstorage"
)

type SensorKey struct {
	NodeID   uint8
	SensorID uint8
}

type SensorData struct {
	Time          time.Time
	SensorDataRaw uint32
	SensorData    float32
	SensorType    canpackets.SensorType
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
	sensors map[SensorKey]types.Sensor
	channel chan SensorUpdate
}

func NewSensorsSystem() *SensorsSystem {
	sensors := new(SensorsSystem)
	sensors.data = make(map[SensorKey]SensorData)
	sensors.channel = make(chan SensorUpdate)
	return sensors
}

func (sensors *SensorsSystem) Reset() {
	sensors.sensors = make(map[SensorKey]types.Sensor)
	var sensorsFromDB []types.Sensor
	sql.Database.Find(&sensorsFromDB)
	for _, sensor := range sensorsFromDB {
		sensors.sensors[SensorKey{
			NodeID:   sensor.NodeID,
			SensorID: sensor.SensorID,
		}] = sensor
	}
}

func (sensors *SensorsSystem) Run() {
	channel := canbackend.CurrentCANBackend.SensorDataChannel()
	for {
		data := <-channel
		sensor, ok := sensors.sensors[SensorKey{uint8(data.NodeId), data.SensorId}]
		if !ok {
			log.Printf("Sensor not registered")
			continue
		}
		out, err := tengo.Eval(context.Background(), sensor.TransformCode, map[string]interface{}{
			"data": int64(data.SensorData),
		})
		var sensorDataTransformed float32 = 0.0

		if err == nil {
			switch v := out.(type) {
			case float32:
				sensorDataTransformed = v
			case float64:
				sensorDataTransformed = float32(v)
			case int:
				sensorDataTransformed = float32(v)
			case int64:
				sensorDataTransformed = float32(v)
			}
		} else {
			log.Printf("Couldn't transform Sensor Data: %s", err)
			continue
		}

		sensors.data[SensorKey{uint8(data.NodeId), data.SensorId}] = SensorData{
			Time:          time.Now(),
			SensorData:    sensorDataTransformed,
			SensorDataRaw: data.SensorData,
			SensorType:    data.SensorType,
		}
		if CurrentEngine.State() == types.ARMED {
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
						Value:     float64(sensorDataTransformed),
					},
				},
			})
			timeseries.Database.InsertRows([]tstorage.Row{
				{
					Metric: "sensor_raw",
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
}

func (sensors *SensorsSystem) GetLatestSensorData(nodeID uint8, sensorID uint8) (SensorData, error) {
	if val, ok := sensors.data[SensorKey{nodeID, sensorID}]; ok {
		return val, nil
	}
	return SensorData{}, errors.New("no Sensor Data found")
}

func (sensors *SensorsSystem) GetSensorData(startTime int64, endTime int64, nodeID uint8, sensorID uint8) ([]*tstorage.DataPoint, error) {
	return timeseries.Database.Select("sensor", []tstorage.Label{
		{Name: "node_id", Value: string(nodeID)},
		{Name: "sensor_id", Value: string(sensorID)},
	}, startTime, endTime)
}
