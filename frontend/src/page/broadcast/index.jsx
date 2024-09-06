import styles from "./broadcast.module.css";
import { useState } from "react";
import TemplateMessage from "./new_template"

const BroadcastPage = () => {
  const [showSecondComponent, setShowSecondComponent] = useState(false);

  return (
    <div>
      <div className={styles.templateContainer}>
        {/* <div className={styles.alert}>
            <p>
              All templates must adhere to WhatsApp's Template Message
              Guidelines. Click here to read.
            </p>
          </div> */}

        <div className={styles.sidebar}>
          <h2 className={styles.sidebarTitle}>Template Messages</h2>
          <ul>
            <li>Broadcast Analytics</li>
            <li>Scheduled Broadcasts</li>
          </ul>
        </div>

        <div className={styles.mainContent}>
          <div className={styles.header}>
            <div className={styles.headerLeft}>
              <h1>Template Messages</h1>
            </div>
            <div className={styles.headerRight}>
              <button className={styles.primaryButton}>Watch Tutorial</button>
              <button className={styles.primaryButton} onClick={() => setShowSecondComponent(true)}>
                New Template Message
              </button>
            </div>
          </div>

          <table className={styles.templateTable}>
            <thead>
              <tr>
                <th>Template Name</th>
                <th>Category</th>
                <th>Status</th>
                <th>Language</th>
                <th>Last Updated</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>onboarding_signoff</td>
                <td>Utility</td>
                <td className={styles.statusApproved}>Approved</td>
                <td>English</td>
                <td>23/5/2024</td>
                <td className={styles.actions}>
                  <button className={styles.actionButton}>
                    Send Broadcast
                  </button>
                  <button className={styles.deleteButton}>Delete</button>
                </td>
              </tr>
              <tr>
                <td>new_chat_v1</td>
                <td>Utility</td>
                <td className={styles.statusApproved}>Approved</td>
                <td>English</td>
                <td>23/5/2024</td>
                <td className={styles.actions}>
                  <button className={styles.actionButton}>
                    Send Broadcast
                  </button>
                  <button className={styles.deleteButton}>Delete</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      {showSecondComponent ? (
        <div className={styles.overlay}>
          <div className={styles.popup}>
            <TemplateMessage/>
            <button className={styles.closeButton} onClick = {() => setShowSecondComponent(false)}>Close</button>
          </div>
        </div>
      ) : (
        ""
      )}
    </div>
  );
};

export default BroadcastPage;
