import { Msg, NatsConnection, connect } from 'nats';
import { INATSConfig } from './nats_service.interface';

class NATSService {
  private _nats_conn!: NatsConnection;
  private _config!: INATSConfig;
  constructor() {}

  LoadConfig(config: INATSConfig) {
    config.LoadNATSCoreConfig();
    this._config = config;
    return this;
  }
  async Connect() {
    try {
      this._nats_conn = await connect({
        servers: [`nats://${this._config.NATS_CORE_HOST}:${this._config.NATS_CORE_PORT}`],
        user: this._config.NATS_CORE_USERNAME,
        pass: this._config.NATS_CORE_PASSWORD,
        reconnectTimeWait: 5000,
      });
      console.log('Connect to NATS server successfully');
    } catch (error) {
      console.log(error, 'Could not connect to NATS server');
      throw error;
    }
  }
  GetConnection() {
    return this._nats_conn;
  }
  async Flush() {
    return await this._nats_conn.flush();
  }
  async Close() {
    await this.Flush();
    return await this._nats_conn.close();
  }
  async Regis(queue: string, handler: (m: Msg) => void) {
    const sub = this._nats_conn.subscribe(queue, { queue: this._config.NATS_CORE_QUEUE_GROUP });
    (async () => {
      for await (const m of sub) {
        handler(m);
      }
    })();
  }
}

export default new NATSService();
