import React, { ReactNode } from "react";
import styles from "./index.module.css";
// import CloseIcon from "@/assets/svg/CloseIcon";

interface ModalProps {
  title: string;
  children: ReactNode;
  actions?: ReactNode;
  isOpen?: boolean;
  onClose: () => void;
}

const Dialog: React.FC<ModalProps> = ({
  title,
  children,
  actions,
  isOpen,
  onClose,
}) => {
  if (!isOpen) return null;

  return (
    <div className={styles["modal-overlay"]}>
      <div className={styles["modal-container"]}>
        <div className={styles["modal-header"]}>
          <h2>{title}</h2>
          <button className={styles["close-button"]} onClick={onClose}>
            ×
          </button>
        </div>
        <div className={styles["modal-content"]}>{children}</div>
        <div className={styles["modal-actions"]}>{actions}</div>
      </div>
    </div>
  );
};

export default Dialog;
