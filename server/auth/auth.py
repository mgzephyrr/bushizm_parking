import os
import httpx
from dotenv import load_dotenv
from fastapi import HTTPException, Response, Request
from utils import create_access_token, verify_token

load_dotenv()

API_TOKEN = os.getenv("UJIN_API_TOKEN")
API_URL = "https://api-uae-test.ujin.tech/api/admin/get-userdata"


def validate_phone(phone: str):
    if not phone.isdigit():
        raise HTTPException(
            status_code=422, detail="Phone number must contain only digits")
    if len(phone) != 11 or not phone.startswith("7"):
        raise HTTPException(
            status_code=422, detail="Phone number must be 11 digits and start with '7'")


async def authenticate_user_by_phone(phone: str, response: Response):
    validate_phone(phone)

    if not API_TOKEN:
        raise HTTPException(
            status_code=500, detail="API token not set in environment")

    params = {
        "app": "crm",
        "token": API_TOKEN,
        "limit": 10,
        "offset": 0,
        "timezone": "Asia/Yekaterinburg"
    }

    timeout = httpx.Timeout(15.0, connect=2.0, read=5.0)
    limits = httpx.Limits(max_connections=10, max_keepalive_connections=5)

    async with httpx.AsyncClient(timeout=timeout, limits=limits) as client:
        try:
            api_response = await client.get(API_URL, params=params)
            api_response.raise_for_status()  # выбросит ошибку если статус не 2xx
            data = api_response.json()
        except httpx.TimeoutException:
            raise HTTPException(status_code=504, detail="External API timeout")
        except httpx.HTTPStatusError as e:
            raise HTTPException(status_code=e.response.status_code,
                                detail=f"External API error: {e.response.text}")
        except Exception as e:
            raise HTTPException(
                status_code=500, detail=f"Unexpected error: {str(e)}")

    users_data = data.get("data", {}).get("userdata", [])
    if not users_data or not isinstance(users_data, list):
        raise HTTPException(
            status_code=404, detail=f"No valid 'userdata' found in API response: {data}")

    for user in users_data:
        if user.get("phone") == phone:
            user_id = str(user.get("id"))
            full_name = f"{user.get('surname', '')} {user.get('name', '')} {user.get('patronymic', '')}".strip(
            )
            jwt_token = create_access_token({"sub": user_id})

            response.set_cookie(
                key="access_token",
                value=jwt_token,
                httponly=True,
                max_age=3600,
                secure=False,
                samesite="Lax"
            )

            return {
                "message": "Authenticated",
                "user_id": user_id,
                "full_name": full_name
            }

    raise HTTPException(
        status_code=404, detail=f"User with phone {phone} not found")


def logout_user(response: Response):
    response.delete_cookie("access_token")
    return {"message": "Logged out successfully"}


def get_user_id_from_cookie(request: Request):
    token = request.cookies.get("access_token")
    if not token:
        raise HTTPException(status_code=401, detail="Missing auth token")

    user_id = verify_token(token)
    if not user_id:
        raise HTTPException(status_code=401, detail="Invalid or expired token")

    return user_id
