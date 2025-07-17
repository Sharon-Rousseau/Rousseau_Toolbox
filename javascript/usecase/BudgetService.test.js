const test = require('node:test');
const assert = require('assert');
const BudgetService = require('./BudgetService');

class MemoryRepo {
  constructor() {
    this.data = [];
  }
  async create(budget) {
    this.data.push({ ...budget, id: this.data.length + 1 });
    return budget;
  }
  async list() {
    return this.data;
  }
}

test('BudgetService creates and lists budgets', async () => {
  const repo = new MemoryRepo();
  const svc = new BudgetService(repo);
  await svc.createBudget({ category: 'Groceries', amount: 100 });
  const list = await svc.listBudgets();
  assert.strictEqual(list.length, 1);
  assert.strictEqual(list[0].category, 'Groceries');
});
