from fastapi import FastAPI, Request, HTTPException, Body
from fastapi.middleware.cors import CORSMiddleware
import httpx
from pydantic import BaseModel
from utils import BOT_TOKEN, AUTH_SERVICE_URL, verify_token

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173", "http://localhost:5174",
                   "http://localhost:8000", "http://localhost:8080", "http://localhost:8001"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

queued_user_ids = set()

USER_CHAT_MAP = {
    "2814589": "35631"
}


@app.post("/queue")
async def add_to_queue(user_id: str = Body(..., embed=True)):
    queued_user_ids.add(user_id)
    print(queued_user_ids)
    return {"status": "queued", "user_id": user_id}


@app.post("/notify")
async def notify(request: Request):
    try:
        token = request.cookies.get("access_token")

        if not token:
            raise HTTPException(status_code=401, detail="Missing access_token")

        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{AUTH_SERVICE_URL}/extract_user_id",
                json={"token": token}
            )

        if response.status_code != 200:
            raise HTTPException(status_code=401, detail="Invalid token")

        user_id = response.json().get("user_id")
        if not user_id:
            raise HTTPException(status_code=401, detail="No user_id extracted")

        if str(user_id) not in queued_user_ids:
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

        # print("HERE!")

        timeout = httpx.Timeout(10.0, connect=5.0, read=10.0)

        async with httpx.AsyncClient(timeout=timeout) as client:
            send_response = await client.post(api_url, json=payload, headers=headers)
            send_response.raise_for_status()

        # print("HERE! END!")

        return {"send": "yes"}

    except httpx.HTTPStatusError as e:
        raise HTTPException(
            status_code=e.response.status_code, detail=e.response.text)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
