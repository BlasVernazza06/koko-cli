from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="[[.ProjectName]] API", version="0.1.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/api/health")
def health_check():
    return {"status": "ok", "backend": "FastAPI"}

# Mock endpoints for CRUD (GET, POST, PUT, DELETE)
items = [{"id": 1, "name": "Item 1"}]

@app.get("/api/items")
def read_items():
    return items

@app.post("/api/items")
def create_item(name: str):
    new_item = {"id": len(items) + 1, "name": name}
    items.append(new_item)
    return new_item

@app.put("/api/items/{item_id}")
def update_item(item_id: int, name: str):
    for item in items:
        if item["id"] == item_id:
            item["name"] = name
            return item
    return {"error": "Not found"}

@app.delete("/api/items/{item_id}")
def delete_item(item_id: int):
    global items
    items = [item for item in items if item["id"] != item_id]
    return {"message": "Deleted"}