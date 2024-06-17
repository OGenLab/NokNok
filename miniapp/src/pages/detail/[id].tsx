import { useRouter } from "next/router";
import styles from "./index.module.css";

export default function Detail() {
  const router = useRouter();
  const { id } = router.query;
  return <main className={styles.main}>id: {id}</main>;
}
