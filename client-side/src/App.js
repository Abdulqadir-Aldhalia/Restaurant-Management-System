import { store } from "./redux/store.js";
import { Provider } from "react-redux";
import ProtectedRoute from "./pages/ProtectedRoute.jsx";

import LoginAdminPortal from "./pages/login/loginAdmin/LoginAdminPortal.jsx";
import Login from "./pages/login/loginUser/Login.jsx";
import LoginVendorPortal from "./pages/login/loginVendor/LoginVendorPortal.jsx";

import RegisterForm from "./pages/register/RegisterForm.jsx";
import Home from "./pages/home/Home.jsx";
import UserHome from "./pages/users/UserHome.jsx";
import AdminDashboard from "./pages/AdminManagment/AdminDashboard.jsx";
import VendorManagement from "./pages/vendorsManagement/VendorManagement.jsx";

import { createBrowserRouter, RouterProvider } from "react-router-dom";
import VendorOrderManagement from "./pages/vendorsManagement/VendorOrderManagement.jsx";
import VendorAdminManagement from "./pages/vendorsManagement/VendorAdminManagement.jsx";
import VendorItemsManagement from "./pages/vendorsManagement/VendorItemsManagement.jsx";
import VendorProfileManagement from "./pages/vendorsManagement/VendorProfileManegement.jsx";

import VendorApp from "./pages/vendorsManagement/VendorApp";
import VendorAdminProfile from "./pages/vendorsManagement/VendorAdminProfile";
import VendorHome from "./pages/vendorsManagement/VendorHome";
import VendorTableManagement from "./pages/vendorsManagement/VendorTableManagement.jsx";

const router = createBrowserRouter([
  {
    path: "/loginAdminPortal",
    element: <LoginAdminPortal />,
  },
  {
    path: "/loginVendorPortal",
    element: <LoginVendorPortal />,
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
  {
    path: "/vendorApp",
    element: (
      <ProtectedRoute>
        <VendorApp />,
      </ProtectedRoute>
    ),
    children: [
      {
        path: "VendorHome",
        element: <VendorHome />,
      },
      {
        path: "VendorAdminProfile",
        element: <VendorAdminProfile />,
      },
      {
        path: "",
        element: <VendorHome />,
      },
    ],
  },

  {
    path: "/vendorManagement/:vendorId",
    element: (
      <ProtectedRoute>
        <VendorManagement />
      </ProtectedRoute>
    ),
    children: [
      {
        path: "VendorOrderManagement",
        element: <VendorOrderManagement />,
      },
      {
        path: "VendorAdminManagement",
        element: <VendorAdminManagement />,
      },
      {
        path: "VendorItemsManagement",
        element: <VendorItemsManagement />,
      },
      {
        path: "VendorTableManagement",
        element: <VendorTableManagement />,
      },

      {
        path: "",
        element: <VendorOrderManagement />,
      },
      {
        path: "VendorProfileManagement",
        element: (
          <ProtectedRoute>
            <VendorProfileManagement />
          </ProtectedRoute>
        ),
      },
    ],
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
