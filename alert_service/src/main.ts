import config from './configs/config';
import nats_service from './nats_service/nats_service';

(async function main() {
  process.on('SIGINT', () => {
    nats_service.Close();
  });
  await nats_service.LoadConfig(config).Connect();
})();
