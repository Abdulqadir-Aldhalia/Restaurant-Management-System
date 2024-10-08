import { useSelector } from "react-redux";
import { Navigate, useLocation } from "react-router-dom";

const ProtectedRoute = ({ children }) => {
  const userToken = useSelector((state) => state.user.userToken); // Get userToken from Redux store
  const location = useLocation();

  if (!userToken) {
    if (location.pathname.startsWith("/adminDashboard")) {
      return <Navigate to="/loginAdminPortal" />;
    } else if (location.pathname.startsWith("/vendorApp")) {
      return <Navigate to="/loginVendorPortal" />;
    }
    return <Navigate to="/login" />;
  }

  return children;
};

export default ProtectedRoute;
