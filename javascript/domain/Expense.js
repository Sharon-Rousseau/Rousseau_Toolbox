class Expense {
  constructor({ id = null, date, amount, vendor, category }) {
    this.id = id;
    this.date = date;
    this.amount = amount;
    this.vendor = vendor;
    this.category = category;
  }
}
module.exports = Expense;
