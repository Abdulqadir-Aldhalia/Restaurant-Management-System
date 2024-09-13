import styles from "./label.module.css";

const Label = ({ htmlFor, children }) => {
  return (
    <label className={styles["label-64"]} htmlFor={htmlFor}>
      {children}
    </label>
  );
};

export default Label;
