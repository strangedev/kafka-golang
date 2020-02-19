package producer

import "encoding/json"

type JSONProducer interface {
	ProduceJSONSync(topic string, partition int32, value interface{}) error
}

func (k *KafkaProducer) ProduceJSONSync(topic string, partition int32, value interface{}) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return k.ProduceSimpleSync(topic, partition, marshaled)
}
