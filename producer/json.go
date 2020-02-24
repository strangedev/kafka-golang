package producer

import "encoding/json"

// JSONProducer is capable of producing messages whose body is JSON encoded.
type JSONProducer interface {
	// ProduceJSONSync produces a message without headers and key by JSON-encoding the given value.
	ProduceJSONSync(topic string, partition int32, value interface{}) error
}

func (k *KafkaProducer) ProduceJSONSync(topic string, partition int32, value interface{}) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return k.ProduceSimpleSync(topic, partition, marshaled)
}
