from typing import Optional
from pydantic import BaseModel, Field
from datetime import datetime

class ItemBase(BaseModel):
    title: str = Field(..., min_length=1, max_length=120, description="Título del item")
    description: Optional[str] = Field(None, max_length=500, description="Descripción opcional")
    status: str = Field(default="pending", description="pending | in_progress | completed")

class ItemCreate(ItemBase):
    pass

class ItemUpdate(BaseModel):
    title: Optional[str] = Field(None, min_length=1, max_length=120)
    description: Optional[str] = None
    status: Optional[str] = None

class ItemResponse(ItemBase):
    id: int
    created_at: datetime = Field(default_factory=datetime.utcnow)

    class Config:
        from_attributes = True
