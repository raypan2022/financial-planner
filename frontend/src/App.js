import { useCallback, useEffect, useState } from 'react';
import { Link, Outlet, useNavigate } from 'react-router-dom';
import Alert from './components/Alert';

function App() {
  const [jwtToken, setJwtToken] = useState('');
  const [alertMessage, setAlertMessage] = useState('');
  const [alertClassName, setAlertClassName] = useState('d-none');

  const [tickInterval, setTickInterval] = useState();

  const navigate = useNavigate();

  const logOut = () => {
    const requestOptions = {
      method: 'GET',
      credentials: 'include',
    };

    fetch(`/logout`, requestOptions)
      .catch((error) => {
        console.log('error logging out', error);
      })
      .finally(() => {
        setJwtToken('');
        toggleRefresh(false);
      });

    navigate('/login');
  };

  const toggleRefresh = useCallback(
    (status) => {
      console.log('clicked');

      if (status) {
        console.log('turning on ticking');
        let i = setInterval(() => {
          const requestOptions = {
            method: 'GET',
            credentials: 'include',
          };

          fetch(`/refresh`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
              if (data.access_token) {
                setJwtToken(data.access_token);
              }
            })
            .catch((error) => {
              console.log('user is not logged in');
            });
        }, 600000);
        setTickInterval(i);
        console.log('setting tick interval to', i);
      } else {
        console.log('turning off ticking');
        console.log('turning off tickInterval', tickInterval);
        setTickInterval(null);
        clearInterval(tickInterval);
      }
    },
    [tickInterval]
  );

  useEffect(() => {
    if (jwtToken === '') {
      const requestOptions = {
        method: 'GET',
        credentials: 'include',
      };

      fetch(`/refresh`, requestOptions)
        .then((response) => response.json())
        .then((data) => {
          if (data.access_token) {
            setJwtToken(data.access_token);
            toggleRefresh(true);
          }
        })
        .catch((error) => {
          console.log('user is not logged in', error);
        });
    }
  }, [jwtToken, toggleRefresh]);

  return (
    <div className="container">
      <div className="row">
        <div className="col">
        <h1 className="mt-3">Financial Planner</h1>
        </div>
        <div className="col text-end">
          {jwtToken === '' ? (
            <div className=''>
              <Link to="/login">
                <span className="badge bg-success me-1">Login</span>
              </Link>
              <Link to="/signup">
                <span className="badge bg-success">Sign Up</span>
              </Link>
            </div>
          ) : (
            <a href="#!" onClick={logOut}>
              <span className="badge bg-danger">Logout</span>
            </a>
          )}
        </div>
        <hr className="mb-3"></hr>
      </div>

      <div className="row">
        <div className="col-md-2">
          <nav>
            <div className="list-group">
              <Link to="/" className="list-group-item list-group-item-action">
                Home
              </Link>
              <Link
                to="/income"
                className="list-group-item list-group-item-action"
              >
                Income
              </Link>
              <Link
                to="/expenses"
                className="list-group-item list-group-item-action"
              >
                Expenses
              </Link>
              <Link
                to="/budget"
                className="list-group-item list-group-item-action"
              >
                Set a Budget
              </Link>
              <Link
                to="/goals"
                className="list-group-item list-group-item-action"
              >
                Set a Goal
              </Link>
            </div>
          </nav>
        </div>
        <div className="col-md-10">
          <Alert message={alertMessage} className={alertClassName} />
          <Outlet
            context={{
              jwtToken,
              setJwtToken,
              setAlertClassName,
              setAlertMessage,
              toggleRefresh,
            }}
          />
        </div>
      </div>
    </div>
  );
}

export default App;
