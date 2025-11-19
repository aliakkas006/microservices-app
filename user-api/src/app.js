const express = require('express');
const cors = require('cors');
const userRoutes = require('./routes/userRoutes');

const app = express();

app.use(cors());
app.use(express.json());

// API routes
app.use('/api/users', userRoutes);

// Add a simple health check endpoint
app.get('/health', (req, res) => {
  res.status(200).json({ status: 'User service is UP' });
});

module.exports = app;
