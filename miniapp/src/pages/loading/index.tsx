import React, { useState } from "react";
import { useRouter } from "next/router";
import styles from "./index.module.css";

const LoadingPage: React.FC = () => {
  const router = useRouter();
  const [progress, setProgress] = useState(0);

  return (
    <div className={styles.loadingContainer}>
      <div className={styles["game-logo"]}>loading logo</div>
      <div className={styles.progressBar}>
        <div
          className={styles.progress}
          style={{ width: `${progress}%` }}
        ></div>
      </div>
      <div className={styles["loading-percent"]} id="percent-text">
        Loading...0%
      </div>
    </div>
  );
};

export default LoadingPage;
