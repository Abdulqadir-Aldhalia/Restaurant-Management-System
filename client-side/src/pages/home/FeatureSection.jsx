import React from "react";

const FeatureSection = () => {
  return (
    <section className="features">
      <div className="feature" id="vendors">
        <h2>Browse Vendors</h2>
        <p>Discover and purchase items from a wide range of vendors.</p>
      </div>
      <div className="feature" id="restaurants">
        <h2>Restaurants</h2>
        <p>Order food from top-rated restaurants nearby.</p>
      </div>
      <div className="feature" id="reservations">
        <h2>Reserve a Table</h2>
        <p>Reserve tables in advance at your favorite restaurants.</p>
      </div>
    </section>
  );
};

export default FeatureSection;
