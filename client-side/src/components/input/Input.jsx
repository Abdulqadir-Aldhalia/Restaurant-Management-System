import styles from "./input.module.css";

const Input = ({ type, placeholder, id, value, onChange }) => (
  <input
    type={type}
    placeholder={placeholder}
    id={id}
    value={value}
    onChange={onChange}
    className={styles["input-64"]}
  />
);

export default Input;
