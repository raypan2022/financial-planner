import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import App from './App';
import EditMovie from './components/EditMovie';
import ErrorPage from './components/ErrorPage';
import Genres from './components/Genres';
import GraphQL from './components/GraphQL';
import Home from './components/Home';
import Login from './components/Login';
import ManageCatalogue from './components/ManageCatalogue';
import Movies from './components/Movies';
import Movie from './components/Movie';
import OneGenre from './components/OneGenre';
import Income from './components/Income';
import Expenses from './components/Expenses';
import Goal from './components/Goal';
import Budget from './components/Budget';
import Signup from './components/Signup';

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <ErrorPage />,
    children: [
      {index: true, element: <Home /> },
      {
        path: "/income",
        element: <Income />,
      },
      // {
      //   path: "/income/:id",
      //   element: <Movie />,
      // },
      {
        path: "/expenses",
        element: <Expenses />,
      },
      // {
      //   path: "/expenses/:id",
      //   element: <OneGenre />
      // },
      {
        path: "/budget",
        element: <Budget />,
      },
      // {
      //   path: "/budget/:id",
      //   element: <EditMovie />,
      // },
      {
        path: "/goals",
        element: <Goal />,
      },
      // {
      //   path: "/goals/:id",
      //   element: <EditMovie />,
      // },
      {
        path: "/login",
        element: <Login />,
      },
      {
        path: "/signup",
        element: <Signup />,
      }
    ]
  }
])

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
