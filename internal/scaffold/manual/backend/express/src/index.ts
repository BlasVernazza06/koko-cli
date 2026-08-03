import express from 'express';
import cors from 'cors';

const app = express();
const port = process.env.PORT || 8080;

app.use(cors());
app.use(express.json());

interface Todo {
  id: number;
  title: string;
  done: boolean;
}

let todos: Todo[] = [
  { id: 1, title: 'Aprender Node.js Express', done: false },
  { id: 2, title: 'Configurar Monorepo', done: true },
];

// Health Check
app.get('/api/health', (req, res) => {
  res.json({
    status: 'ok',
    project: '[[.ProjectName]]-backend',
  });
});

// 1. GET - Obtener todos los todos
app.get('/api/todos', (req, res) => {
  res.json(todos);
});

// 2. POST - Crear un nuevo todo
app.post('/api/todos', (req, res) => {
  const { title } = req.body;
  if (!title) {
    res.status(400).json({ error: 'Title is required' });
    return;
  }

  const newTodo: Todo = {
    id: todos.length + 1,
    title,
    done: false,
  };
  todos.push(newTodo);
  res.status(201).json(newTodo);
});

// 3. PUT - Actualizar un todo existente
app.put('/api/todos/:id', (req, res) => {
  const id = parseInt(req.params.id, 10);
  const { title, done } = req.body;

  const todoIndex = todos.findIndex((t) => t.id === id);
  if (todoIndex === -1) {
    res.status(404).json({ error: 'Todo not found' });
    return;
  }

  if (title !== undefined) todos[todoIndex].title = title;
  if (done !== undefined) todos[todoIndex].done = done;

  res.json(todos[todoIndex]);
});

// 4. DELETE - Eliminar un todo
app.delete('/api/todos/:id', (req, res) => {
  const id = parseInt(req.params.id, 10);
  const todoIndex = todos.findIndex((t) => t.id === id);
  if (todoIndex === -1) {
    res.status(404).json({ error: 'Todo not found' });
    return;
  }

  todos = todos.filter((t) => t.id !== id);
  res.json({ message: 'Todo deleted successfully' });
});

app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});
