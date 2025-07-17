class BudgetService {
  constructor(repo) {
    this.repo = repo; // repo must implement create, list methods
  }

  async createBudget(data) {
    return this.repo.create(data);
  }

  async listBudgets() {
    return this.repo.list();
  }
}

module.exports = BudgetService;
