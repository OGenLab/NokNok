import React, { useEffect, useState } from "react";
import { useRouter } from "next/router";
import styles from "./index.module.css";
import { login, getUserInfo, getUserHammer } from "src/services/user";
import { initInitData } from "@tma.js/sdk";
import { useUserInfo } from "src/hooks/useUserInfo";

export default function Home() {
  const { setUserData } = useUserInfo();
  const router = useRouter();

  const parseStr = (obj: any) => {
    return Object.keys(obj)
      .map(key => {
        const val =
          typeof obj[key] === "object" ? JSON.stringify(obj[key]) : obj[key];
        return `${key}=${val}`;
      })
      .join("\n");
  };

  const initLoadData = async () => {
    const tgWebAppData = initInitData();
    const tgWebAppDataRaw = parseStr(tgWebAppData);
    const result = tgWebAppDataRaw.replace("initData=", "");
    const data = JSON.parse(result);
    console.log("tg web app data:", data);

    try {
      await login({
        token: `${JSON.stringify(data)}`,
        referralCode: data.startParam,
        channel: "",
      });

      const res = await getUserInfo();
      console.log("user info:", res);
      setUserData(res);
    } catch (err) {
      console.log("init load:", err);
    }
  };

  useEffect(() => {
    initLoadData();
  }, []);

  return (
    <main className={styles.main} onClick={() => router.push("/hammer")}>
      1234
    </main>
  );
}
