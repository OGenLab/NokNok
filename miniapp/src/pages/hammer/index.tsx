import HammerItem from "src/components/HammerItem";
import styles from "./index.module.css";
import { useEffect, useState } from "react";
import { HammerInfo } from "src/components/HammerItem";
import { getUserHammer } from "src/services/user";
import { equipHammer } from "src/services/battle";

export default function HammerPage() {
  const [curEquipHammerId, setCurEquipHammerId] = useState<number>(0);
  const [curSelectHammer, setCurSelectHammer] = useState<HammerInfo>();

  useEffect(() => {
    // 获取用户锤子详情
    getHammerInfo();
  }, []);

  const getHammerInfo = async () => {
    try {
      const hammers = await getUserHammer();
      console.log("hammer list:", hammers);
    } catch (err) {
      console.log("hammer list request err: ", err);
    }
  };

  const selectHammer = (info: HammerInfo | undefined) => {
    if (info) {
      setCurSelectHammer(info);
    }
  };

  const handleEquip = async () => {
    if (curSelectHammer?.id) {
      await equipHammer({
        hammerId: curSelectHammer?.id,
      });
    }
  };

  return (
    <div className={styles["hammer-page-wrap"]}>
      <div className={styles["hammer-title"]}>Hammer</div>
      <div className={styles["hammer-info"]}>
        <div className={styles["hammer-area"]}></div>
        <div className={styles["hammer-content"]}>wood</div>
      </div>
      <div className={styles["hammer-list"]}>
        <div className={styles["select-box"]}></div>
        <div className={styles["hammer-data"]}>
          {/* <HammerItem isSelected={false} onSelectFn={selectHammer}></HammerItem> */}
        </div>
      </div>
      <div className={styles["bottom-btn"]} onClick={handleEquip}>
        {curSelectHammer?.id !== curEquipHammerId ? "EQUIP" : "EQUIPED"}
      </div>
    </div>
  );
}
