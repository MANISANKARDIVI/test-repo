// src/index.js
const express = require('express');
const app = express();
const port = process.env.PORT || 3000; // Set port to 3000 or use PORT env variable

app.get('/', (req, res) => {
  res.send("Hello, GitHub Actions CI!");
});

app.listen(port, () => {
  console.log(`Starting server on port ${port}...`);
});
