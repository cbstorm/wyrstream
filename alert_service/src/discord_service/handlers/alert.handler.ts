import { codeBlock, DiscordEventAdapter, DiscordFromEventAdapter, IDiscordEventAdapterContext } from 'nithjs-discord';

export const AlertEventAdapter = new DiscordEventAdapter();
export const alert = new DiscordFromEventAdapter({
  name: 'alert',
  event: AlertEventAdapter,
  handler: async (ctx: IDiscordEventAdapterContext<any>) => {
    await ctx.SendTo('wyr-stream-alert', codeBlock(JSON.stringify(JSON.parse(ctx.GetData()), null, 2)));
  },
});
