import React, { useEffect, useState } from "react";
import Dialog from "../Dialog";
import styles from "./EarnDailog.module.css";
import { getTasks } from "src/services/battle";
import { TaskInfo } from "src/services/battle.d";

type PropsType = {
  onClose: () => void;
  isOpen?: boolean;
};

interface Task {
  icon: string;
  description: string;
  reward: number;
  status: "not-started" | "in-progress" | "completed";
  actionLabel: string;
  action: () => void;
}

const EarnDailog = ({ onClose, isOpen }: PropsType) => {
  const [dailyTasks, setDailyTasks] = useState<TaskInfo[]>([]);
  const [basicTasks, setBasicTasks] = useState<TaskInfo[]>([]);

  useEffect(() => {
    // 获取任务详情
    if (isOpen) {
      getTaskDetail();
    }
  }, [isOpen]);

  const getTaskDetail = async () => {
    const { taskInfos } = await getTasks();
    const dailys = taskInfos.filter(item => item.type === 1);
    setDailyTasks(dailys);
    const basics = taskInfos.filter(item => item.type === 2);
    setBasicTasks(basics);
    console.log("task list data:", taskInfos);
  };

  return (
    <Dialog title="Earn" onClose={onClose} isOpen={isOpen}>
      {/* <div className={styles["earn-content"]}>
        <div className={styles["task-section"]}>
          <h3>Daily</h3>
          <ul className={styles["task-list"]}>
            {dailyTasks.map((task, index) => (
              <li
                key={index}
                className={`${styles["task-item"]} ${styles[task.status]}`}
              >
                <div className={styles["task-icon"]}>
                  <img src={task.icon} alt="" />
                </div>
                <div className={styles["task-description"]}>
                  {task.description}
                </div>
                <div className={styles["task-reward"]}>
                  <span role="img" aria-label="coin"></span>
                  {task.reward}
                </div>
                <button className={styles["task-action"]} onClick={task.action}>
                  {task.actionLabel}
                </button>
              </li>
            ))}
          </ul>
        </div>
        <div className={styles["task-section"]}>
          <h3>Basic tasks</h3>
          <ul className={styles["task-list"]}>
            {basicTasks.map((task, index) => (
              <li
                key={index}
                className={`${styles["task-item"]} ${styles[task.status]}`}
              >
                <div className={styles["task-icon"]}>
                  <img src={task.icon} alt="" />
                </div>
                <div className={styles["task-description"]}>
                  {task.description}
                </div>
                <div className={styles["task-reward"]}>
                  <span role="img" aria-label="coin"></span> {task.reward}
                </div>
                <button className={styles["task-action"]} onClick={task.action}>
                  {task.actionLabel}
                </button>
              </li>
            ))}
          </ul>
        </div>
      </div> */}
      <div></div>
    </Dialog>
  );
};

export default EarnDailog;
