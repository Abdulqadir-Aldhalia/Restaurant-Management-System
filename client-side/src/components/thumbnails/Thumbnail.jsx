import React, { useState } from "react";
import styles from "./thumbnail.module.css";

const Thumbnail = ({ onImageChange }) => {
  const [imagePreview, setImagePreview] = useState(null);

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result);
        onImageChange(file);
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <div className={styles["thumbnail-container"]}>
      {imagePreview ? (
        <img
          src={imagePreview}
          alt="Thumbnail Preview"
          className={styles.thumbnail}
        />
      ) : (
        <div className={styles["thumbnail-placeholder"]}>
          <p>Image Preview</p>
        </div>
      )}
      <input type="file" accept="image/*" onChange={handleImageChange} />
    </div>
  );
};

export default Thumbnail;
