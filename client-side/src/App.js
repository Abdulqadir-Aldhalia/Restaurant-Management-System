import { store } from "./redux/store.js";
import { Provider } from "react-redux";
import ProtectedRoute from "./pages/ProtectedRoute.jsx";
import LoginAdminPortal from "./pages/login/loginAdmin/LoginAdminPortal.jsx";
import Login from "./pages/login/loginUser/Login.jsx";
import RegisterForm from "./pages/register/RegisterForm.jsx";
import Home from "./pages/home/Home.jsx";
import UserHome from "./pages/users/UserHome.jsx";
import AdminDashboard from "./pages/AdminManagment/AdminDashboard.jsx";

import { createBrowserRouter, RouterProvider } from "react-router-dom";

const router = createBrowserRouter([
  {
    path: "/loginAdminPortal",
    element: <LoginAdminPortal />,
  },
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/userHome",
    element: (
      <ProtectedRoute>
        <UserHome />
      </ProtectedRoute>
    ),
  },
  {
    path: "/register",
    element: <RegisterForm />,
  },
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/adminDashboard",
    element: (
      <ProtectedRoute>
        <AdminDashboard />
      </ProtectedRoute>
    ),
  },
]);

function App() {
  return (
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  );
}

export default App;
