package kafka

// KafkaConfig holds the configuration for Kafka producer and consumer.
type KafkaConfig struct {
    Brokers []string
    GroupID string
    Topic   string
}

// NewKafkaConfig initializes a new KafkaConfig with default values.
func NewKafkaConfig() *KafkaConfig {
    return &KafkaConfig{
        Brokers: []string{"localhost:9092"}, // Default brokers
        GroupID: "default-group",            // Default group ID
        Topic:   "default-topic",            // Default topic
    }
}

// SetBrokers sets the Kafka brokers for the config.
func (kc *KafkaConfig) SetBrokers(brokers ...string) *KafkaConfig {
    kc.Brokers = brokers
    return kc
}

// SetGroupID sets the group ID for the Kafka config.
func (kc *KafkaConfig) SetGroupID(groupID string) *KafkaConfig {
    kc.GroupID = groupID
    return kc
}

// SetTopic sets the topic for the Kafka config.
func (kc *KafkaConfig) SetTopic(topic string) *KafkaConfig {
    kc.Topic = topic
    return kc
}
