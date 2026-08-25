from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List, Optional

router = APIRouter()

class Todo(BaseModel):
    id: int
    title: string = ""
    done: bool = False

class TodoCreate(BaseModel):
    title: str

class TodoUpdate(BaseModel):
    title: Optional[str] = None
    done: Optional[bool] = None

todos_db: List[dict] = [
    {"id": 1, "title": "Learn FastAPI & Pydantic v2", "done": False},
    {"id": 2, "title": "Configure Koko Monorepo", "done": True},
]

@router.get("/todos")
def get_todos():
    return todos_db

@router.post("/todos", status_code=201)
def create_todo(payload: TodoCreate):
    new_todo = {
        "id": len(todos_db) + 1,
        "title": payload.title,
        "done": False,
    }
    todos_db.append(new_todo)
    return new_todo

@router.put("/todos/{todo_id}")
def update_todo(todo_id: int, payload: TodoUpdate):
    for todo in todos_db:
        if todo["id"] == todo_id:
            if payload.title is not None:
                todo["title"] = payload.title
            if payload.done is not None:
                todo["done"] = payload.done
            return todo
    raise HTTPException(status_code=404, detail="Todo not found")

@router.delete("/todos/{todo_id}")
def delete_todo(todo_id: int):
    global todos_db
    initial_len = len(todos_db)
    todos_db = [t for t in todos_db if t["id"] != todo_id]
    if len(todos_db) == initial_len:
        raise HTTPException(status_code=404, detail="Todo not found")
    return {"message": "Todo deleted successfully"}
