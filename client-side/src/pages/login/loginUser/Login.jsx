import React, { useState } from "react";
import axios from "axios";
import { setUserToken } from "../../../redux/user/userSlice.js";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router";
import Button from "../../../components/button/Button.jsx";
import Input from "../../../components/input/Input.jsx";
import Container from "../../../components/container/Container.jsx";
import { baseUrl } from "../../../const.js";
import "./login.css";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setErrors({}); // Clear previous errors
    let formErrors = {};

    const trimmedEmail = email.trim();
    const trimmedPassword = password.trim();

    if (!trimmedEmail) {
      formErrors.email = "Email is required";
    }

    if (!trimmedPassword) {
      formErrors.password = "Password is required";
    }

    console.log(trimmedPassword + " " + trimmedEmail);

    if (Object.keys(formErrors).length > 0) {
      setErrors(formErrors);
      setLoading(false);
      return;
    }

    try {
      const response = await axios.post(
        `${baseUrl}/signin?username=${encodeURIComponent(trimmedEmail)}&password=${encodeURIComponent(trimmedPassword)}`,
        {},
      );
      if (response.status === 200) {
        dispatch(setUserToken(response.data.token));
        localStorage.setItem("token", response.data.token);
        navigate("/userHome");
      } else if (response.data.message) {
        console.log("Login message:", response.data.message);
      } else {
        console.log("Unexpected response structure:", response.data);
      }
    } catch (error) {
      if (error.response?.status === 400) {
        setErrors({
          general: error.response.data.error,
        });
      } else if (error.response?.status === 404) {
        setErrors({
          general: error.response.data.error,
        });
      } else {
        // General error message for other cases
        setErrors({
          general: "Oops! Something went wrong. Please try again later.",
        });
      }
      console.error("Login error:", error.response?.data || error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container>
      <div className="welcome-message">
        <h2>Welcome</h2>
        <p>Please login to continue</p>
      </div>

      <form onSubmit={handleSubmit}>
        <Input
          type="email"
          placeholder="Enter your email"
          id="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        {errors.email && <p className="error-message">{errors.email}</p>}

        <Input
          type="password"
          placeholder="Enter your password"
          id="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        {errors.password && <p className="error-message">{errors.password}</p>}

        {errors.general && <p className="error-message">{errors.general}</p>}

        <Button
          title={loading ? "Logging in..." : "Login"}
          disabled={loading}
        />

        <div className="register-option">
          <span>
            Don't have an account? <a href="/register">Create one here</a>
          </span>
        </div>
      </form>
    </Container>
  );
};

export default Login;
