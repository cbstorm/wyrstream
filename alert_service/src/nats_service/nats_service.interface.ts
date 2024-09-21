export interface INATSConfig {
  LoadNATSCoreConfig: () => void;
  NATS_CORE_HOST: string;
  NATS_CORE_PORT: string;
  NATS_CORE_USERNAME: string;
  NATS_CORE_PASSWORD: string;
  NATS_CORE_QUEUE_GROUP: string;
}
