import { Link, useOutletContext } from 'react-router-dom';
import { useEffect, useState } from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell
} from 'recharts';

const COLORS = [
  '#0088FE', // Blue
  '#00C49F', // Green
  '#FFBB28', // Yellow
  '#FF8042', // Orange
  '#FF6384', // Pink
  '#36A2EB', // Light Blue
  '#FFCE56', // Light Yellow
  '#4BC0C0', // Teal
  '#9966FF', // Purple
  '#FF9F40', // Light Orange
  '#C9CBCF', // Grey
  '#8A2BE2', // Blue Violet
  '#DC143C', // Crimson
  '#00FF7F', // Spring Green
  '#FF4500'  // Orange Red
];

const Home = () => {
  const { jwtToken } = useOutletContext();
  const [summary, setSummary] = useState({});
  const [incomeExpenseLineData, setIncomeExpenseLineData] = useState([]);
  const [sourcePie, setSourcePie] = useState([]);
  const [categoryPie, setCategoryPie] = useState([]);

  const getLineData = (months) => {
    return months.map((m) => ({
      name: m.month,
      Income: m.income_sum,
      Expense: m.expense_sum,
    }));
  };

  const getPieData = (distribution) => {
    return Object.keys(distribution).map((key) => ({
      name: key,
      value: distribution[key],
    }));
  };

  useEffect(() => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('Authorization', 'Bearer ' + jwtToken);

    const requestOptions = {
      method: 'GET',
      headers: headers,
      credentials: 'include',
    };

    fetch(`/admin/summary`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        console.log('summary', data);
        setSummary(data);
        setIncomeExpenseLineData(getLineData(data.months));
        setSourcePie(getPieData(data.overall_income_by_source));
        setCategoryPie(getPieData(data.overall_expense_by_category));
      })
      .catch((err) => {
        console.log(err);
      });
  }, [jwtToken]);

  return (
    <div>
      <div className="text-center">
        <h2>Dashboard</h2>
        <hr />
      </div>

      <div className="container">
        <div className="row gx-3">
          <div className="col">
            <div className="card">
              <div className="card-body">
                <h4 className="card-title">${summary.account_balance}</h4>
                <p>Account Balance</p>
              </div>
            </div>
          </div>

          <div className="col">
            <div className="card">
              <div className="card-body">
                <h4 className="card-title">${summary.income_sum_total}</h4>
                <p>Total Income</p>
              </div>
            </div>
          </div>

          <div className="col">
            <div className="card">
              <div className="card-body">
                <h4 className="card-title">${summary.expense_sum_total}</h4>
                <p>Total Expense</p>
              </div>
            </div>
          </div>
        </div>

        <div className="row mt-5">
          <h5 className="mb-4">Metrics From The Last 12 Months</h5>
          <ResponsiveContainer width="100%" height={400}>
            <LineChart
              width={500}
              height={300}
              data={incomeExpenseLineData}
              margin={{
                top: 5,
                right: 30,
                left: 20,
                bottom: 5,
              }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" padding={{ left: 20, right: 20 }} />
              <YAxis padding={{ top: 20, bottom: 20 }} />
              <Tooltip />
              <Legend />
              <Line
                type="monotone"
                dataKey="Income"
                stroke="#8884d8"
                strokeWidth={2}
                activeDot={{ r: 8 }}
              />
              <Line
                type="monotone"
                dataKey="Expense"
                stroke="#82ca9d"
                strokeWidth={2}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>

        <div className="row mt-5">
          <div className="col">
            <h5 className="mb-4">Distribution of Total Income</h5>
            <ResponsiveContainer width="100%" height={200}>
              <PieChart width={400} height={400}>
                <Pie
                  dataKey="value"
                  isAnimationActive={false}
                  data={sourcePie}
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  fill="#8884d8"
                  label
                />
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
          <div className="col">
            <h5 className="mb-4">Distribution of Total Expense</h5>
            <ResponsiveContainer width="100%" height={200}>
              <PieChart width={400} height={400}>
                <Pie
                  dataKey="value"
                  isAnimationActive={false}
                  data={categoryPie}
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  fill="#8884d8"
                  label
                >
                  {categoryPie.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* <div className="row mt-5">
          <div className='col'>
            <h5 className="mb-4">Top 3 Income of </h5>
            <table className="table table-hover">
              <thead>
                <tr>
                  <th>From</th>
                  <th>Date</th>
                  <th>Amount</th>
                </tr>
              </thead>
              <tbody>
                {paycheques && paycheques.map((p) => (
                  <tr key={p.id}>
                    <td>
                      {p.source.name}
                    </td>
                    <td>{p.date?.split('T')[0]}</td>
                    <td>{p.amount}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div> */}
      </div>
    </div>
  );
};

export default Home;
