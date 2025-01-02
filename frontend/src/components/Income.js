import { useEffect, useState } from 'react';
import { useOutletContext } from 'react-router-dom';
import { FaPlus } from 'react-icons/fa';

import PaychequeForm from './PaychequeForm';

const Income = () => {
  const { jwtToken } = useOutletContext();

  const [paycheques, setPaycheques] = useState([]);
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

    fetch(`/admin/incomes`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        setPaycheques(data);
      })
      .catch(err => {
        console.log(err);
      })

  }, [jwtToken]);

  return (
    <div className="container">
      <div className="row">
        <div className="col">
          <h2>Income</h2>
        </div>
        <div className="col text-end mt-3">
          <a
            href="#!"
            className="text-decoration-none d-inline-flex align-items-center"
            onClick={handleShowModal}
          >
            <FaPlus className="me-2" />
            <span>Add a paycheque</span>
          </a>
        </div>
        <hr />
      </div>

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

      <PaychequeForm show={showModal} handleClose={handleCloseModal} />
    </div>
  );
};

export default Income;
