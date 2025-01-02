import React, { useState, useEffect } from 'react';
import { useOutletContext } from 'react-router-dom';
import { Modal, Button } from 'react-bootstrap';

import Input from './form/Input';
import TextArea from './form/TextArea';
import ReactCreatable from './form/ReactCreatable';
import ReactSelectable from './form/ReactSelectable';

const ExpenseForm = ({ show, handleClose }) => {
  const { jwtToken } = useOutletContext();

  // paycheque object
  const [expense, setExpense] = useState({
    amount: '',
    category: '',
    payment_method: '',
    date: '',
    description: '',
  });

  // react-select options
  const [categoryOptions, setCategoryOptions] = useState([]);
  const paymentOptions = [{ value: 'Cash', label: 'Cash' }, { value: 'Cheque', label: 'Cheque' }, { value: 'Credit Card', label: 'Credit Card' }, { value: 'Debit Card', label: 'Debit Card' }, { value: 'E-Transfer', label: 'E-Transfer' }];

  // state to toggle preview visibility
  const [showPreview, setShowPreview] = useState(false);

  useEffect(() => {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", "Bearer " + jwtToken);

    const requestOptions = {
      method: "GET",
      headers: headers,
      credentials: "include",
    }

    fetch(`/admin/categories`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        data?.forEach((category) => {
          setCategoryOptions((prevOptions) => [...prevOptions, { value: category.name, label: category.name }]);
        })
      })
      .catch(err => {
        console.log(err);
      })

  }, [jwtToken]);

  const handleChange = (event, value) => {
    let name = event.target.name;
    setExpense({
      ...expense,
      [name]: value,
    });
  };

  const onBlur = (event, amount) => {
    const value = parseFloat(amount);

    if (isNaN(value)) {
      handleChange(event, '');
      return;
    }

    handleChange(event, value.toFixed(2));
  };

  const handleOptionChange = (newValue, fieldName) => {
    setExpense({
      ...expense,
      [fieldName]: newValue ? newValue.value : null,
    });
  };

  const handleOptionCreate = (inputValue) => {
    const newOption = { value: inputValue, label: inputValue };
    setCategoryOptions((prevOptions) => [...prevOptions, newOption]);
    handleOptionChange(newOption, 'category');
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", "Bearer " + jwtToken);

    let requestBody = expense;

    requestBody.category = {
      name: expense.category,
    };
    requestBody.date = new Date(expense.date);
    requestBody.amount = parseFloat(expense.amount);

    const requestOptions = {
      body: JSON.stringify(requestBody),
      method: "POST",
      headers: headers,
      credentials: "include",
    };

    fetch(`/admin/expenses/new`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          console.log(data.error);
        } else {
          handleClose();
          window.location.reload();
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }

  return (
    <>
      {show && (
        <Modal show={show} onHide={handleClose}>
          <Modal.Header closeButton>
            <Modal.Title>Expense</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <form>
              <ReactCreatable
                mandatory={true}
                title="Category"
                name="category"
                value={{ value: expense.category, label: expense.category }}
                onChange={(newValue) => handleOptionChange(newValue, "category")}
                onCreateOption={handleOptionCreate}
                options={categoryOptions}
              />
              <ReactSelectable
                mandatory={true}
                title="Payment Method"
                name="payment_method"
                value={{ value: expense.payment_method, label: expense.payment_method }}
                onChange={(newValue) => handleOptionChange(newValue, "payment_method")}
                options={paymentOptions}
              />
              <Input
                mandatory={true}
                title="Amount"
                dollar={true}
                className="form-control"
                type="text"
                name="amount"
                value={expense.amount}
                onBlur={(event) => onBlur(event, expense.amount)}
                onChange={(event) => handleChange(event, event.target.value)}
              />
              <Input
                mandatory={true}
                title="Date of Transaction"
                className="form-control"
                type="date"
                name="date"
                value={expense.date}
                onChange={(event) => handleChange(event, event.target.value)}
              />
              <TextArea
                mandatory={false}
                title="Notes"
                name="description"
                value={expense.description}
                onChange={(event) => handleChange(event, event.target.value)}
              />
            </form>
            {showPreview && (
              <div className="mt-4">
                <p className="">
                  You spent{' '}
                  {expense.amount ? expense.amount : '0.00'} for{' '}
                  {expense.category ? expense.category : 'specified company'} on{' '}
                  {expense.date ? expense.date : 'specified date'}.
                </p>
              </div>
            )}
            <button
              type="button"
              className="btn btn-link mt-2"
              onClick={() => setShowPreview(!showPreview)}
            >
              {showPreview ? 'Hide Preview' : 'Show Preview'}
            </button>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleClose}>
              Close
            </Button>
            <Button variant="primary" onClick={handleSubmit}>
              Save changes
            </Button>
          </Modal.Footer>
        </Modal>
      )}
    </>
  );
};

export default ExpenseForm;