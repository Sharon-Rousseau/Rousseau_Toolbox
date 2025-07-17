import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardContent } from '@/components/ui/card';
import { Table, Thead, Tbody, Tr, Th, Td } from '@/components/ui/table';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';

const defaultCategories = [
  'Rent', 'Groceries', 'Utilities', 'Transport', 'Entertainment',
  'Savings', 'Insurance', 'Debt', 'Miscellaneous'
];

export default function BudgetingApp() {
  const [budgetRows, setBudgetRows] = useState([
    { id: 1, category: 'Groceries', description: '', amount: 0 }
  ]);
  const [expensesData, setExpensesData] = useState([]);

  const addBudgetRow = () => {
    setBudgetRows([...budgetRows, { id: Date.now(), category: '', description: '', amount: 0 }]);
  };

  return (
    <div className="p-6 bg-gray-50 min-h-screen">
      <header className="mb-6">
        <h1 className="text-3xl font-bold">Budgeting App</h1>
      </header>
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Budget Breakdown</h2>
          </CardHeader>
          <CardContent>
            <Table>
              <Thead>
                <Tr>
                  <Th>Category</Th>
                  <Th>Description</Th>
                  <Th>Amount (ZAR)</Th>
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {budgetRows.map(row => (
                  <Tr key={row.id}>
                    <Td>
                      <select
                        className="w-full border rounded p-1"
                        value={row.category}
                        onChange={() => {}}
                      >
                        {defaultCategories.map(cat => (
                          <option key={cat}>{cat}</option>
                        ))}
                        <option value="custom">+ Custom Category</option>
                      </select>
                    </Td>
                    <Td>
                      <input
                        type="text"
                        className="w-full border rounded p-1"
                        placeholder="Description"
                        value={row.description}
                        onChange={() => {}}
                      />
                    </Td>
                    <Td>
                      <input
                        type="number"
                        className="w-full border rounded p-1"
                        placeholder="0"
                        value={row.amount}
                        onChange={() => {}}
                      />
                    </Td>
                    <Td>
                      <Button variant="ghost" onClick={() => {}}>Delete</Button>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
            <div className="mt-4">
              <Button onClick={addBudgetRow}>Add Row</Button>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Spending Graph</h2>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={expensesData}>
                <XAxis dataKey="day" />
                <YAxis />
                <Tooltip />
                <Line type="monotone" dataKey="amount" />
              </LineChart>
            </ResponsiveContainer>
            <p className="mt-4 text-sm text-gray-600">
              Net Income: R26 400 + R69 000; RA Fund: R13 000 (pre-salary)
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
