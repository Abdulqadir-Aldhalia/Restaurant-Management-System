import styles from "./container.module.css";

const Container = ({ children }) => {
  return (
    <div className={styles["container-64-wrapper"]}>
      <div className={styles["container-64"]}>{children}</div>
    </div>
  );
};

export default Container;
