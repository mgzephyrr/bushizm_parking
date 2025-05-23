from jose import JWTError, jwt
import os
from dotenv import load_dotenv

load_dotenv()

BOT_TOKEN = os.getenv("UJIN_BOT_TOKEN")
CHANNEL_KEY = os.getenv("CHANNEL_KEY")
AUTH_SERVICE_URL = os.getenv("AUTH_SERVICE_URL")
SECRET_KEY = os.getenv("SECRET_KEY")
ALGORITHM = os.getenv("ALGORITHM")

def verify_token(token: str):
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        return payload.get("sub")
    except JWTError:
        return None