import os
from dotenv import load_dotenv

load_dotenv()

BOT_TOKEN = os.getenv("UJIN_BOT_TOKEN")
CHANNEL_KEY = os.getenv("CHANNEL_KEY")
AUTH_SERVICE_URL = os.getenv("AUTH_SERVICE_URL")