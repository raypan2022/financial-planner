import { useEffect, useState } from 'react';
import { useOutletContext } from 'react-router-dom';
import { FaPlus } from 'react-icons/fa';

import ExpenseForm from './ExpenseForm';

const Expenses = () => {
  const { jwtToken } = useOutletContext();

  const [expenses, setExpenses] = useState([]);
  const [showModal, setShowModal] = useState(false);

  const handleShowModal = () => setShowModal(true);
  const handleCloseModal = () => setShowModal(false);

  useEffect( () => {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", "Bearer " + jwtToken);

    const requestOptions = {
      method: "GET",
      headers: headers,
      credentials: "include",
    }

    fetch(`/admin/expenses`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        setExpenses(data);
      })
      .catch(err => {
        console.log(err);
      })

  }, [jwtToken]);

  return (
    <div className="container">
      <div className="row">
        <div className="col">
          <h2>Expenses</h2>
        </div>
        <div className="col text-end mt-3">
        <a
            href="#!"
            className="text-decoration-none d-inline-flex align-items-center"
            onClick={handleShowModal}
          >
            <FaPlus className="me-2" />
            <span>Add an expense</span>
          </a>
        </div>
        <hr />
      </div>

      <table className="table table-hover">
        <thead>
          <tr>
            <th>Category</th>
            <th>Date</th>
            <th>Amount</th>
          </tr>
        </thead>
        <tbody>
          {expenses && expenses.map((e) => (
            <tr key={e.id}>
              <td>
                {e.category.name}
              </td>
              <td>{e.date?.split('T')[0]}</td>
              <td>{e.amount}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <ExpenseForm show={showModal} handleClose={handleCloseModal} />
    </div>
  );
};

export default Expenses;
