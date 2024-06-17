import React, { useEffect, useState } from "react";
import Dialog from "../Dialog";

import styles from "./HammerBoardDailog.module.css";

type PropsType = {
  onClose: () => void;
  isOpen?: boolean;
};

interface Leader {
  rank: number;
  icon: string;
  name: string;
  score: string;
}

const HammerBoardDailog = ({ onClose, isOpen }: PropsType) => {
  const [leaders] = useState<Leader[]>([]);
  const [others] = useState<Leader[]>([]);

  return (
    <Dialog title="Chart" onClose={onClose} isOpen={isOpen}>
      <div className={styles["chart-content"]}>
        <div className={styles["leaders-section"]}>
          {leaders.map((leader, index) => (
            <div key={index} className={styles["leader"]}>
              <div className={`${styles["rank"]} rank-${leader.rank}`}>
                <img src={leader.icon} alt={leader.name} />
                <div className={styles["score"]}>{leader.score}</div>
              </div>
            </div>
          ))}
        </div>
        <div className={styles["others-section"]}>
          {others.map((leader, index) => (
            <div key={index} className={styles["leader"]}>
              <div className={styles["rank"]}>{leader.rank}</div>
              <div className={styles["icon"]}>
                <img src={leader.icon} alt={leader.name} />
              </div>
              <div className={styles["score"]}>{leader.score}</div>
            </div>
          ))}
        </div>
      </div>
    </Dialog>
  );
};

export default HammerBoardDailog;
