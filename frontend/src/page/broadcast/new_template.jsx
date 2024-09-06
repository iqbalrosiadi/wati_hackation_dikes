import styles from "./new_template.module.css";


function TemplateMessage() {
  return (
<div className={styles.container}>
      <section className={styles.templateMessage}>
        <div className={styles.templateHeader}>
          <h2>Create Template Message</h2>
          <button className={styles.useSampleBtn}>Use a sample</button>
        </div>

        <div className={styles.templateForm}>
          <label htmlFor="templateName">Template Name</label>
          <input type="text" id="templateName" className={styles.inputField} />

          <label htmlFor="category">Category</label>
          <select id="category" className={styles.inputField}>
            <option value="marketing">Marketing</option>
          </select>

          <label htmlFor="language">Language</label>
          <select id="language" className={styles.inputField}>
            <option value="en">English</option>
          </select>
        </div>

        <div className={styles.broadcastSection}>
          <label htmlFor="broadcastTitle">Broadcast title (Optional)</label>
          <input type="text" id="broadcastTitle" className={styles.inputField} />
        </div>

        <div className={styles.messageBody}>
          <label htmlFor="body">Body</label>
          <textarea id="body" rows="4" className={styles.inputField}></textarea>
          <button className={styles.addVariableBtn}>Add Variable</button>
        </div>

        <div className={styles.footerSection}>
          <label htmlFor="footer">Footer (Optional)</label>
          <input type="text" id="footer" className={styles.inputField} placeholder="Powered by wati.io" />
        </div>

        <div className={styles.buttonSection}>
          <label>Buttons (Recommended)</label>
          <button className={styles.addButtonBtn}>Add Button</button>
        </div>

        <div className={styles.submitSection}>
          <button className={styles.saveDraftBtn}>Save as Draft</button>
          <button className={styles.submitBtn}>Save and Submit</button>
        </div>
      </section>

      <section className={styles.preview}>
        <h3>Preview</h3>
        <div className={styles.previewWindow}>
          <img src="preview_image.png" alt="Preview of template message" />
        </div>
      </section>
    </div>
  );
}

export default TemplateMessage;
