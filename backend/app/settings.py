from dataclasses import dataclass
import os


def required(name: str) -> str:
    value = os.environ.get(name, "").strip()
    if not value:
        raise RuntimeError(f"{name} is required")
    return value


@dataclass(frozen=True)
class Settings:
    frontend_token: str
    data_api_url: str
    data_api_token: str

    @classmethod
    def from_environment(cls) -> "Settings":
        return cls(
            frontend_token=required("FRONTEND_TOKEN"),
            data_api_url=required("DATA_API_URL").rstrip("/"),
            data_api_token=required("DATA_API_TOKEN"),
        )
