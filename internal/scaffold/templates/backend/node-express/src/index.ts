import express from 'express';
import cors from 'cors';

const app = express();
const port = process.env.PORT || 8080;

app.use(cors());
app.use(express.json());

app.get('/api/health', (req, res) => {
  res.json({
    status: 'ok',
    message: 'Welcome to Node.js Express API!',
  });
});

app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});
