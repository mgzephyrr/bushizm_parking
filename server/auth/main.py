from fastapi import Body, FastAPI, HTTPException, Response, Request
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from auth import authenticate_user_by_phone, get_user_id_from_cookie, logout_user
from utils import verify_token

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173", "http://localhost:5174", "http://localhost:8001", "http://localhost:8080", "http://localhost:8000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class AuthRequest(BaseModel):
    phone: str

@app.post("/login")
async def login(data: AuthRequest, response: Response):
    return await authenticate_user_by_phone(data.phone, response)

@app.post("/logout")
def logout(response: Response):
    return logout_user(response)

@app.post("/extract_user_id")
def extract_user_id(token: str = Body(..., embed=True)):
    user_id = verify_token(token)
    if not user_id:
        raise HTTPException(status_code=401, detail="Invalid or expired token")
    return {"user_id": user_id}


@app.get("/me")
def get_current_user(request: Request):
    user_id = get_user_id_from_cookie(request)
    return {"user_id": user_id}

# async def create_new_ticket(token: str):
#     api_url = f'https://api-uae-test.ujin.tech/api/v1/tck/bms/tickets/create/?token={token}'
    
#     async with httpx.AsyncClient() as client:
#         respns = await client.post(api_url, json={
#             "title": "Заявка на сантехническое обслуживание",
#             "description": "Есть вероятность засора канализации, необходимо вызвать сантехническую службу",
#             "priority": "high",
#             "class": "inspection",
#             "status": "new",
#             "initiator.id": 739111,
#             "types": [],
#             "assignees": [],
#             "contracting_companies": [],
#             "objects": [
#                 {
#                 "type": "building",
#                 "id": 47
#                 }
#             ],
#             "planned_start_at": "",
#             "planned_end_at": "",
#             "hide_planned_at_from_resident": "",
#             "extra": ""
#         })

#     respns = respns.json()

#     print(respns)
#     return respns