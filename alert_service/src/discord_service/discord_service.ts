import { DiscordApp } from 'nithjs-discord';
import path = require('path');

export const AlertApp = new DiscordApp({
  discordBotToken: process.env.DISCORD_BOT_TOKEN!,
  handlerPath: path.join(process.cwd(), '.out', 'discord_service', 'handlers'),
  handlerPattern: '.handler.js',
});
