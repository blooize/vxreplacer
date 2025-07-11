import discord
from discord.ext import commands
import logging
from dotenv import load_dotenv
import os
import re

load_dotenv()
token = os.getenv('DISCORD_TOKEN')

# Create logs directory if it doesn't exist
os.makedirs('logs', exist_ok=True)
handler = logging.FileHandler(filename='logs/discord.log', encoding='utf-8', mode='w')
intents = discord.Intents.default()
intents.message_content = True
intents.members = True

bot = commands.Bot(command_prefix='!', intents = intents)

twitter_link_pattern = re.compile(r"(https?://)?(www\.)?(twitter\.com/\S+)", re.IGNORECASE)
x_link_pattern = re.compile(r"(https?://)?(www\.)?(x\.com/\S+)", re.IGNORECASE)


@bot.event
async def on_message(message):
    if message.author == bot.user:
        return
    
    matches_twitter = twitter_link_pattern.findall(message.content)
    matches_x = x_link_pattern.findall(message.content)
    if (matches_twitter or matches_x):
        modified_links = []
        for match in matches_twitter:
            full_url = ''.join(match)
            vx_url = full_url.replace("twitter.com","vxtwitter.com")
            modified_links.append(vx_url)

        for match in matches_x:
            full_url = ''.join(match)
            vx_url = full_url.replace("x.com","vxtwitter.com")
            modified_links.append(vx_url)


        if modified_links:
            await message.channel.send("\n".join(modified_links))

    await bot.process_commands(message)

bot.run(token, log_handler=handler, log_level=logging.DEBUG)
