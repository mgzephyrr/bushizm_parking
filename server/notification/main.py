from fastapi import FastAPI, Request, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import httpx
from pydantic import BaseModel
from utils import BOT_TOKEN, AUTH_SERVICE_URL, verify_token

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173", "http://localhost:5174", "http://localhost:8000", "http://localhost:8080"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


queued_user_ids = set()

USER_CHAT_MAP = {
    "2814589": "35631"
}

class UserRequest(BaseModel):
    user_id: str

@app.post("/queue")
async def add_to_queue(data: dict):
    user_id = data.get("user_id")
    if not user_id:
        raise HTTPException(status_code=400, detail="Missing user_id")

    queued_user_ids.add(str(user_id))
    return {"status": "queued", "user_id": user_id}

@app.post("/notify")
async def notify(request_data: UserRequest):
    try:
        if request_data.user_id not in queued_user_ids:
            return {"send": "no"}

        queued_user_ids.remove(str(user_id))

        chat_id = USER_CHAT_MAP.get(str(user_id))
        if not chat_id:
            raise HTTPException(status_code=404, detail="Chat ID not found")

        api_url = "https://api-uae-test.ujin.tech/sendMessage"
        payload = {
            "channel_key": chat_id,
            "text": "На парковке появилось свободное место.",
            "attachment": {},
            "is_hidden": False
        }
        headers = {
            "Authorization": f"Token {BOT_TOKEN}"
        }

        async with httpx.AsyncClient() as client:
            send_response = await client.post(api_url, json=payload, headers=headers)
            send_response.raise_for_status()

        return {"send": "yes"}

    except httpx.HTTPStatusError as e:
        raise HTTPException(status_code=e.response.status_code, detail=e.response.text)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))