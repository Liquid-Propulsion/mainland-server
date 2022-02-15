package systems

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/database/timeseries"
	"github.com/Liquid-Propulsion/mainland-server/types"
	"github.com/d5/tengo/v2"
	"github.com/nakabonne/tstorage"
)

type SensorData struct {
	Time          time.Time
	SensorDataRaw uint32
	SensorData    float32
}

type SensorUpdate struct {
	Time       time.Time
	SensorID   uint
	SensorData uint32
}

type SensorsSystem struct {
	data    map[uint]SensorData
	sensors map[uint]types.Sensor
	channel chan SensorUpdate
}

func NewSensorsSystem() *SensorsSystem {
	sensors := new(SensorsSystem)
	sensors.data = make(map[uint]SensorData)
	sensors.channel = make(chan SensorUpdate)
	return sensors
}

func (sensors *SensorsSystem) Reset() {
	sensors.data = make(map[uint]SensorData)
	sensors.sensors = make(map[uint]types.Sensor)
	var sensorsFromDB []types.Sensor
	sql.Database.Find(&sensorsFromDB)
	for _, sensor := range sensorsFromDB {
		sensors.sensors[sensor.ID] = sensor
	}
}

func (sensors *SensorsSystem) Run() {
	channel := canbackend.CurrentCANBackend.SensorDataChannel()
	for {
		data := <-channel
		sensor, ok := sensors.sensors[uint(data.SensorId)]
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

		sensors.data[uint(data.SensorId)] = SensorData{
			Time:          time.Now(),
			SensorData:    sensorDataTransformed,
			SensorDataRaw: data.SensorData,
		}
		if CurrentEngine.State() == types.ARMED {
			timeseries.Database.InsertRows([]tstorage.Row{
				{
					Metric: "sensor",
					Labels: []tstorage.Label{
						{Name: "sensor_id", Value: string(data.SensorId)},
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
						{Name: "sensor_id", Value: string(data.SensorId)},
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

func (sensors *SensorsSystem) GetLatestSensorData(sensorID uint) (SensorData, error) {
	if val, ok := sensors.data[uint(sensorID)]; ok {
		return val, nil
	}
	return SensorData{}, errors.New("no Sensor Data found")
}

func (sensors *SensorsSystem) GetSensorData(startTime int64, endTime int64, sensorID uint) ([]*tstorage.DataPoint, error) {
	return timeseries.Database.Select("sensor", []tstorage.Label{
		{Name: "sensor_id", Value: fmt.Sprint(sensorID)},
	}, startTime, endTime)
}
