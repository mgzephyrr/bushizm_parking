import React from "react";
import styles from "./footer.module.css"

interface FooterProps {
    children?: React.ReactNode;
}

export default function Footer({children}: FooterProps) {
    return <div className={styles.footer}>{children}</div>;
}