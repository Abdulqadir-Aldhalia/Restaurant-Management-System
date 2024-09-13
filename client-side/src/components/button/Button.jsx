import styles from "./button.module.css";

const Button = ({ title }) => {
  return (
    <button className={styles["button-64"]} role="button">
      <span className={styles.text}>{title}</span>
    </button>
  );
};

export default Button;
