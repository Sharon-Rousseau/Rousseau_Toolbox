class Budget {
  constructor({ id = null, category, description = '', amount = 0 }) {
    this.id = id;
    this.category = category;
    this.description = description;
    this.amount = amount;
  }
}
module.exports = Budget;
