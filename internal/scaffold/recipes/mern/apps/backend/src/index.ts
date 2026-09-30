import express from 'express';
import cors from 'cors';
import mongoose from 'mongoose';
import dotenv from 'dotenv';
import { Item } from './models/item';
import { createItemSchema, updateItemSchema } from './schemas/item.schema';

dotenv.config();

const app = express();
const port = process.env.PORT || 8080;
const mongoUri = process.env.MONGO_URI || 'mongodb://localhost:27017/[[.ProjectName]]';

app.use(cors());
app.use(express.json());

// MongoDB connection
mongoose
  .connect(mongoUri)
  .then(() => console.log('🍃 [MongoDB] Connected successfully to:', mongoUri))
  .catch((err) => console.warn('⚠️ [MongoDB] Connection warning (running in offline mode):', err.message));

// Health check endpoint
app.get('/api/health', (req, res) => {
  res.json({
    status: 'ok',
    stack: 'MERN (React + Express + MongoDB + Mongoose)',
    database: mongoose.connection.readyState === 1 ? 'connected' : 'disconnected',
    timestamp: new Date().toISOString(),
  });
});

// CRUD endpoints with Mongoose & Zod
app.get('/api/items', async (req, res) => {
  try {
    if (mongoose.connection.readyState === 1) {
      const items = await Item.find().sort({ createdAt: -1 });
      return res.json(items);
    }
    // In-memory fallback if MongoDB is not running locally yet
    return res.json([
      { _id: '1', title: 'Explorar arquitectura MERN en Koko CLI', status: 'completed', createdAt: new Date() },
      { _id: '2', title: 'Conectar base de datos MongoDB con Docker', status: 'in_progress', createdAt: new Date() },
    ]);
  } catch (error) {
    res.status(500).json({ error: 'Error fetching items', details: error });
  }
});

app.post('/api/items', async (req, res) => {
  try {
    const validated = createItemSchema.parse(req.body);
    if (mongoose.connection.readyState === 1) {
      const newItem = await Item.create(validated);
      return res.status(201).json(newItem);
    }
    return res.status(201).json({ _id: Date.now().toString(), ...validated, createdAt: new Date() });
  } catch (error: any) {
    if (error.errors) {
      return res.status(400).json({ error: 'Validation Error', details: error.errors });
    }
    res.status(500).json({ error: 'Failed to create item' });
  }
});

app.put('/api/items/:id', async (req, res) => {
  try {
    const validated = updateItemSchema.parse(req.body);
    if (mongoose.connection.readyState === 1) {
      const updated = await Item.findByIdAndUpdate(req.params.id, validated, { new: true });
      if (!updated) return res.status(404).json({ error: 'Item not found' });
      return res.json(updated);
    }
    return res.json({ _id: req.params.id, ...validated, updatedAt: new Date() });
  } catch (error: any) {
    if (error.errors) {
      return res.status(400).json({ error: 'Validation Error', details: error.errors });
    }
    res.status(500).json({ error: 'Failed to update item' });
  }
});

app.delete('/api/items/:id', async (req, res) => {
  try {
    if (mongoose.connection.readyState === 1) {
      await Item.findByIdAndDelete(req.params.id);
    }
    res.json({ message: 'Item deleted successfully', id: req.params.id });
  } catch (error) {
    res.status(500).json({ error: 'Failed to delete item' });
  }
});

app.listen(port, () => {
  console.log(`🚀 [MERN Backend] Running on http://localhost:${port}`);
});