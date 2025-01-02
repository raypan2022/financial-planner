import React, { useState, useEffect } from 'react';
import { useOutletContext } from 'react-router-dom';
import { Modal, Button } from 'react-bootstrap';

import Input from './form/Input';
import TextArea from './form/TextArea';
import ReactCreatable from './form/ReactCreatable';

const PaychequeForm = ({ show, handleClose }) => {
  const { jwtToken } = useOutletContext();

  // paycheque object
  const [paycheque, setPaycheque] = useState({
    amount: '',
    source: '',
    date: '',
    description: '',
  });

  // react-select options
  const [sourceOptions, setSourceOptions] = useState([]);

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

    fetch(`/admin/sources`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        data?.forEach((source) => {
          setSourceOptions((prevOptions) => [...prevOptions, { value: source.name, label: source.name }]);
        })
      })
      .catch(err => {
        console.log(err);
      })

  }, [jwtToken]);

  const handleChange = (event, value) => {
    let name = event.target.name;
    setPaycheque({
      ...paycheque,
      [name]: value,
    });
  };

  const onBlur = (event, price) => {
    const value = parseFloat(price);

    if (isNaN(value)) {
      handleChange(event, '');
      return;
    }

    handleChange(event, value.toFixed(2));
  };

  const handleOptionChange = (newValue) => {
    setPaycheque({
      ...paycheque,
      source: newValue ? newValue.value : null,
    });
  };

  const handleOptionCreate = (inputValue) => {
    const newOption = { value: inputValue, label: inputValue };
    setSourceOptions((prevOptions) => [...prevOptions, newOption]);
    handleOptionChange(newOption);
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Authorization", "Bearer " + jwtToken);

    let requestBody = paycheque;

    requestBody.source = {
      name: paycheque.source,
    };
    requestBody.date = new Date(paycheque.date);
    requestBody.amount = parseFloat(paycheque.amount);

    const requestOptions = {
      body: JSON.stringify(requestBody),
      method: "POST",
      headers: headers,
      credentials: "include",
    };

    fetch(`/admin/incomes/new`, requestOptions)
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
            <Modal.Title>Paycheque</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <form>
              <ReactCreatable
                mandatory={true}
                title="From"
                name="source"
                value={{ value: paycheque.source, label: paycheque.source }}
                onChange={handleOptionChange}
                onCreateOption={handleOptionCreate}
                options={sourceOptions}
              />
              <Input
                mandatory={true}
                title="Amount"
                dollar={true}
                className="form-control"
                type="text"
                name="amount"
                value={paycheque.amount}
                onBlur={(event) => onBlur(event, paycheque.amount)}
                onChange={(event) => handleChange(event, event.target.value)}
              />
              <Input
                mandatory={true}
                title="Date Received"
                className="form-control"
                type="date"
                name="date"
                value={paycheque.date}
                onChange={(event) => handleChange(event, event.target.value)}
              />
              <TextArea
                mandatory={false}
                title="Notes"
                name="description"
                value={paycheque.description}
                onChange={(event) => handleChange(event, event.target.value)}
              />
            </form>
            {showPreview && (
              <div className="mt-4">
                <p className="">
                  You received{' '}
                  {paycheque.amount ? paycheque.amount : '0.00'} from{' '}
                  {paycheque.source ? paycheque.source : 'specified company'} on{' '}
                  {paycheque.date ? paycheque.date : 'specified date'}.
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

export default PaychequeForm;