import { INATSConfig } from '../nats_service/nats_service.interface';

class Config implements INATSConfig {
  public NATS_CORE_HOST!: string;
  public NATS_CORE_PASSWORD!: string;
  public NATS_CORE_PORT!: string;
  public NATS_CORE_USERNAME!: string;
  constructor() {}
  LoadNATSCoreConfig() {
    this.NATS_CORE_HOST = process.env.NATS_CORE_HOST!;
    this.NATS_CORE_PORT = process.env.NATS_CORE_PORT!;
    this.NATS_CORE_PASSWORD = process.env.NATS_CORE_PASSWORD!;
    this.NATS_CORE_USERNAME = process.env.NATS_CORE_USERNAME!;
  }
}

export default new Config();
