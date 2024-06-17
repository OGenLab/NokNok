import React from "react";
import styles from "./index.module.css";

export type HammerInfo = {
  id: number;
  equipmentCount: number;
  name: string;
  imageUrl: string;
  count: number;
};

type HammerItemPropsType = {
  isSelected: boolean;
  hammerInfo?: HammerInfo;
  onSelectFn?: (info: HammerInfo | undefined) => void;
};

const HammerItem = ({ hammerInfo, onSelectFn }: HammerItemPropsType) => {
  const selectHammer = () => {
    if (onSelectFn) {
      onSelectFn(hammerInfo);
    }
  };

  return (
    <div className={styles["hammer-item"]} onClick={selectHammer}>
      {hammerInfo?.equipmentCount === 1 && (
        <div className={styles["equip-status"]}></div>
      )}
      <div className={styles["hammer-content"]}>
        <div className={styles["hammer-image"]}>
          <image />
          <div className={styles["hammer-count"]}></div>
        </div>
        <div className={styles["hammer-name"]}></div>
      </div>
    </div>
  );
};

export default HammerItem;
