import { Msg } from 'nats';
import config from './configs/config';
import { AlertApp } from './discord_service/discord_service';
import { AlertEventAdapter } from './discord_service/handlers/alert.handler';
import nats_service from './nats_service/nats_service';

(async function main() {
  process.on('SIGINT', () => {
    nats_service.Close();
  });
  await nats_service.LoadConfig(config).Connect();
  await nats_service.Regis('alert', async (m: Msg) => {
    AlertEventAdapter.Emit(m.data);
  });
  (await AlertApp.Init())
    .LoadHandler((ev: string) => {
      console.log(ev + ' Loaded');
    })
    .Listen((c) => {
      console.log(c.user.displayName + ' Started');
    });
})();
