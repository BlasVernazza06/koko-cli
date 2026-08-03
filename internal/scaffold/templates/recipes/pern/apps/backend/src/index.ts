import express from 'express'
import cors from 'cors'
import { PrismaClient } from '@prisma/client'

const app = express()
const prisma = new PrismaClient()
const port = process.env.PORT || 8080

app.use(cors())
app.use(express.json())

app.get('/api/health', (req, res) => {
  res.json({ status: 'ok', stack: 'PERN' })
})

// Mocks for CRUD using Prisma
app.get('/api/users', async (req, res) => {
  try {
    const users = await prisma.user.findMany()
    res.json(users)
  } catch (err) {
    res.status(500).json({ error: 'Database error' })
  }
})

app.post('/api/users', async (req, res) => {
  try {
    const user = await prisma.user.create({
      data: { email: req.body.email, name: req.body.name }
    })
    res.status(201).json(user)
  } catch (err) {
    res.status(400).json({ error: 'Failed to create user' })
  }
})

app.listen(port, () => {
  console.log(`PERN Backend running on port ${port}`)
})