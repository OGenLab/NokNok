import React, { useEffect, useState } from "react";
import Dialog from "../Dialog";
import styles from "./FrenZoneDailog.module.css";

type PropsType = {
  onClose: () => void;
  isOpen?: boolean;
};

interface InvitationReward {
  icon: string;
  count: number;
  reward: number;
  claimed: boolean;
  onClick: () => void;
}

const FrenZoneDailog = ({ onClose, isOpen }: PropsType) => {
  const [rewards] = useState<InvitationReward[]>([]);
  const [invitations] = useState(0);

  const onInvite = () => {};

  return (
    <Dialog title="Fren Zone" onClose={onClose} isOpen={isOpen}>
      <div className={styles["fren-zone-content"]}>
        <div className={styles["invite-section"]}>
          <h3>Invite frens to get bonuses</h3>
          <div className={styles["invite-rule"]}>
            <div className={styles["rule"]}>
              Invite fren
              <div>
                <span></span>
                +100 For you and fren
              </div>
            </div>
            <div className={styles["rule"]}>
              Fren with Premium
              <div>
                <span></span>
                +500 For you and fren
              </div>
            </div>
          </div>
        </div>
        <div className={styles["cumulative-invitation"]}>
          <h3>Cumulative invitation</h3>
          <div className={styles["rewards"]}>
            <div className={styles["invitations-count"]}>
              <div className={styles["current-invitations"]}>{invitations}</div>
            </div>
            {rewards.map((reward, index) => (
              <div
                key={index}
                className={`${styles["reward"]} ${reward.claimed ? styles["claimed"] : ""}`}
                onClick={reward.onClick}
              >
                <div className={styles["reward-icon"]}>
                  <img src={reward.icon} alt="" />
                  <div className={styles["reward-amount"]}>{reward.reward}</div>
                </div>
                <div className={styles["reward-info"]}>
                  <div className={styles["reward-count"]}>{reward.count}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
        <div className={styles["invite-btn"]}>
          <button className={styles["invite-button"]} onClick={onInvite}>
            INVITE A FREN
          </button>
        </div>
      </div>
    </Dialog>
  );
};

export default FrenZoneDailog;
