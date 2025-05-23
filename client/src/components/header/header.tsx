import React from "react";
import styles from "./header.module.css";

interface IHeaderProps {
    children?: React.ReactNode;
}

export default function Header({children}: IHeaderProps) {
    return (
        <div className={styles.header}>
            {children}
        </div>
    )
}