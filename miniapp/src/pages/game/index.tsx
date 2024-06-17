import React, { useEffect } from "react";
import styles from "./index.module.css";

const GamePage: React.FC = () => {
  return (
    <div className={styles["game-page-wrap"]}>
      <div className={styles["top-info"]}>
        <div className={styles["coin-info"]}>
          <div className={styles["coin-bar"]}></div>
          <div className={styles["energy-bar"]}></div>
        </div>
        <div className={styles["mole-info"]}>
          <div className={styles["hammer-icon"]}></div>
          <div className={styles["statis-data"]}></div>
          <div className={styles["rank-list"]}></div>
        </div>
      </div>
      <div className={styles["game-area"]} id="game-area"></div>
      <div className={styles["bottom-toolbar"]}>
        <div className={styles["hammer-list"]}></div>
        <div className={styles["boost-info"]}></div>
        <div className={styles["treasure-chest"]}></div>
        <div className={styles["earn-list"]}></div>
        <div className={styles["fren-list"]}></div>
      </div>
    </div>
  );
};

export default GamePage;
