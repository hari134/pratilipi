package kafka

// KafkaConfig holds the configuration for Kafka producer and consumer.
type KafkaConfig struct {
    Brokers     []string   // List of Kafka brokers
    GroupID     string     // Consumer group ID
    Topic       string     // Single topic (optional)
    GroupTopics []string   // Multiple topics for consumer groups
}

// NewKafkaConfig initializes a new KafkaConfig with default values.
func NewKafkaConfig() *KafkaConfig {
    return &KafkaConfig{
        Brokers: []string{"localhost:9092"}, // Default brokers
        GroupID: "default-group",            // Default group ID
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

// SetTopic sets a single topic for the Kafka config.
func (kc *KafkaConfig) SetTopic(topic string) *KafkaConfig {
    kc.Topic = topic
    return kc
}

// SetGroupTopics sets multiple topics for the Kafka config.
func (kc *KafkaConfig) SetGroupTopics(topics ...string) *KafkaConfig {
    kc.GroupTopics = topics
    return kc
}
