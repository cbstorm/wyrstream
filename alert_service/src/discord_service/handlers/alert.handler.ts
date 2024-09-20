import { codeBlock, DiscordEventAdapter, DiscordFromEventAdapter, IDiscordEventAdapterContext } from 'nithjs-discord';

export const AlertEventAdapter = new DiscordEventAdapter();
export const alert = new DiscordFromEventAdapter({
  name: 'alert',
  event: AlertEventAdapter,
  handler: async (ctx: IDiscordEventAdapterContext<any>) => {
    await ctx.SendTo('wyr_stream_alert', codeBlock(ctx.GetData()));
  },
});
