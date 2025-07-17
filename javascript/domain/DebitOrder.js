class DebitOrder {
  constructor({ id = null, name, amount, dueDate, paid = false }) {
    this.id = id;
    this.name = name;
    this.amount = amount;
    this.dueDate = dueDate;
    this.paid = paid;
  }
}
module.exports = DebitOrder;
