import React, { useEffect, useState } from "react";
import Dialog from "../Dialog";
import { getBoostInfo } from "src/services/user";
import styles from "./MultitapDailog.module.css";
import { BoostInfo } from "src/services/user.d";
import { useUserInfo } from "src/hooks/useUserInfo";

type PropsType = {
  onClose: () => void;
  isOpen?: boolean;
};

const MultitapDailog = ({ onClose, isOpen }: PropsType) => {
  const { userData } = useUserInfo();
  const [curBoostInfo, setCurBoostInfo] = useState<BoostInfo>();
  const [nextBoostInfo, setNextBoostInfo] = useState<BoostInfo>();
  const onLevelUp = () => {};
  useEffect(() => {
    if (isOpen) {
      getBoostDetailData();
    }
  }, [isOpen]);

  const getBoostDetailData = async () => {
    try {
      // 获取升级的列表详情
      const { boosts } = await getBoostInfo();
      setCurBoostInfo(
        boosts.find(item => item.boost === userData?.accountInfo?.boost)
      );
      setNextBoostInfo(
        boosts.find(item => item.boost === userData?.accountInfo?.nextBoost)
      );
      console.log(boosts, "boost");
    } catch (err) {
      console.log("get boost info: ", err);
    }
  };
  return (
    <Dialog title="Multitap" onClose={onClose} isOpen={isOpen}>
      <div className={styles["multitap-dailog-wrap"]}>
        <div className={styles["boost-content"]}>
          <div className={styles["boost-tip-title"]}>
            <span role="img" aria-label="hammer">
              🔨
            </span>
            Hammer count for one tap
          </div>
          <div className={styles["level-info"]}>
            <div className={styles["current-level"]}>
              <span>Lv.{curBoostInfo?.boost}</span>
              <span>
                <span role="img" aria-label="hammer">
                  🔨
                </span>
                x{curBoostInfo?.coinsRate}
              </span>
            </div>
            <div className={styles["arrow"]}>➜</div>
            <div className={styles["next-level"]}>
              <span>Lv.{nextBoostInfo?.boost}</span>
              <span>
                <span role="img" aria-label="hammer">
                  🔨
                </span>
                x{nextBoostInfo?.coinsRate}
              </span>
            </div>
          </div>
        </div>
        {curBoostInfo?.boost !== 5 && (
          <button className="level-up-button" onClick={onLevelUp}>
            <span className="cost">
              <span role="img" aria-label="coin">
                🪙
              </span>
              {curBoostInfo?.neededCoins}
            </span>
            LEVEL UP
          </button>
        )}
      </div>
    </Dialog>
  );
};

export default MultitapDailog;
