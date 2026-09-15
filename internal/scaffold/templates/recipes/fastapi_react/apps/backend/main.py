from fastapi import FastAPI, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from typing import List
from datetime import datetime
from app.schemas.item import ItemCreate, ItemUpdate, ItemResponse

app = FastAPI(
    title="[[.ProjectName]] API",
    description="FastAPI Backend con validación Pydantic v2 y arquitectura modular",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# In-memory store (or connected with SQLAlchemy / PostgreSQL)
items_db: List[dict] = [
    {
        "id": 1,
        "title": "Configurar FastAPI con Pydantic v2 en Koko CLI",
        "description": "Explorar documentación automática en /docs y /redoc",
        "status": "completed",
        "created_at": datetime.utcnow()
    },
    {
        "id": 2,
        "title": "Conectar frontend React con API FastAPI",
        "description": "Utilizar cliente Axios / Fetch con tipado estricto",
        "status": "in_progress",
        "created_at": datetime.utcnow()
    }
]

@app.get("/api/health", tags=["Health"])
def health_check():
    return {
        "status": "ok",
        "backend": "Python / FastAPI",
        "version": "1.0.0",
        "timestamp": datetime.utcnow().isoformat()
    }

@app.get("/api/items", response_model=List[ItemResponse], tags=["Items"])
def get_items():
    return items_db

@app.post("/api/items", response_model=ItemResponse, status_code=status.HTTP_201_CREATED, tags=["Items"])
def create_item(item_in: ItemCreate):
    new_id = max([i["id"] for i in items_db], default=0) + 1
    new_item = {
        "id": new_id,
        "title": item_in.title,
        "description": item_in.description,
        "status": item_in.status,
        "created_at": datetime.utcnow()
    }
    items_db.append(new_item)
    return new_item

@app.put("/api/items/{item_id}", response_model=ItemResponse, tags=["Items"])
def update_item(item_id: int, item_in: ItemUpdate):
    for item in items_db:
        if item["id"] == item_id:
            if item_in.title is not None:
                item["title"] = item_in.title
            if item_in.description is not None:
                item["description"] = item_in.description
            if item_in.status is not None:
                item["status"] = item_in.status
            return item
    raise HTTPException(status_code=404, detail="Item not found")

@app.delete("/api/items/{item_id}", status_code=status.HTTP_204_NO_CONTENT, tags=["Items"])
def delete_item(item_id: int):
    global items_db
    initial_len = len(items_db)
    items_db = [i for i in items_db if i["id"] != item_id]
    if len(items_db) == initial_len:
        raise HTTPException(status_code=404, detail="Item not found")
    return None