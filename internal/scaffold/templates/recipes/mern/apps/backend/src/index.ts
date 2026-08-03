import express from 'express'
import cors from 'cors'
import mongoose from 'mongoose'
import jwt from 'jsonwebtoken'
import bcrypt from 'bcryptjs'

const app = express()
const port = process.env.PORT || 8080

app.use(cors())
app.use(express.json())

// Health check endpoint
app.get('/api/health', (req, res) => {
  res.json({ status: 'ok', stack: 'MERN' })
})

// Mongoose connection mock
const mongoUri = process.env.MONGO_URI || 'mongodb://localhost:27017/[[.ProjectName]]'
mongoose.connect(mongoUri)
  .then(() => console.log('MongoDB connected'))
  .catch(err => console.error('MongoDB connection error:', err))

// CRUD endpoints mock (GET, POST, PUT, DELETE)
let items = [{ id: 1, name: 'Item 1' }]

app.get('/api/items', (req, res) => {
  res.json(items)
})

app.post('/api/items', (req, res) => {
  const newItem = { id: items.length + 1, name: req.body.name }
  items.push(newItem)
  res.status(201).json(newItem)
})

app.put('/api/items/:id', (req, res) => {
  const id = parseInt(req.params.id)
  const item = items.find(i => i.id === id)
  if (item) {
    item.name = req.body.name
    res.json(item)
  } else {
    res.status(404).json({ error: 'Not found' })
  }
})

app.delete('/api/items/:id', (req, res) => {
  const id = parseInt(req.params.id)
  items = items.filter(i => i.id !== id)
  res.json({ message: 'Deleted' })
})

app.listen(port, () => {
  console.log(`MERN Backend running on port ${port}`)
})