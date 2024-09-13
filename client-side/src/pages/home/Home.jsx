import React from "react";
import Header from "./Header";
import HeroSection from "./HeroSection";
import FeatureSection from "./FeatureSection";
import Footer from "./Footer";
import "./home.css";

function Home() {
  return (
    <div>
      <Header />
      <HeroSection />
      <FeatureSection />
      <Footer />
    </div>
  );
}

export default Home;
