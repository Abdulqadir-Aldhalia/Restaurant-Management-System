import React from "react";
import { useNavigate } from "react-router";

const HeroSection = () => {
  const navigate = useNavigate();
  return (
    <section className="hero">
      <div className="hero-content">
        <h1>Welcome to App</h1>
        <p>
          Browse vendors, purchase items, order food, and reserve tables—all in
          one place!
        </p>
        <a onClick={() => navigate("/userHome")} className="button">
          Start Exploring
        </a>
      </div>
    </section>
  );
};

export default HeroSection;
