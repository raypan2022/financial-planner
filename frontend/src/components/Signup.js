import { useState } from 'react';
import { useNavigate, useOutletContext } from 'react-router-dom';
import Input from './form/Input';

const Signup = () => {
  const [email, setEmail] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [password, setPassword] = useState('');

  const { setJwtToken } = useOutletContext();
  const { setAlertClassName } = useOutletContext();
  const { setAlertMessage } = useOutletContext();
  const { toggleRefresh } = useOutletContext();

  const navigate = useNavigate();

  const handleSubmit = (event) => {
    event.preventDefault();

    // build the request payload
    let payload = {
      email: email,
      first_name: firstName,
      last_name: lastName,
      password: password,
    };

    const requestOptions = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    };

    fetch(`/signup`, requestOptions)
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          setAlertClassName('alert-danger');
          setAlertMessage(data.message);
        } else {
          setJwtToken(data.access_token);
          setAlertClassName('d-none');
          setAlertMessage('');
          toggleRefresh(true);
          navigate('/');
        }
      })
      .catch((error) => {
        setAlertClassName('alert-danger');
        setAlertMessage(error);
      });
  };

  return (
    <div className="col-md-6 offset-md-3">
      <h2>Sign Up</h2>
      <hr />

      <form onSubmit={handleSubmit}>
        <Input
          title="Email Address"
          type="email"
          className="form-control"
          name="email"
          autoComplete="email-new"
          onChange={(event) => setEmail(event.target.value)}
        />

        <Input
          title="First Name"
          type="text"
          className="form-control"
          name="firstName"
          autoComplete="given-name"
          onChange={(event) => setFirstName(event.target.value)}
        />

        <Input
          title="Last Name"
          type="text"
          className="form-control"
          name="lastName"
          autoComplete="family-name"
          onChange={(event) => setLastName(event.target.value)}
        />

        <Input
          title="Password"
          type="password"
          className="form-control"
          name="password"
          autoComplete="password-new"
          onChange={(event) => setPassword(event.target.value)}
        />

        <hr />

        <input type="submit" className="btn btn-primary" value="Login" />
      </form>
    </div>
  );
};

export default Signup;
