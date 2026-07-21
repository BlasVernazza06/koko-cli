import { serve } from '@hono/node-server'
import { Hono } from 'hono'
import { cors } from 'hono/cors'

const app = new Hono()

app.use('/api/*', cors())

interface Todo {
  id: number
  title: string
  done: boolean
}

let todos: Todo[] = [
  { id: 1, title: 'Aprender Hono', done: false },
  { id: 2, title: 'Configurar Monorepo', done: true }
]

// Health Check
app.get('/api/health', (c) => {
  return c.json({
    status: 'ok',
    project: '[[.ProjectName]]-backend'
  })
})

// 1. GET - Obtener todos los todos
app.get('/api/todos', (c) => {
  return c.json(todos)
})

// 2. POST - Crear un nuevo todo
app.post('/api/todos', async (c) => {
  const { title } = await c.req.json()
  if (!title) {
    return c.json({ error: 'Title is required' }, 400)
  }

  const newTodo: Todo = {
    id: todos.length + 1,
    title,
    done: false
  }
  todos.push(newTodo)
  return c.json(newTodo, 201)
})

// 3. PUT - Actualizar un todo existente
app.put('/api/todos/:id', async (c) => {
  const id = parseInt(c.req.param('id'), 10)
  const { title, done } = await c.req.json()

  const todoIndex = todos.findIndex((t) => t.id === id)
  if (todoIndex === -1) {
    return c.json({ error: 'Todo not found' }, 404)
  }

  if (title !== undefined) todos[todoIndex].title = title
  if (done !== undefined) todos[todoIndex].done = done

  return c.json(todos[todoIndex])
})

// 4. DELETE - Eliminar un todo
app.delete('/api/todos/:id', (c) => {
  const id = parseInt(c.req.param('id'), 10)
  const todoIndex = todos.findIndex((t) => t.id === id)
  if (todoIndex === -1) {
    return c.json({ error: 'Todo not found' }, 404)
  }

  todos = todos.filter((t) => t.id !== id)
  return c.json({ message: 'Todo deleted successfully' })
})

const port = 3000
console.log(`Server is running on port ${port}`)

serve({
  fetch: app.fetch,
  port
})
