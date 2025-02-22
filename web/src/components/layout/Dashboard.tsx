import { ReactNode } from "react";
import { Navbar } from "../shared/Navbar";
import styles from "./Dashboard.module.scss";

interface DashboardProps {
  children: ReactNode;
}

export const Dashboard = ({ children }: DashboardProps) => {
  return (
    <div className={styles.dashboardContainer}>
      <Navbar />
      <main className={styles.content}>{children}</main>
    </div>
  );
};
