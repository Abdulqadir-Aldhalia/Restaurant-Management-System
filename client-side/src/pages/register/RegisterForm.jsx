import React, { useState } from "react";
import "./registerForm.css";
import Button from "../../components/button/Button.jsx";
import Input from "../../components/input/Input";
import { useNavigate } from "react-router";
import axios from "axios";
import { baseUrl } from "../../const";

function RegisterForm() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");
  const [imageUrl, setImageUrl] = useState("");
  const [imageFile, setImageFile] = useState(null);
  const [errors, setErrors] = useState({});
  const [isChecked, setIsChecked] = useState(false); // State for the terms checkbox
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    let formErrors = {};

    // Frontend validation checks
    if (!name.trim()) formErrors.name = "Name is required";
    if (!email.trim()) formErrors.email = "Email is required";
    if (!password.trim()) formErrors.password = "Password is required";
    if (password !== confirmPassword)
      formErrors.confirmPassword = "Passwords do not match";
    if (!phoneNumber.trim())
      formErrors.phoneNumber = "Phone number is required";

    if (Object.keys(formErrors).length > 0) {
      setErrors(formErrors);
      return;
    }

    // Form data for image file
    const formData = new FormData();
    formData.append("name", name);
    formData.append("email", email);
    formData.append("password", password);
    formData.append("phone", phoneNumber);
    if (imageFile) {
      formData.append("img", imageFile);
    }

    try {
      const response = await axios.post(`${baseUrl}/signup`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      if (response.status === 200 || response.status === 201) {
        alert("Registration successful!");
        navigate("/login");
      } else {
        setErrors({ general: "Failed to register. Please try again." });
      }
    } catch (error) {
      console.log(error);
      if (error.response) {
        const status = error.response.status;

        if (status === 400) {
          const serverErrors = error.response.data.error;
          setErrors({
            ...serverErrors,
            general: serverErrors,
          });
        } else if (status === 409) {
          setErrors({ general: "Email already exists." });
        } else if (status === 500) {
          setErrors({ general: "Server error. Please try again later." });
        } else {
          setErrors({
            general: "An unexpected error occurred. Please try again.",
          });
        }
      } else if (error.request) {
        setErrors({ general: "Network error. Please check your connection." });
      } else {
        setErrors({ general: "An unknown error occurred." });
      }
    }
  };

  const handleImageUpload = (e) => {
    const file = e.target.files[0];
    setImageUrl(URL.createObjectURL(file));
    setImageFile(file);
  };

  return (
    <div className="container">
      <div className="form-container">
        <div className="info-section">
          <h2>INFOMATION</h2>
          <p>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Et
            molestie ac feugiat sed. Diam volutpat commodo.
          </p>
          <p>
            Eu ultrices: Vitae auctor eu augue ut. Malesuada nunc vel risus
            commodo viverra. Praesent elementum facilisis leo vel.
          </p>
        </div>
        <div className="register-section">
          <h2>REGISTER FORM</h2>
          <form onSubmit={handleSubmit}>
            <div className="input-group">
              <Input
                type="text"
                placeholder="Full Name"
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
              {errors.name && <p className="error-message">{errors.name}</p>}
            </div>
            <div className="input-group">
              <Input
                type="email"
                placeholder="Email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              {errors.email && <p className="error-message">{errors.email}</p>}{" "}
            </div>
            <div className="input-group">
              <Input
                type="password"
                placeholder="Password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              {errors.password && (
                <p className="error-message">{errors.password}</p>
              )}{" "}
            </div>
            <div className="input-group">
              <Input
                type="password"
                placeholder="Confirm Password"
                id="confirmPassword"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
              {errors.confirmPassword && (
                <p className="error-message">{errors.confirmPassword}</p>
              )}{" "}
            </div>
            <div className="input-group">
              <Input
                type="text"
                placeholder="Phone Number"
                id="phoneNumber"
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value)}
              />
              {errors.phoneNumber && (
                <p className="error-message">{errors.phoneNumber}</p>
              )}{" "}
            </div>
            <div className="input-group">
              <input type="file" onChange={handleImageUpload} />

              {errors.general && (
                <p className="error-message">{errors.general}</p>
              )}
            </div>
            <div className="checkbox-group">
              <input
                type="checkbox"
                id="terms"
                name="terms"
                onChange={(e) => setIsChecked(e.target.checked)} // Handle checkbox state
                required
              />
              <label htmlFor="terms">
                I agree to the <a href="#">Terms and Conditions</a>
              </label>
            </div>
            <Button title={"Register"} disabled={!isChecked} />{" "}
          </form>
        </div>
      </div>
    </div>
  );
}

export default RegisterForm;
