const express = require('express');
const { getUsers, getUserById, createUser } = require('../controllers/userController');

const router = express.Router();

router.get('/', getUsers);             // GET /api/users
router.get('/:id', getUserById);      // GET /api/users/:id
router.post('/', createUser);         // POST /api/users

module.exports = router;
