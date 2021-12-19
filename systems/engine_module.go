package systems

import (
	"log"

	"github.com/Liquid-Propulsion/mainland-server/types"
	"github.com/d5/tengo/v2"
)

func EngineModule() map[string]tengo.Object {
	return map[string]tengo.Object{
		"sensor_data_raw": &tengo.UserFunction{Value: sensorDataRaw},
		"sensor_data":     &tengo.UserFunction{Value: sensorData},
		"safe_system":     &tengo.UserFunction{Value: safeSystem},
		"log_warn":        &tengo.UserFunction{Value: logWarn},
		"log_error":       &tengo.UserFunction{Value: logError},
	}
}

// Definition: `sensor_data_raw(node_id int, sensor_id int) int`
func sensorDataRaw(args ...tengo.Object) (tengo.Object, error) {
	if len(args) >= 2 {
		switch node_id := args[0].(type) {
		case *tengo.Int:
			switch sensor_id := args[1].(type) {
			case *tengo.Int:
				data, err := CurrentEngine.SensorsSystem.GetLatestSensorData(uint8(node_id.Value), uint8(sensor_id.Value))
				return &tengo.Int{Value: int64(data.SensorDataRaw)}, err
			default:
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "second",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
		default:
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "int",
				Found:    args[0].TypeName(),
			}
		}
	}
	return nil, tengo.ErrWrongNumArguments
}

// Definition: `sensor_data(node_id int, sensor_id int) float`
func sensorData(args ...tengo.Object) (tengo.Object, error) {
	if len(args) >= 2 {
		switch node_id := args[0].(type) {
		case *tengo.Int:
			switch sensor_id := args[1].(type) {
			case *tengo.Int:
				data, err := CurrentEngine.SensorsSystem.GetLatestSensorData(uint8(node_id.Value), uint8(sensor_id.Value))
				return &tengo.Float{Value: float64(data.SensorData)}, err
			default:
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "second",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
		default:
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "int",
				Found:    args[0].TypeName(),
			}
		}
	}
	return nil, tengo.ErrWrongNumArguments
}

// Definition: `safe_system()`, forces the system into safe state.
func safeSystem(args ...tengo.Object) (tengo.Object, error) {
	return tengo.UndefinedValue, CurrentEngine.SetState(types.SAFE)
}

// Definition: `log_error(string)`, logs an error.
func logError(args ...tengo.Object) (tengo.Object, error) {
	if len(args) >= 1 {
		switch v := args[0].(type) {
		case *tengo.String:
			log.Printf("Error: %s", v.Value)
			return tengo.UndefinedValue, nil
		default:
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string",
				Found:    args[0].TypeName(),
			}
		}
	}
	return nil, tengo.ErrWrongNumArguments
}

// Definition: `log_warn(string)`, logs an error.
func logWarn(args ...tengo.Object) (tengo.Object, error) {
	if len(args) >= 1 {
		switch v := args[0].(type) {
		case *tengo.String:
			log.Printf("Warn: %s", v.Value)
			return tengo.UndefinedValue, nil
		default:
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string",
				Found:    args[0].TypeName(),
			}
		}
	}
	return nil, tengo.ErrWrongNumArguments
}
