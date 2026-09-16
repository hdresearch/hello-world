import asyncio
import hmac
from typing import Annotated, Literal

import httpx
from fastapi import FastAPI, Header, HTTPException, Response, status

from .settings import Settings

settings = Settings.from_environment()
app = FastAPI(title="Vers hello-world backend", docs_url=None, redoc_url=None)


def require_frontend(authorization: Annotated[str | None, Header()] = None) -> None:
    expected = f"Bearer {settings.frontend_token}"
    if authorization is None or not hmac.compare_digest(authorization, expected):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Unauthorized")


async def call_data_api(method: Literal["GET", "POST"], path: str) -> dict[str, int]:
    headers = {"Authorization": f"Bearer {settings.data_api_token}"}
    last_error = "Database service unavailable"
    async with httpx.AsyncClient(timeout=httpx.Timeout(5.0, connect=3.0)) as client:
        for attempt in range(3):
            try:
                response = await client.request(method, f"{settings.data_api_url}{path}", headers=headers)
                if response.status_code < 500:
                    response.raise_for_status()
                    payload = response.json()
                    if not isinstance(payload.get("count"), int):
                        raise ValueError("Invalid database response")
                    return {"count": payload["count"]}
                last_error = f"Database service returned {response.status_code}"
            except (httpx.HTTPError, ValueError) as error:
                last_error = str(error)
            if attempt < 2:
                await asyncio.sleep(0.2 * (2**attempt))
    raise HTTPException(status_code=status.HTTP_502_BAD_GATEWAY, detail=last_error)


@app.get("/healthz", include_in_schema=False)
async def healthz() -> Response:
    return Response(content="ok\n", media_type="text/plain")


@app.get("/readyz")
async def readyz(authorization: Annotated[str | None, Header()] = None) -> dict[str, str]:
    require_frontend(authorization)
    await call_data_api("GET", "/visits")
    return {"status": "ready"}


@app.get("/api/visits")
async def get_visits(authorization: Annotated[str | None, Header()] = None) -> dict[str, int]:
    require_frontend(authorization)
    return await call_data_api("GET", "/visits")


@app.post("/api/visits")
async def increment_visits(authorization: Annotated[str | None, Header()] = None) -> dict[str, int]:
    require_frontend(authorization)
    return await call_data_api("POST", "/visits")
