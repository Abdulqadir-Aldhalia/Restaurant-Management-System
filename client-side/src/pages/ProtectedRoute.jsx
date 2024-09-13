import { useSelector } from "react-redux";
import { Navigate, useLocation } from "react-router-dom";

const ProtectedRoute = ({ children }) => {
  const userToken = useSelector((state) => state.user.userToken); // Get userToken from Redux store
  const location = useLocation();

  if (!userToken) {
    if (location.pathname === "/adminDashboard") {
      return <Navigate to="/loginAdminPortal" />;
    }
    return <Navigate to="/login" />;
  }

  return children;
};

export default ProtectedRoute;
