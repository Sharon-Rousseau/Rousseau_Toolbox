const express = require('express');
const { MongoClient, ObjectId } = require('mongodb');
const BudgetService = require('../../usecase/BudgetService');

class MongoBudgetRepo {
  constructor(db) {
    this.collection = db.collection('budgets');
  }

  async create(budget) {
    const res = await this.collection.insertOne(budget);
    return { ...budget, id: res.insertedId };
  }

  async list() {
    const docs = await this.collection.find().toArray();
    return docs.map(d => ({
      id: d._id,
      category: d.category,
      description: d.description,
      amount: d.amount,
    }));
  }
}

async function createServer(mongoUri) {
  const client = new MongoClient(mongoUri);
  await client.connect();
  const db = client.db();
  const repo = new MongoBudgetRepo(db);
  const svc = new BudgetService(repo);

  const app = express();
  app.use(express.json());

  app.get('/budgets', async (req, res) => {
    const budgets = await svc.listBudgets();
    res.json(budgets);
  });

  app.post('/budgets', async (req, res) => {
    const created = await svc.createBudget(req.body);
    res.status(201).json(created);
  });

  return { app, client };
}

module.exports = { createServer };
