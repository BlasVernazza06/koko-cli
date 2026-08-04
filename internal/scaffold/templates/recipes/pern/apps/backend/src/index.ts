import express from 'express';
import cors from 'cors';
import authRouter from './controllers/auth.controller';
import taskRouter from './controllers/task.controller';

const app = express();
const port = process.env.PORT || 8080;

app.use(cors());
app.use(express.json());

// Routes
app.use('/api/auth', authRouter);
app.use('/api/tasks', taskRouter);

app.get('/api/health', (req, res) => {
  res.json({ status: 'ok', stack: 'PERN' });
});

app.listen(port, () => {
  console.log(`PERN Backend running on port ${port}`);
});